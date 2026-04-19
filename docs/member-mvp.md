# 医美会员积分兑换后台 MVP

## 范围

本轮仅实现 GVA 后台 MVP：

- 会员管理
- 积分账户管理
- 积分流水管理
- 积分商品管理
- 积分兑换订单管理

不包含：

- 小程序
- 视频面诊
- 预约
- 充值 / 支付
- 营销活动
- Excel 导入

## 数据表

服务启动时会通过 `server/initialize/gorm_biz.go` 自动迁移以下表：

- `mm_members`
- `mm_point_accounts`
- `mm_point_transactions`
- `mm_point_products`
- `mm_redemption_orders`

## 启动期自动注册

服务启动时会通过 `server/service/member/bootstrap.go` 自动补齐：

- 业务 API 到 `sys_apis`
- 后台菜单到 `sys_base_menus`
- 超级管理员角色 `888` 的菜单关联
- 超级管理员角色 `888` 的按钮级权限
- 超级管理员角色 `888` 的 casbin API 权限

如果当前数据库已经初始化过 GVA，只需要重新启动后端，新的会员积分菜单就会自动出现到超级管理员账号下。

如果数据库里还保留旧版会员积分表 `member_*`，可以执行下面的命令迁移历史数据到新 `mm_*` 表：

```bash
cd server
go run ./cmd/member_migrate -c config.local.yaml
```

迁移命令只会复制旧表数据到新表，并补齐会员模块 API / 菜单 / 按钮权限，不会删除旧表。

## 运行方式

后端：

```bash
cd server
go mod tidy
go build ./...
go run main.go
```

前端：

```bash
cd web
npm install
npm run dev
```

生产构建验证：

```bash
cd web
npm run build
```
