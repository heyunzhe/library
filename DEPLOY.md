# 图书管理系统 - 部署文档

## 📋 项目概述

基于 Go + MySQL 的图书管理系统，提供图书管理、借阅管理、用户管理、在线公告、在线阅读等功能。

---

## 🚀 快速开始（Docker 部署 - 推荐）

### 前置条件
- Docker & Docker Compose

### 部署步骤

```bash
# 1. 克隆项目
git clone https://github.com/heyunzhe/library.git
cd library

# 2. 配置环境变量（可选，有默认值）
cp .env.example .env
# 编辑 .env 文件中的配置

# 3. 启动服务
docker-compose up -d

# 4. 访问系统
# http://localhost:8080
```

### 默认管理员账号
- **管理员 ID:** `a`
- **密码:** `1`

---

## 🛠 手动部署

### 前置条件

| 依赖 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.21+ | 编译运行后端 |
| MySQL | 8.0+ | 数据存储 |
| Nginx（可选）| 任意 | 反向代理 |

### 1. 数据库配置

```sql
CREATE DATABASE library CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `MYSQL_DSN` | MySQL 连接字符串 | `root:123456@tcp(127.0.0.1:3306)/library?charset=utf8mb4&parseTime=True&loc=Local` |
| `COOKIE_SECRET` | Session 加密密钥 | `default-cookie-secret-key-change-in-production` |
| `JWT_SECRET` | JWT 签名密钥 | `your-256-bit-secret-key-change-in-production!!` |
| `SMTP_HOST` | 邮件服务器地址 | `smtp.qq.com` |
| `SMTP_PORT` | 邮件服务器端口 | `587` |
| `SMTP_USER` | 邮箱账号（需你自己申请） | - |
| `SMTP_PASS` | 邮箱密码/授权码（需你自己申请） | - |

### 3. 编译运行

```bash
# 编译
go build -o library-server .

# 运行（可添加环境变量）
MYSQL_DSN="root:密码@tcp(数据库IP:3306)/library?charset=utf8mb4&parseTime=True&loc=Local" \
COOKIE_SECRET="自定义密钥" \
JWT_SECRET="自定义JWT密钥" \
SMTP_HOST="smtp.qq.com" \
SMTP_USER="your@qq.com" \
SMTP_PASS="邮箱授权码" \
./library-server
```

### 4. 访问
- 浏览器打开 `http://localhost:8080`

---

## 🐳 Docker Compose 配置详解

`docker-compose.yml` 已配置好以下服务：

| 服务 | 端口 | 说明 |
|------|------|------|
| `backend` | `8080` | Go 应用 |
| `mysql` | `3306` | 数据库 |

```

## 📝 配置文件

### 环境变量 (.env.example)

```env
# MySQL 数据库
MYSQL_DSN=root:123456@tcp(mysql:3306)/library?charset=utf8mb4&parseTime=True&loc=Local

# Session 密钥
COOKIE_SECRET=your-cookie-secret-key

# JWT 密钥
JWT_SECRET=your-jwt-secret-key

# 邮箱配置（用于发送验证码）
SMTP_HOST=smtp.qq.com
SMTP_PORT=587
SMTP_USER=your_email@qq.com
SMTP_PASS=your_email_auth_code
```

---

## 📁 项目结构

```
library/
├── main.go              # 入口文件，路由注册
├── go/                  # 后端处理器
│   ├── adminbooks.go    # 图书管理（增删改查）
│   ├── adminnotices.go  # 公告管理
│   ├── email.go         # 邮件发送
│   ├── jwt.go           # JWT 认证
│   ├── login.go         # 登录/登出/中间件
│   ├── readbook.go      # 在线阅读
│   ├── register.go      # 用户注册/验证码
│   ├── table.go         # 表结构初始化
│   ├── userbook.go      # 借书/还书/排行榜
│   └── users.go         # 用户管理
├── html/                # 前端页面模板
├── css/                 # 样式文件
├── js/                  # 前端 JavaScript
├── images/              # 图书封面图片
├── font/                # 图标字体
├── Dockerfile           # Docker 构建文件
└── docker-compose.yml   # Docker Compose 配置
```

---

## ⚠️ 注意事项

### 1. 首次部署
- 程序启动时会自动创建数据库表，**无需手动建表**
- 默认管理员账号 **`a` / `1`** 已自动插入，开箱即用

### 2. 邮箱配置（可选，默认不配置）
- 注册验证码、密码重置需要配置 SMTP，**需使用你自己的邮箱账号**
- QQ邮箱需开启 SMTP 服务并使用 **授权码**（不是QQ密码），参考：[QQ邮箱帮助中心](https://help.mail.qq.com/)
- 如果不配置邮箱，管理员可以直接在数据库中添加用户

### 3. MySQL 模式
- 项目依赖 MySQL 8.0+
- 如果遇到 `ONLY_FULL_GROUP_BY` 错误，可以调整 SQL 模式：
  ```sql
  SET GLOBAL sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
  ```

### 4. 生产环境建议
- 修改默认的 `COOKIE_SECRET` 和 `JWT_SECRET`
- 修改默认管理员密码
- 使用 HTTPS
- MySQL 不要使用 root 账户，创建一个专用数据库用户

---

## 🔧 常见问题

**Q: 启动后访问页面报 500？**
- 检查 MySQL 是否正常运行
- 检查 `MYSQL_DSN` 配置是否正确

**Q: 发送验证码失败？**
- 检查 SMTP 配置是否正确
- QQ邮箱需要使用授权码，不是密码

**Q: 如何添加图书？**
- 使用管理员账号登录后台 → 图书管理 → 添加图书
- 需要上传封面图片
