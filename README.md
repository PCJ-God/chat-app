# Chat App · 全栈即时聊天

一个功能完整的即时聊天应用：**Go + Gin** 后端、**React + TypeScript** 前端，支持私聊、群聊、实时推送、在线状态、通知中心与未读数，配套 Docker 一键部署。

> 截图待补充：可运行后把界面截图放到 `docs/` 目录，并在此处引用。

---

## 功能特性

- **账号体系**：注册 / 登录（JWT），密码 bcrypt 哈希。
- **私聊**：会话列表、消息气泡、Enter 发送 / Shift+Enter 换行、历史消息分页加载。
- **群聊**：创建群聊、按房间 ID 加入、群内消息（显示发送者昵称）。
- **实时通信**：基于 WebSocket 的消息与在线状态推送（断线指数退避重连）。
- **在线状态**：在线 / 离开 / 离线，会话列表实时显示。
- **通知中心**：未读通知列表、标记已读；顶栏铃铛显示私聊未读数。
- **未读 / 已读**：打开会话自动把收到的未读消息标记为已读。
- **Docker 部署**：前端 nginx 托管 + 反向代理，后端、MySQL、Redis 一体化编排。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.27、Gin、GORM（MySQL）、go-redis、gorilla/websocket、golang-jwt、bcrypt |
| 前端 | React 18、TypeScript、Vite、Ant Design、TanStack Query、Zustand、axios、dayjs、react-router |
| 基础设施 | MySQL 8、Redis 7、Nginx、Docker Compose |

## 架构概览

```
浏览器 (React SPA)
  │  REST /api/*        WebSocket /ws
  ▼                        ▼
Nginx ──反代──► Go 后端 (Gin) ──► MySQL (持久化)
                  │  ▲            Redis (在线状态/会话)
                  │  └─ WebSocket Hub (消息/在线状态广播)
                  ▼
              业务 Service ──► Repository (GORM)
```

- 分层：`handler`（HTTP）→ `service`（业务）→ `repository`（数据访问）。
- 实时：服务端 `websocket.Hub` 维护连接，私聊推给接收方、群聊广播给房间成员。
- 鉴权：JWT 中间件解析 `Authorization: Bearer <token>`，WebSocket 通过 `?token=` 校验。

## 目录结构

```
.
├── cmd/server/main.go        # 入口：依赖装配 + 路由注册
├── internal/
│   ├── handler/              # HTTP 处理器（auth/user/message/room/notification/ws）
│   ├── service/              # 业务逻辑
│   ├── repository/           # 数据访问 (GORM) + 自动迁移
│   └── model/                # 数据模型 (User/Room/Message/Notification/Session)
├── pkg/
│   ├── auth/                 # JWT 签发与中间件
│   ├── config/               # 配置加载 (viper)
│   ├── websocket/            # Hub / Client / Upgrader
│   └── utils/                # 密码哈希等
├── config/                   # 配置（含 config.example.yaml）
├── frontend/                 # React + TS + Vite 前端
│   └── src/
│       ├── api/              # HTTP 接口封装（axios）
│       ├── components/       # 会话列表 / 消息线程 / 群聊
│       ├── layouts/ pages/   # 布局与页面
│       ├── realtime/         # WebSocket 客户端 + useRealtime
│       ├── stores/           # Zustand 状态（auth/presence/ui）
│       └── types/ utils/     # 类型与工具
├── Dockerfile                # 后端镜像（多阶段）
├── frontend/Dockerfile       # 前端镜像（构建 + nginx）
├── frontend/nginx.conf       # SPA 回退 + 反代 /api、/ws
├── docker-compose.yml        # 本地依赖（仅 mysql + redis）
└── docker-compose.full.yml   # 全栈部署（frontend/backend/mysql/redis）
```

## 快速开始（本地开发）

### 1. 准备依赖（MySQL + Redis）

```bash
docker compose up -d          # 使用 docker-compose.yml，默认 root 密码 your_password
```

### 2. 配置后端

```bash
# 复制示例配置
cp config/config.example.yaml config/config.yaml     # Windows: copy config\config.example.yaml config\config.yaml
```

编辑 `config/config.yaml`，确认 `database.password` 与上一步的 `MYSQL_ROOT_PASSWORD` 一致；`jwt.secret` 改成足够长的随机串。

### 3. 启动后端

```bash
go run ./cmd/server            # 监听 http://localhost:8080
```

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev                    # http://localhost:5173（已代理 /api、/ws 到 8080）
```

用两个浏览器（或无痕窗口）分别注册账号即可体验私聊 / 群聊。

## Docker 部署（全栈）

```bash
# 构建镜像（单独构建以避开部分 compose 版本的 bake 问题）
docker build -t chat-app-backend:local  -f Dockerfile .
docker build -t chat-app-frontend:local -f frontend/Dockerfile frontend

# 启动
docker compose -f docker-compose.full.yml up -d --no-build

# 访问：http://127.0.0.1:8099
```

停止 / 清理：

```bash
docker compose -f docker-compose.full.yml down            # 保留数据卷
docker compose -f docker-compose.full.yml down -v         # 连数据卷一起删除
```

> 若你的 `docker compose` 版本能正常用 bake，也可直接 `docker compose -f docker-compose.full.yml up -d --build`。

## 配置说明

- 后端配置见 `config/config.example.yaml`；`config/config.yaml` 为本地真实配置，**已被 `.gitignore` 忽略，不要提交**。
- `config/config.docker.yaml` 供容器内使用（`host: 0.0.0.0`、服务名 `mysql`/`redis`），端口与密码为示例占位，请按需修改。

## API 概览

**认证（无需登录）**

- `POST /api/auth/register` — 注册
- `POST /api/auth/login` — 登录，返回 `{ token, user }`

**需要 `Authorization: Bearer <token>`**

- 用户：`GET /api/users`、`GET /api/users/me`
- 消息：`POST /api/messages`、`GET /api/messages?target_id=|room_id=&page=&page_size=`、`POST /api/messages/:id/read`、`GET /api/messages/unread/count`
- 群聊：`POST /api/rooms`、`GET /api/rooms`、`GET /api/rooms/:id`、`POST /api/rooms/:id/join`、`POST /api/rooms/:id/leave`、`POST /api/rooms/:id/messages`
- 通知：`GET /api/notifications`、`POST /api/notifications/:id/read`

**WebSocket**：`ws://<host>/ws?token=<jwt>`，事件类型：

- `message`：`{ type, id, sender_id, content, created_at, receiver_id? | room_id? }`
- `presence`：`{ type, user_id, status }`

## 常见问题

- **容器内构建拉不到依赖**：镜像构建已清空代理环境变量并使用国内源（Go `goproxy.cn`、npm `registry.npmmirror.com`、apk 阿里云）。不同网络环境可按需调整 `Dockerfile`。
- **`docker compose build` 报 `failed to execute bake`**：部分 compose 版本 bake 不可用，按上文用 `docker build` 逐镜像构建后 `up -d --no-build`。
- **本机 `localhost:8099` 打不开**：优先用 `127.0.0.1:8099`（个别环境 IPv6 转发异常）。

## 许可证

MIT
