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
  - `cd server && go test ./service/member/...` 通过
  - `cd web && npm run build` 通过
  - `cd server && go run ./cmd/member_migrate -c config.local.yaml` 已执行，旧 `member_*` 数据已迁移到当前 MySQL 库的 `mm_*` 新表

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
- 补充 `server/service/member/member_service_test.go` 的 MySQL 集成测试：
  - 会员创建自动开户
  - 积分账户缺失后的恢复
  - 手工加减积分流水
  - 兑换订单创建 / 完成 / 取消
  - 禁删有积分流水会员
  - 禁用会员与库存不足场景校验
- 补充 `server/service/member/bootstrap_test.go` 的 MySQL 集成测试：
  - 会员模块 API / 菜单 / 菜单按钮自动初始化
  - 超级管理员菜单挂载
  - 超级管理员按钮级权限与 Casbin 规则自动补齐
- 新增旧表迁移能力：
  - `server/service/member/migration.go` 提供旧 `member_*` 表到新 `mm_*` 表的数据迁移
  - `server/service/member/migration_test.go` 覆盖会员 / 账户 / 商品 / 订单 / 流水的字段映射
  - `server/cmd/member_migrate/main.go` 提供手动执行命令入口
- 修复积分账户被软删除后无法恢复的问题，避免 `mm_point_accounts.member_id` 唯一索引与软删除冲突。

### 前端

- 重写 `web/src/api/member.js`，改成新的会员 / 积分账户 / 积分流水 / 积分商品 / 兑换订单接口。
- 对齐以下页面字段与接口：
  - `web/src/view/member/member/index.vue`
  - `web/src/view/member/account/index.vue`
  - `web/src/view/member/log/index.vue`
  - `web/src/view/member/goods/index.vue`
  - `web/src/view/member/order/index.vue`
- 为以下弹窗补充表单校验：
  - 会员新增 / 编辑：姓名、手机号、等级、状态、备注长度
  - 积分商品新增 / 编辑：商品名称、积分价格、库存、状态、排序、封面 URL
  - 积分账户手工调整：积分数量、调整备注
  - 兑换订单创建：会员、商品、数量、收货人、联系电话、备注长度
- 为会员模块页面补齐按钮级权限联动：
  - 会员管理：新增、编辑、状态切换、查看积分账户、删除
  - 积分账户：手工加分、手工扣分
  - 积分商品：新增、编辑、上下架、删除
  - 兑换订单：新建、查看详情、完成、取消
- 更新会员模块初始化逻辑，启动时自动写入对应菜单按钮，并为超级管理员补齐按钮授权。

## 仍未完成

- 还没有补接口层、前端层或端到端自动化测试，当前主要覆盖会员服务层 MySQL 集成测试。
- 没有做小程序端或其他业务扩展。

## 已知问题

- 工作区里仍有一个与本轮任务无关的旧差异：`deploy/nginx.conf.example` 存在尾随空格问题。
- 这不影响当前后端或前端构建，但如果后续要追求 `git diff --check` 全绿，需要单独清理那份部署文件。
- `npm run build` 有大 chunk warning，但构建已成功，不阻塞当前 MVP。
- 会员服务测试默认会读取 `server/config.local.yaml` 的 MySQL 配置，并自动创建 / 使用 `<db-name>_codex_test` 测试库；如需覆盖，可设置 `MEMBER_TEST_DSN` 或 `MEMBER_TEST_DB_NAME`。
- 生产 / 联调库里的旧 `member_*` 表目前仍然保留，迁移命令只做复制到 `mm_*` 新表，不会删除旧表；如后续确认稳定，可再单独规划旧表下线。

## 下次继续建议

如果下一次继续开发，可以直接用下面这句作为提示词：

> 请先阅读 `docs/codex-progress.md` 和 `git diff`，基于当前仓库继续完善医美会员积分兑换系统后台 MVP。优先补接口层 / 前端层自动化测试，并在确认新表稳定后评估旧 `member_*` 表的下线方案，同时保持 `cd server && go build ./...`、`cd server && go test ./service/member/...` 与 `cd web && npm run build` 通过。不要处理部署、GitHub Actions、小程序和支付。
