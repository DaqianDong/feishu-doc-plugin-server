# 飞书云文档检查器插件

用于统计和分析飞书文档内容的浏览器插件，提供字数统计、块类型分类、表格统计等功能。

## 目录

- [功能特性](#功能特性) | [技术栈](#技术栈) | [快速开始](#快速开始)
- [后端服务](#后端服务) | [项目结构](#项目结构) | [开发指南](#开发指南)
- [常见问题](#常见问题) | [贡献指南](#贡献指南)

---

## 功能特性

- 📊 **文档统计**：字数统计与 Token 估算
- 🏷️ **智能分类**：块类型自动分类（友好块/问题块）
- 📈 **表格分析**：表格统计与快速定位
- 🔄 **实时更新**：文档变化时自动刷新统计数据

---

## 技术栈

### 前端插件
- React 18.2 + TypeScript 4.9
- Webpack 5 + Babel
- Less + CSS Modules
- @lark-opdev/block-docs-addon-api

### 后端服务（可选）
- Go 语言
- 阿里云 OCR API（图片识别）

---

## 快速开始

### 前置要求
- Node.js >= 14
- npm 或 yarn
- 飞书开发者账号

### 安装依赖
```bash
npm install
```

### 配置环境变量
```bash
# 在飞书开发者后台获取
APP_ID=your_app_id
APP_SECRET=your_app_secret
```

### 开发模式
```bash
npm start              # 启动开发服务器（支持 HMR）
npm run build          # 生产构建到 dist/ 目录
npm run upload         # 构建并上传到飞书平台
```

---

## 后端服务

如需使用图片识别等功能，可启动后端服务：

### 环境变量
```bash
APP_ID=xxx                                    # 飞书应用 ID
APP_SECRET=xxx                                # 飞书应用 Secret
OCR_URL=https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions
OCR_MODEL=qwen-vl-ocr-latest                  # 阿里 OCR 模型
OCR_KEY=sk-xxx                                # 阿里 OCR API Key
PORT=8081                                     # 服务端口
```

### 运行
```bash
go mod tidy           # 安装 Go 依赖
go run main.go        # 运行服务
```

### Docker 部署
```bash
make build_img version=0.0.1    # 构建镜像
docker compose up -d             # 启动服务
```

---

## 项目结构

```
.
├── src/                    # 前端源码
│   ├── App.tsx            # 主应用组件
│   ├── index.css          # 全局样式
│   ├── index.html         # HTML 模板
│   └── index.tsx          # 应用入口
├── dist/                   # 构建产物（自动生成）
├── package.json           # 项目配置
├── webpack.config.js      # Webpack 配置
├── tsconfig.json          # TypeScript 配置
└── .prettierrc.js         # Prettier 配置
```

---

## 开发指南

### 代码规范
- 使用函数组件 + Hooks
- Props 必须定义 TypeScript 接口
- 禁止使用 `any` 类型
- 遵循 Prettier 代码格式化规则

### 版本管理
遵循语义化版本规范（SemVer）：
- **Major**：不兼容的 API 修改
- **Minor**：向下兼容的功能性新增
- **Patch**：向下兼容的问题修正

### 飞书 API 使用

插件将飞书块类型分为两类：

**友好块**（适合 AI 处理）：
文本、标题（1-5级）、列表、引用、待办、代码块、分割线、图片、表格、文件

**问题块**（可能影响 AI 处理）：
标题（6-9级）、内嵌网页、云文档小组件、任务、OKR、白板、议程、AI 模板、降级块

---

## 常见问题

**Q: 统计结果不更新？**
A: 检查文档是否已保存，插件会在文档更新时自动重新统计

**Q: 某些块类型显示"未知"？**
A: 可能是飞书新增的块类型，需要在 `App.tsx` 的 `getBlockTypeName` 函数中添加中文名称

**Q: 开发服务器启动失败？**
A: 检查端口是否被占用，确保 Node.js 版本 >= 14

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

---

## 许可证

MIT License

---

## 相关链接

- [飞书开放平台](https://open.feishu.cn/)
- [飞书插件开发文档](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM)
