# 自动部署说明

本文档说明当前项目如何通过 GitHub Actions 在 `main` 分支推送后，自动构建 `server/` 与 `web/`，上传服务器并执行远端部署脚本。

重要边界：

- GitHub Actions 只负责后续代码包更新，不负责生成生产密钥、数据库配置或真实域名配置。
- 第一次部署前，必须手动准备服务器上的 `.env`、`config.yaml`、数据库、systemd 和 Nginx。
- 真实 `.env`、真实 `config.yaml`、数据库密码、JWT key、域名证书都不能提交到 GitHub。

## 1. 服务器目录结构

推荐使用以下目录结构：

```text
/opt/meimei-member/
├── backups/
├── deploy/
│   ├── deploy.sh
│   ├── member-api.service
│   └── nginx.conf.example
├── release.tar.gz
├── releases/
├── server/
│   ├── .env
│   ├── config.yaml
│   └── member-api
└── web/
    └── dist/
```

说明：

- `server/.env` 只保存在服务器，不提交 GitHub。
- `server/config.yaml` 只保存在服务器，不提交 GitHub。
- 可以参考仓库内的 `server/.env.example` 和 `server/config.prod.example.yaml` 手动生成真实生产配置。

## 2. GitHub Secrets 配置

进入仓库 `Settings -> Secrets and variables -> Actions`，新增以下 Secrets：

- `DEPLOY_HOST`: 服务器公网 IP 或域名
- `DEPLOY_PORT`: SSH 端口，例如 `22`
- `DEPLOY_USER`: 执行部署的 SSH 用户
- `DEPLOY_SSH_KEY`: 该用户对应的私钥全文

说明：

- 工作流会上传 `release.tar.gz` 与 `deploy/` 下的部署文件。
- 工作流随后会远程执行 `/opt/meimei-member/deploy/deploy.sh`。
- Secrets 只解决 SSH 连接问题，不会替你创建生产 `.env`、`config.yaml` 或数据库。

## 3. 首次手动部署步骤

第一次部署必须先手动完成以下准备，再让 GitHub Actions 接手后续更新。

1. 创建服务器目录：

```bash
sudo mkdir -p /opt/meimei-member/{deploy,server,web,backups,releases}
sudo chown -R <deploy-user>:<deploy-user> /opt/meimei-member
```

2. 提前创建数据库：

- 在 MySQL 中手动创建业务库。
- `server/config.prod.example.yaml` 里的 `mysql.db-name`、`username`、`password` 需要替换成真实值。

3. 手动创建 `/opt/meimei-member/server/.env`：

可以先参考仓库模板：

```bash
cp server/.env.example /opt/meimei-member/server/.env
```

然后确认内容至少包含：

```bash
GIN_MODE=release
GVA_CONFIG=/opt/meimei-member/server/config.yaml
```

4. 手动创建 `/opt/meimei-member/server/config.yaml`：

可以先参考仓库模板：

```bash
cp server/config.prod.example.yaml /opt/meimei-member/server/config.yaml
```

然后至少修改这些真实值：

- `jwt.signing-key`
- `mysql.path`
- `mysql.port`
- `mysql.db-name`
- `mysql.username`
- `mysql.password`

5. 手动上传部署文件到服务器：

- `deploy/deploy.sh`
- `deploy/member-api.service`
- `deploy/nginx.conf.example`

6. 安装 systemd 服务：

```bash
sudo cp /opt/meimei-member/deploy/member-api.service /etc/systemd/system/member-api.service
sudo systemctl daemon-reload
sudo systemctl enable member-api
```

`member-api.service` 当前约定为：

- `WorkingDirectory=/opt/meimei-member/server`
- `ExecStart=/opt/meimei-member/server/member-api`
- `EnvironmentFile=/opt/meimei-member/server/.env`

7. 安装 Nginx 配置：

```bash
sudo cp /opt/meimei-member/deploy/nginx.conf.example /etc/nginx/conf.d/meimei-member.conf
```

安装前必须手动修改：

- 将 `your-domain.com` 替换成真实域名
- 将证书路径替换成真实证书文件路径

然后执行：

```bash
sudo nginx -t
sudo systemctl reload nginx
```

8. 在 GitHub 配置好四个 Secrets 后，再执行：

```bash
git push origin main
```

## 4. 后续自动部署步骤

首次部署准备完成后，后续流程如下：

1. 本地开发并提交代码
2. 推送到 GitHub `main`
3. GitHub Actions 自动执行：
   - 构建 `server/member-api`
   - 构建 `web/dist`
   - 打包 `release.tar.gz`
   - 上传到 `/opt/meimei-member`
   - 执行 `/opt/meimei-member/deploy/deploy.sh`
4. 远端脚本只替换二进制和前端静态文件，不会覆盖真实 `.env` 和 `config.yaml`

## 5. deploy.sh 的前置检查

当前部署脚本会在真正替换文件前检查：

- `/opt/meimei-member/server/.env` 是否存在
- `/opt/meimei-member/server/config.yaml` 是否存在
- `/etc/systemd/system/member-api.service` 是否存在

如果上述任一项缺失，脚本会直接报错退出，并提示你先手动准备。

另外：

- 如果 `/etc/nginx/conf.d/meimei-member.conf` 尚未安装，脚本只会提示并跳过 `nginx -t` 与 reload
- 脚本不会自动覆盖 Nginx 配置，更不会生成真实域名或证书配置

## 6. Nginx 与 router-prefix 说明

当前后端配置模板约定：

- `system.router-prefix: ""`

当前 Nginx 示例约定：

- 前端请求走 `/api/`
- Nginx 通过 `rewrite ^/api/?(.*)$ /$1 break;` 去掉 `/api`
- 再反代到本地 `127.0.0.1:8888`

因此两者是配套的，若你未来修改 `router-prefix`，需要同步修改 Nginx 反代规则。

## 7. 回滚方式

每次部署前，脚本都会备份旧版本到：

```text
/opt/meimei-member/backups/<timestamp>/
```

回滚步骤：

1. 查看备份目录：

```bash
ls -lah /opt/meimei-member/backups
```

2. 回滚后端二进制：

```bash
cp /opt/meimei-member/backups/<timestamp>/server/member-api /opt/meimei-member/server/member-api
```

3. 回滚前端静态文件：

```bash
rm -rf /opt/meimei-member/web/dist
cp -a /opt/meimei-member/backups/<timestamp>/web/dist /opt/meimei-member/web/dist
```

4. 重新加载服务：

```bash
sudo systemctl restart member-api
sudo nginx -t
sudo systemctl reload nginx
```

## 8. 必须留在服务器、不能进 GitHub 的内容

以下内容必须只保留在服务器：

- `/opt/meimei-member/server/.env`
- `/opt/meimei-member/server/config.yaml`
- MySQL 真实连接信息
- JWT 真实签名密钥
- 域名证书文件
- 微信 `AppSecret`
- 任何生产环境私钥、密码、Token

## 9. 当前边界

当前自动部署方案明确不包含以下内容：

- 不做 webhook 监听扩展
- 不做 Docker 镜像构建
- 不做蓝绿部署
- 不自动上传微信小程序体验版
