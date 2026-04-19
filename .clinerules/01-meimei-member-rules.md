# 医美会员积分兑换系统开发规则

## 项目背景

当前项目是 DBkevin/meimei-member，基于 Gin-Vue-Admin 开发。

视频面诊项目已经封存，不再开发。

当前项目目标是：医美会员积分兑换系统。

当前已经完成第一版后台 MVP，并且服务器上已经通过主流程验收：

- 后台可以打开
- 后端可以启动
- 数据库可以写入
- 可以新建会员
- 可以自动创建积分账户
- 可以手动增加积分
- 可以查看积分流水
- 可以新建积分商品
- 可以创建兑换订单
- 可以扣减会员积分
- 可以扣减商品库存
- 可以核销订单
- 可以取消/退款订单

## 当前开发优先级

当前只允许优先处理：

1. 阶段 0：基线健康修正
2. 阶段 1：后端业务保护

暂时不要做：

- 小程序
- 视频面诊
- 预约系统
- 支付
- 微信登录
- 短信
- dashboard
- 导入导出
- GitHub Actions
- 部署脚本
- Nginx
- MCP
- 复杂会员等级权益
- 优惠券系统
- 营销活动

## 技术规则

必须遵守 GVA 分层结构：

- model：数据模型
- request：请求结构
- service：业务逻辑
- api：参数绑定、调用 service、返回 response
- router：路由注册
- initialize / enter.go：初始化注册

API 层不能直接写业务逻辑。

Service 层不能引用 gin.Context。

业务规则必须收口在 service 层。

不要修改 go.mod module。

不要重构 GVA 原有系统。

不要删除已有会员积分 MVP 代码。

不要提交真实数据库密码、服务器 IP、JWT key、域名、证书路径。

不要修改：

- .github/workflows
- deploy/
- Nginx 配置
- 生产 config.yaml
- 生产 .env

## 后端验证命令

每次修改后必须执行：

```bash
cd server
gofmt -w ./api ./model ./request ./router ./service ./initialize
go test ./...
go build ./...