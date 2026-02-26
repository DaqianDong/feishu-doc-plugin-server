## 使用

### 环境变量
- 根目录新建`.env`文件，填入以下配置
```
APP_ID=xxx # 仅为示例值，请使用你的应用的 App ID，获取方式：开发者后台 -> 基础信息 -> 凭证与基础信息 -> 应用凭证 -> App ID
APP_SECRET=xxx # 仅为示例值，请使用你的应用的 App Secret，获取方式：开发者后台 -> 基础信息 -> 凭证与基础信息 -> 应用凭证 -> App Secret
OCR_KEY=sk-xxx # 阿里的OCR_API_KEY
PORT=8081 # 服务端口
```

### 快速开始
- 安装依赖`go mod tidy`
- 运行`go run main.go`

### 镜像
- 构建`make build_img version=0.0.1`
- 运行`docker compose up -d`