# catiSip

对标南康功能的，基于 SIP 的系统实现，集成 FreeSWITCH，前端使用 React，后端使用 Go 编写业务逻辑。

A SIP-based system implementation similar to Nankang functionality, integrating FreeSWITCH, with a React frontend and Go backend for business logic.

## 系统架构 / System Architecture

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  React Frontend │ ───> │   Go Backend    │ ───> │   FreeSWITCH    │
│   (Port 3000)   │ HTTP │   (Port 8080)   │ ESL  │   (Port 5060)   │
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

## 功能特性 / Features

- ✅ SIP 呼叫管理 / SIP Call Management
- ✅ 实时呼叫状态监控 / Real-time Call Status Monitoring
- ✅ 分机管理 / Extension Management
- ✅ FreeSWITCH 集成 / FreeSWITCH Integration
- ✅ RESTful API
- ✅ 现代化 React UI / Modern React UI

## 技术栈 / Tech Stack

### 后端 / Backend
- Go 1.21+
- net/http (HTTP server)
- 自定义 SIP 客户端 / Custom SIP client

### 前端 / Frontend
- React 18
- Modern JavaScript (ES6+)
- CSS3

### 电话系统 / Telephony
- FreeSWITCH
- SIP Protocol

## 快速开始 / Quick Start

### 前置要求 / Prerequisites

- Go 1.21 或更高版本 / Go 1.21 or higher
- Node.js 16+ and npm
- FreeSWITCH (可选，用于完整功能 / Optional, for full functionality)

### 安装步骤 / Installation

#### 1. 克隆仓库 / Clone the repository

```bash
git clone https://github.com/phughe11/catiSip.git
cd catiSip
```

#### 2. 启动后端 / Start Backend

```bash
cd backend

# 配置环境变量 / Configure environment variables (optional)
export SIP_HOST=localhost
export SIP_USERNAME=1000
export SIP_PASSWORD=1234
export SIP_DOMAIN=localhost
export PORT=8080

# 构建并运行 / Build and run
go build -o catiSip
./catiSip
```

后端服务将在 http://localhost:8080 启动

Backend service will start at http://localhost:8080

#### 3. 启动前端 / Start Frontend

```bash
cd frontend

# 安装依赖 / Install dependencies
npm install

# 启动开发服务器 / Start development server
npm start
```

前端应用将在 http://localhost:3000 启动

Frontend application will start at http://localhost:3000

## API 接口 / API Endpoints

### Health Check
```
GET /api/health
```

### 发起呼叫 / Make a Call
```
POST /api/call/make
Content-Type: application/json

{
  "from": "1000",
  "to": "1001"
}
```

### 挂断呼叫 / Hangup Call
```
POST /api/call/hangup
Content-Type: application/json

{
  "call_id": "call-1234567890"
}
```

### 查询呼叫状态 / Get Call Status
```
GET /api/call/status?call_id=call-1234567890
```

### 获取分机列表 / List Extensions
```
GET /api/extensions
```

## 配置 / Configuration

### 环境变量 / Environment Variables

| 变量名 / Variable | 描述 / Description | 默认值 / Default |
|------------------|-------------------|-----------------|
| `SIP_HOST` | FreeSWITCH 服务器地址 / FreeSWITCH server host | `localhost` |
| `SIP_PORT` | SIP 端口 / SIP port | `5060` |
| `SIP_USERNAME` | SIP 用户名 / SIP username | `1000` |
| `SIP_PASSWORD` | SIP 密码 / SIP password | `1234` |
| `SIP_DOMAIN` | SIP 域 / SIP domain | `localhost` |
| `PORT` | HTTP 服务端口 / HTTP server port | `8080` |

### 配置文件 / Configuration File

也可以使用 JSON 配置文件：

You can also use a JSON configuration file:

```json
{
  "sip": {
    "host": "localhost",
    "port": 5060,
    "username": "1000",
    "password": "1234",
    "domain": "localhost"
  },
  "server": {
    "port": 8080
  }
}
```

使用配置文件：

Use the configuration file:

```bash
export CONFIG_FILE=/path/to/config.json
./catiSip
```

## FreeSWITCH 集成 / FreeSWITCH Integration

本系统设计用于与 FreeSWITCH 集成。完整功能需要：

This system is designed to integrate with FreeSWITCH. For full functionality:

1. 安装并配置 FreeSWITCH / Install and configure FreeSWITCH
2. 配置 SIP 分机 / Configure SIP extensions
3. 启用 ESL (Event Socket Library) / Enable ESL
4. 更新环境变量以指向 FreeSWITCH 服务器 / Update environment variables to point to FreeSWITCH server

目前的实现包含一个模拟的 SIP 客户端用于演示。在生产环境中，应该使用真实的 FreeSWITCH ESL 库。

The current implementation includes a simulated SIP client for demonstration. In production, use a real FreeSWITCH ESL library.

## 项目结构 / Project Structure

```
catiSip/
├── backend/              # Go 后端 / Go backend
│   ├── main.go          # 主入口 / Main entry point
│   ├── config/          # 配置管理 / Configuration management
│   ├── handlers/        # HTTP 处理器 / HTTP handlers
│   └── sip/             # SIP 客户端 / SIP client
├── frontend/            # React 前端 / React frontend
│   ├── public/          # 静态资源 / Static assets
│   └── src/             # 源代码 / Source code
│       ├── App.js       # 主应用组件 / Main app component
│       └── App.css      # 样式 / Styles
├── docs/                # 文档 / Documentation
└── README.md            # 本文件 / This file
```

## 开发 / Development

### 后端开发 / Backend Development

```bash
cd backend

# 运行测试 / Run tests
go test ./...

# 格式化代码 / Format code
go fmt ./...

# 运行代码检查 / Run linter
go vet ./...
```

### 前端开发 / Frontend Development

```bash
cd frontend

# 启动开发服务器 / Start development server
npm start

# 构建生产版本 / Build for production
npm run build

# 运行测试 / Run tests
npm test
```

## 部署 / Deployment

### 后端部署 / Backend Deployment

```bash
cd backend
go build -o catiSip
./catiSip
```

### 前端部署 / Frontend Deployment

```bash
cd frontend
npm run build
# 将 build/ 目录部署到静态文件服务器
# Deploy the build/ directory to a static file server
```

## 贡献 / Contributing

欢迎提交 Pull Request 和 Issue！

Pull requests and issues are welcome!

## 许可证 / License

MIT License

## 联系方式 / Contact

- GitHub: [phughe11/catiSip](https://github.com/phughe11/catiSip)

## 致谢 / Acknowledgments

- FreeSWITCH 社区 / FreeSWITCH Community
- React 团队 / React Team
- Go 社区 / Go Community

