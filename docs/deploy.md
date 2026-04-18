# 自动部署说明

本文档说明当前项目如何通过 GitHub Actions 在 `main` 分支推送后，自动构建 GVA 的 `server/` 与 `web/`，上传服务器并执行远端部署脚本。

## 1. 服务器目录结构

推荐使用以下目录结构：

```text
/opt/meimei-member/
├── backups/                 # 每次部署前自动备份旧版本
├── deploy/
│   ├── deploy.sh
│   ├── member-api.service
│   └── nginx.conf.example
├── release.tar.gz           # GitHub Actions 上传的发布包
├── releases/                # 每次部署解压目录
├── server/
│   ├── .env                 # 仅服务器保存，不提交 GitHub
│   ├── config.yaml          # GVA 后端配置
│   └── member-api           # Go 后端二进制
└── web/
    └── dist/                # GVA 管理后台静态文件
```

说明：

- `miniprogram/` 这一阶段仍然手动上传体验版，不进入自动部署链路。
- `.env`、数据库密码、微信 `AppSecret` 只保存在服务器，不提交到 GitHub。

## 2. GitHub Secrets 配置

进入仓库 `Settings -> Secrets and variables -> Actions`，新增以下 Secrets：

- `DEPLOY_HOST`: 服务器公网 IP 或域名
- `DEPLOY_PORT`: SSH 端口，例如 `22`
- `DEPLOY_USER`: 执行部署的 SSH 用户
- `DEPLOY_SSH_KEY`: 该用户对应的私钥全文

说明：

- 工作流会使用 `scp` 上传 `release.tar.gz` 与 `deploy/` 下的部署文件。
- 工作流随后会远程执行 `/opt/meimei-member/deploy/deploy.sh`。

## 3. SSH Key 配置

在本地生成部署专用密钥：

```bash
ssh-keygen -t ed25519 -C "github-actions@meimei-member"
```

把公钥写入服务器目标用户的 `~/.ssh/authorized_keys`：

```bash
mkdir -p ~/.ssh
chmod 700 ~/.ssh
cat id_ed25519.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

把私钥内容完整复制到 GitHub Secret `DEPLOY_SSH_KEY`。

## 4. systemd 配置

安装 systemd 服务文件：

```bash
sudo cp /opt/meimei-member/deploy/member-api.service /etc/systemd/system/member-api.service
sudo systemctl daemon-reload
sudo systemctl enable member-api
```

在 `/opt/meimei-member/server/.env` 中准备运行时环境变量，例如：

```bash
GIN_MODE=release
GVA_CONFIG=/opt/meimei-member/server/config.yaml
```

同时准备后端配置文件：

```bash
cp server/config.yaml /opt/meimei-member/server/config.yaml
```

## 5. Nginx 配置

安装 Nginx 配置：

```bash
sudo cp /opt/meimei-member/deploy/nginx.conf.example /etc/nginx/conf.d/meimei-member.conf
sudo nginx -t
sudo systemctl reload nginx
```

配置说明：

- `root` 指向 `/opt/meimei-member/web/dist`
- `/api/` 请求会先去掉 `/api` 前缀，再反代到本地 `127.0.0.1:8888`
- 这样可以兼容 GVA 默认 `router-prefix: ""` 的原始后端配置
- `try_files` 已支持 Vue history fallback

如果你有 HTTPS 或域名证书，可以在此配置基础上继续扩展。

## 6. 首次部署步骤

首次部署建议按下面顺序执行：

1. 创建服务器目录：

```bash
sudo mkdir -p /opt/meimei-member/{deploy,server,web,backups,releases}
sudo chown -R <deploy-user>:<deploy-user> /opt/meimei-member
```

2. 手动上传以下文件到服务器：

- `deploy/deploy.sh`
- `deploy/member-api.service`
- `deploy/nginx.conf.example`
- `server/config.yaml`

3. 创建 `/opt/meimei-member/server/.env`，写入运行时环境变量。

4. 安装并启用 systemd 服务：

```bash
sudo cp /opt/meimei-member/deploy/member-api.service /etc/systemd/system/member-api.service
sudo systemctl daemon-reload
sudo systemctl enable member-api
```

5. 安装并启用 Nginx 配置：

```bash
sudo cp /opt/meimei-member/deploy/nginx.conf.example /etc/nginx/conf.d/meimei-member.conf
sudo nginx -t
sudo systemctl reload nginx
```

6. 在 GitHub 配置好四个 Secrets。

7. 本地代码推送到 `main`：

```bash
git push origin main
```

随后 GitHub Actions 会自动：

- 构建 `server/member-api`
- 构建 `web/dist`
- 打包 `release.tar.gz`
- 上传到 `/opt/meimei-member`
- 远程执行 `/opt/meimei-member/deploy/deploy.sh`

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

2. 选择目标备份，例如 `20260418173000`

3. 回滚后端二进制：

```bash
cp /opt/meimei-member/backups/20260418173000/server/member-api /opt/meimei-member/server/member-api
```

4. 回滚前端静态文件：

```bash
rm -rf /opt/meimei-member/web/dist
cp -a /opt/meimei-member/backups/20260418173000/web/dist /opt/meimei-member/web/dist
```

5. 重新加载服务：

```bash
sudo systemctl restart member-api
sudo nginx -t
sudo systemctl reload nginx
```

## 8. 当前边界

当前自动部署方案明确不包含以下内容：

- 不做 webhook 监听扩展
- 不做 Docker 镜像构建
- 不做蓝绿部署
- 不自动上传微信小程序体验版
