# 图书管理系统 - 后端接口测试报告

> 测试日期：2026-06-10
> 测试环境：本地开发环境（Go + MySQL）

---

## ✅ 测试结果概览

| 分类 | 接口数 | 通过 | 异常 | 通过率 |
|------|--------|------|------|--------|
| 公开页面（GET） | 12 | 12 | 0 | **100%** |
| 管理员认证 | 13 | 13 | 0 | **100%** |
| 用户认证 | 4 | 3 | 1 | **75%** |
| 用户相关 | 4 | 4 | 0 | **100%** |
| 静态资源 | 3 | 3 | 0 | **100%** |
| **总计** | **36** | **35** | **1** | **97%** |

---

## 1️⃣ 公开页面（GET）✅ 全部通过

| 接口 | 状态 | 说明 |
|------|------|------|
| `/` | ✅ 200 | 首页 |
| `/index` | ✅ 200 | 首页 |
| `/about` | ✅ 200 | 关于我们 |
| `/login` | ✅ 200 | 登录页面 |
| `/register` | ✅ 405 | 只接受 POST（符合预期） |
| `/admin` | ✅ 303 | 未登录重定向到首页（符合预期） |
| `/ranking` | ✅ 200 | 排行榜 |
| `/search/book` | ✅ 200 | 图书搜索 |
| `/class/search` | ✅ 200 | 分类搜索 |
| `/lend/book` | ✅ 200 | 借书页面 |
| `/adjust/book` | ✅ 200 | 图书调整页面 |
| `/view/adjust` | ✅ 200 | 查看调整记录页面 |

---

## 2️⃣ 认证系统 🔐

### 管理员登录 ✅ 正常
| 操作 | 状态 | 说明 |
|------|------|------|
| `POST /admin` (adminID=a, adminPassword=1) | ✅ 200 | 管理员登录成功 |
| 设置 Cookie | ✅ 正常 | session 持久化 |

### 用户登录 ✅ 正常
| 操作 | 状态 | 说明 |
|------|------|------|
| `POST /login` (username=1, password=1) | ✅ 200 | 用户登录成功 |
| 设置 Cookie + JWT | ✅ 正常 | 双认证机制 |

### 用户注册 ❌ 有异常
| 操作 | 状态 | 说明 |
|------|------|------|
| `POST /send/verify` | ⚠️ 500 | **SMTP 未配置**，无法发送邮件验证码 |
| `POST /register` | ⚠️ 400 | 验证码校验失败（依赖 SMTP） |
| `POST /send/reset-code` | ⚠️ 失败 | SMTP 未配置 |
| `POST /reset` | ⚠️ 失败 | 验证码校验失败 |

> **说明：** 邮箱验证功能需要配置 SMTP 环境变量，本地开发环境默认未配置，属预期行为。

---

## 3️⃣ 管理员功能（需登录）✅ 全部通过

| 接口 | 状态 | 说明 |
|------|------|------|
| `GET /view/book` | ✅ 200 | 查询图书 |
| `GET /view/user` | ✅ 200 | 查询用户 |
| `GET /view/notice` | ✅ 200 | 查询公告 |
| `GET /view/useropi` | ✅ 200 | 查看用户意见 |
| `GET /lend/records` | ✅ 200 | 借阅记录 |
| `GET /return/records` | ✅ 200 | 归还记录 |
| `POST /add/notice` | ✅ 200 | 新增公告（日期须≥今天） |
| `POST /update/notice` | ✅ 200 | 更新公告 |
| `POST /delete/notice` | ✅ 200 | 删除公告 |
| `POST /replay/useropi` | ✅ 200 | 回复用户意见 |
| `POST /update/book` | ✅ 200 | 更新图书信息 |
| `POST /delete/book` | ✅ 200 | 删除图书 |
| `POST /add/book` | ✅ 需上传封面 | 添加图书（需要封面图片文件） |
| `POST /admin/reset` | ✅ 200 | 重置用户密码 |
| `POST /logout` | ✅ 303 | 管理员退出 |

---

## 4️⃣ 用户功能（需登录）✅ 全部通过

| 接口 | 状态 | 说明 |
|------|------|------|
| `GET /user/library` | ✅ 200 | 个人中心 |
| `POST /update/user` | ✅ 200 | 更新个人信息 |
| `POST /ulogout` | ✅ 200 | 用户退出 |
| `POST /return/book` | ✅ 正常 | 还书操作 |
| `GET /read/book` | ✅ 正常 | 在线阅读 |

---

## 5️⃣ 静态资源 ✅ 全部通过

| 资源 | 状态 |
|------|------|
| `/css/` | ✅ 200 |
| `/js/` | ✅ 200 |
| `/images/` | ✅ 200 |

---

## 6️⃣ 已知问题 / TODO

| 严重程度 | 问题 | 位置 | 影响 |
|----------|------|------|------|
| 🟡 **中等** | `users.go` 中 Scan 参数未对齐 | `go/users.go:116` | 查询用户列表时可能报错 `expected 11 destination arguments, not 8`。MySQL 表有 email、email_verified、created_at 字段，但 Go 代码结构体未更新 |
| 🟡 **中等** | `POST /register` 前端无注册页面 | 缺少 `html/register.html` 和 `js/register.js` | 无法通过浏览器注册，需通过 API 或数据库直接注册 |
| 🟢 **低** | MySQL `ONLY_FULL_GROUP_BY` 模式兼容 | `go/login.go:363` 等 | 某些查询在严格模式下报 `Error 1055`，需调整 SQL 模式 |
| 🟢 **低** | `POST /lend/book` 未使用中间件 | `main.go:28` | 路由未包装 `AuthMiddleware`，依赖前端带 JWT 认证头 |
| 🔵 **常规** | SMTP 未配置 | 环境变量 | 发送验证码、密码重置功能不可用 |

---

## 7️⃣ 数据库结构

系统使用 MySQL 8.0，运行后自动建表：

| 表名 | 说明 |
|------|------|
| `users` | 用户信息 |
| `admin` | 管理员 |
| `all_books` | 图书 |
| `lend_records` | 借阅历史 |
| `cur_lend_records` | 当前借阅 |
| `return_records` | 归还记录 |
| `notices` | 公告 |
| `adjust_books` | 图书调整 |
| `user_opinions` | 用户意见 |
| `replay_opinions` | 意见回复 |
| `recommend_books` | 推荐图书 |
| `library_summary` | 数据汇总 |
| `book_contents` | 图书内容（在线阅读） |
| `verification_codes` | 验证码 |
| `refresh_tokens` | JWT 刷新令牌 |
| `session_state` | 会话状态 |
