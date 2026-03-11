# CLAUDE.md

本文件为 Claude Code 提供项目上下文和开发指导。

## 快速导航

- [项目概述](#项目概述) | [技术栈](#技术栈) | [开发工作流](#开发工作流)
- [项目结构](#项目结构) | [环境变量](#环境变量) | [API 端点](#api-端点)

---

## 项目概述

**飞书白板 OCR 服务** - 为飞书云文档检查器插件提供后端支持，实现白板图片的 OCR 文字识别功能。

### 核心功能
- 🔐 **飞书 OAuth 认证** - 用户授权登录
- 📸 **白板图片获取** - 从飞书 API 下载白板截图
- 🤖 **智能 OCR 识别** - 使用阿里云通义千问视觉模型识别图片文字
- 📊 **结构化输出** - 返回 JSON 格式的识别结果

---

## 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25.1 |
| Web 框架 | Gin | 1.11.0 |
| Session 管理 | gin-contrib/sessions | 1.0.4 |
| 环境变量 | godotenv | 1.5.1 |
| 飞书 SDK | larksuite/oapi-sdk-go | v3.5.3 |
| OCR 服务 | 阿里云通义千问视觉 | - |

---

## 开发工作流

```bash
# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go

# 构建二进制文件（Linux）
make build

# 构建并打包 Docker 镜像
make build_img version=0.0.1

# 使用 Docker Compose 启动
docker compose up -d
```

---

## 项目结构

```
fs-doc-plugin-server/
├── main.go                 # 应用入口，路由配置
├── Makefile                # 构建脚本
├── go.mod                  # Go 模块定义
├── controller/             # 控制器层
│   ├── auth.go            # 认证相关
│   ├── index.go           # 首页
│   └── whiteboard.go      # 白板 OCR 处理
├── infra/                  # 基础设施层
│   ├── httpclient/        # HTTP 客户端
│   ├── image/             # 图片处理
│   ├── larkclient/        # 飞书 API 客户端
│   └── ocr/               # OCR 服务集成
├── docker/                 # Docker 配置
└── .env                    # 环境变量（不提交到 Git）
```

---

## 环境变量

在 `.env` 文件中配置以下变量：

```bash
# 飞书应用配置（必需）
APP_ID=cli_xxxxxxxxx           # 飞书应用 ID
APP_SECRET=xxxxxxxxxx          # 飞书应用 Secret

# OCR 服务配置（必需）
OCR_URL=https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions
OCR_MODEL=qwen-vl-ocr-latest   # 通义千问 OCR 模型
OCR_KEY=sk-xxxxxxxx            # 阿里云 API Key

# 服务配置（可选）
PORT=8081                      # 服务监听端口，默认 8081
```

---

## API 端点

| 端点 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/` | GET | 服务健康检查 | 否 |
| `/login` | GET | 飞书 OAuth 登录入口 | 否 |
| `/callback` | GET | OAuth 回调处理 | 否 |
| `/whiteboard` | GET | 白板图片 OCR 识别 | 是（Session） |

### `/whiteboard` 端点

**Query 参数**：
- `documentId` - 飞书文档 ID（必需）
- `recordId` - 白板记录 ID（必需）

**响应格式**：
```json
{
  "标题": "提取的标题",
  "描述": "详细描述内容",
  "类型": "分类",
  "核心概念": "主要概念",
  "关键元素": "元素1, 元素2",
  "重要机制": "机制说明",
  "触发条件": "触发条件",
  "标签": "标签1, 标签2"
}
```

---

## 开发注意事项

### 飞书 API
- 使用 `larksuite/oapi-sdk-go/v3` SDK
- 需要正确配置 APP_ID 和 APP_SECRET
- 图片下载可能受网络环境影响，建议添加超时处理

### OCR 调用
- 当前使用阿里云通义千问视觉模型
- 支持自动清理 Markdown 代码块标记
- 响应时间约 1-5 秒，建议异步处理

### Session 管理
- 使用 Cookie 存储 Session（开发环境）
- 生产环境建议使用 Redis 等持久化存储
- Session 密钥不应硬编码（当前仅为示例）

### CORS 配置
- 当前允许所有来源（`*`）
- 生产环境需指定允许的域名

---

## 常见问题

**Q: OCR 识别返回空结果？**
A: 检查 OCR_KEY 是否有效，确保 OCR_URL 和 OCR_MODEL 配置正确

**Q: 飞书 API 调用失败？**
A: 确认 APP_ID 和 APP_SECRET 正确，检查应用是否已发布

**Q: Session 频繁过期？**
A: Cookie Session 在服务重启后会丢失，生产环境建议使用 Redis

**Q: 构建失败？**
A: 确保 Go 版本 >= 1.25，运行 `go mod tidy` 更新依赖

---

## 快速参考

**关键文件**：
- `main.go:main()` - 应用入口，路由配置
- `controller/whiteboard.go:WhiteboardController()` - OCR 核心逻辑
- `infra/ocr/ocr.go:OCR()` - OCR API 调用

**开发提示**：
- ⚡ Go 1.25+ 支持，使用 Gin 框架
- 🔐 生产环境需更换 Session 存储方案
- 🐳 Docker 部署使用 `docker-compose.yaml`
- 🔄 修改代码后无需重启（使用 `air` 可实现热重载）

---

## 相关链接

- [飞书开放平台](https://open.feishu.cn/)
- [通义千问 API 文档](https://help.aliyun.com/zh/dashscope/)
- [Gin 框架文档](https://gin-gonic.com/)
