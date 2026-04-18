# 自动部署说明

本文档说明当前项目如何通过 GitHub Actions 在 `main` 分支推送后，自动构建 `server/` 和 `web/`，再经由 SSH 上传到服务器执行部署。

## 1. 服务器目录结构

推荐服务器目录如下：

```text
/opt/meimei-member/
├── backups/                # 每次部署前自动备份旧版本
├── deploy/
│   ├── deploy.sh
│   ├── member-api.service
│   └── nginx.conf.example
├── release.tar.gz          # GitHub Actions 上传的发布包
├── releases/               # 每次部署解压目录
├── server/
│   ├── .env                # 仅服务器保存，不提交 GitHub
│   ├── config.yaml         # 后端运行配置
│   └── member-api          # Go 后端二进制
└── web/
    └── dist/               # GVA 管理后台静态文件
```

说明：

- `miniprogram/` 第一阶段仍然手动上传体验版，不进入自动部署链路。
- `.env`、数据库密码、微信 `AppSecret` 只放在服务器，不提交到仓库。

## 2. GitHub Secrets 配置

在 GitHub 仓库 `Settings -> Secrets and variables -> Actions` 中新增以下 Secrets：

- `DEPLOY_HOST`: 服务器公网 IP 或域名
- `DEPLOY_PORT`: SSH 端口，默认常见为 `22`
- `DEPLOY_USER`: 用于部署的 SSH 用户
- `DEPLOY_SSH_KEY`: 对应私钥内容

建议使用具备部署目录权限的专用账号；如果部署脚本需要执行 `systemctl` 和 `nginx -t`，该账号应具备对应权限。

## 3. SSH Key 配置

本地生成部署密钥：

```bash
ssh-keygen -t ed25519 -C "github-actions@meimei-member"
```

将生成的公钥内容追加到服务器目标用户的 `~/.ssh/authorized_keys`：

```bash
mkdir -p ~/.ssh
chmod 700 ~/.ssh
cat id_ed25519.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

将私钥内容完整复制到 GitHub Secret `DEPLOY_SSH_KEY`。

## 4. systemd 安装

先把服务文件安装到 systemd：

```bash
sudo cp /opt/meimei-member/deploy/member-api.service /etc/systemd/system/member-api.service
sudo systemctl daemon-reload
sudo systemctl enable member-api
```

在 `/opt/meimei-member/server/.env` 中至少放入运行时环境变量，例如：

```bash
GIN_MODE=release
GVA_CONFIG=/opt/meimei-member/server/config.yaml
```

然后确保后端配置文件存在：

```bash
cp server/config.yaml /opt/meimei-member/server/config.yaml
```

注意：

- 当前 Nginx 方案约定后端接口前缀为 `/api`
- 仓库里的 `server/config.yaml` 已调整为 `system.router-prefix: /api`
- 如果你在服务器上自定义了配置文件，也请保持该值一致

## 5. Nginx 配置

将示例配置复制到 Nginx：

```bash
sudo cp /opt/meimei-member/deploy/nginx.conf.example /etc/nginx/conf.d/meimei-member.conf
sudo nginx -t
sudo systemctl reload nginx
```

说明：

- `root` 指向 `/opt/meimei-member/web/dist`
- `/api/` 反代到本地 `127.0.0.1:8888`
- `try_files` 已支持 Vue history fallback

如果你有域名和 HTTPS，可以在此基础上追加证书配置。

## 6. 首次部署步骤

首次部署建议按下面顺序执行：

1. 在服务器创建目录：

```bash
sudo mkdir -p /opt/meimei-member/{deploy,server,web,backups,releases}
sudo chown -R <deploy-user>:<deploy-user> /opt/meimei-member
```

2. 将以下文件先手动上传到服务器：

- `deploy/deploy.sh`
- `deploy/member-api.service`
- `deploy/nginx.conf.example`
- `server/config.yaml`

3. 在服务器创建 `/opt/meimei-member/server/.env`，填写运行时环境变量。

4. 安装并启用 systemd 服务。

5. 安装并启用 Nginx 配置。

6. 在 GitHub 配置好四个 Secrets。

7. 本地开发完成后，推送到 `main`：

```bash
git push origin main
```

GitHub Actions 会自动：

- 构建 `server/member-api`
- 构建 `web/dist`
- 打包 `release.tar.gz`
- 通过 SSH 上传到 `/opt/meimei-member`
- 远程执行 `/opt/meimei-member/deploy/deploy.sh`

## 7. 回滚方式

每次部署前，脚本都会把旧版本备份到：

```text
/opt/meimei-member/backups/<timestamp>/
```

回滚步骤：

1. 查看备份目录：

```bash
ls -lah /opt/meimei-member/backups
```

2. 选择目标时间戳，例如 `20260418173000`

3. 回滚后端二进制：

```bash
cp /opt/meimei-member/backups/20260418173000/server/member-api /opt/meimei-member/server/member-api
```

4. 回滚前端静态文件：

```bash
rm -rf /opt/meimei-member/web/dist
cp -a /opt/meimei-member/backups/20260418173000/web/dist /opt/meimei-member/web/dist
```

5. 重启服务并重新加载 Nginx：

```bash
sudo systemctl restart member-api
sudo nginx -t
sudo systemctl reload nginx
```

## 8. 当前边界

当前自动部署方案明确不包含以下内容：

- 不监听 webhook 以外的额外部署入口
- 不做 Docker 镜像构建
- 不做蓝绿部署
- 不自动上传微信小程序体验版
