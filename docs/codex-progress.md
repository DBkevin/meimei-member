# Codex Progress

## 当前状态

- 项目已基于 GVA 原有分层完成会员积分后台 MVP 的后端改造：
  - `mm_members`
  - `mm_point_accounts`
  - `mm_point_transactions`
  - `mm_point_products`
  - `mm_redemption_orders`
- 后端接口、路由、服务、模型、自动迁移、Swagger 注释已经对齐新 schema。
- 前端会员模块页面和 `web/src/api/member.js` 已对齐新的接口路径和字段命名。
- 当前验证结果：
  - `cd server && go build ./...` 通过
  - `cd web && npm run build` 通过

## 本次完成

### 后端

- 重构 `server/model/member` 下的业务模型与常量，旧 `member_*` 表改为 `mm_*` 表。
- 重写 `server/service/member` 的核心逻辑：
  - 创建会员自动创建积分账户
  - 手工加减积分统一走事务和积分流水
  - 积分商品增删改查和上下架
  - 兑换订单创建、完成、取消
  - 订单创建时在同一事务内扣积分、扣库存、建单、写流水
- 更新 `server/api/v1/member` 和 `server/router/member` 的接口与路由注册。
- 更新 `server/initialize/gorm_biz.go` 与 `server/initialize/router.go`。

### 前端

- 重写 `web/src/api/member.js`，改成新的会员 / 积分账户 / 积分流水 / 积分商品 / 兑换订单接口。
- 对齐以下页面字段与接口：
  - `web/src/view/member/member/index.vue`
  - `web/src/view/member/account/index.vue`
  - `web/src/view/member/log/index.vue`
  - `web/src/view/member/goods/index.vue`
  - `web/src/view/member/order/index.vue`

## 仍未完成

- 没有补自动化测试。
- 没有做旧表到新 `mm_*` 表的数据迁移脚本。
- 没有处理会员模块的更细致表单校验与按钮级权限联动。
- 没有做小程序端或其他业务扩展。

## 已知问题

- 工作区里仍有一个与本轮任务无关的旧差异：`deploy/nginx.conf.example` 存在尾随空格问题。
- 这不影响当前后端或前端构建，但如果后续要追求 `git diff --check` 全绿，需要单独清理那份部署文件。
- `npm run build` 有大 chunk warning，但构建已成功，不阻塞当前 MVP。

## 下次继续建议

如果下一次继续开发，可以直接用下面这句作为提示词：

> 请先阅读 `docs/codex-progress.md` 和 `git diff`，基于当前仓库继续完善医美会员积分兑换系统后台 MVP。优先补自动化测试、细化前端表单校验，并保持 `cd server && go build ./...` 与 `cd web && npm run build` 通过。不要处理部署、GitHub Actions、小程序和支付。
