# 部署指南 / Deployment Guide

## 部署方式 / Deployment Methods

本指南介绍三种部署方式：

This guide covers three deployment methods:

1. Docker Compose（推荐 / Recommended）
2. 独立部署 / Standalone Deployment
3. Kubernetes 部署 / Kubernetes Deployment

---

## 方式 1: Docker Compose 部署 / Docker Compose Deployment

### 前置要求 / Prerequisites

- Docker 20.10+
- Docker Compose 2.0+

### 步骤 / Steps

1. **克隆仓库 / Clone repository**

```bash
git clone https://github.com/phughe11/catiSip.git
cd catiSip
```

2. **启动所有服务 / Start all services**

```bash
docker-compose up -d
```

这将启动：
- FreeSWITCH (端口 5060, 8021)
- 后端 API (端口 8080)
- 前端 Web (端口 3000)

This will start:
- FreeSWITCH (ports 5060, 8021)
- Backend API (port 8080)
- Frontend Web (port 3000)

3. **查看日志 / View logs**

```bash
docker-compose logs -f
```

4. **停止服务 / Stop services**

```bash
docker-compose down
```

### 访问服务 / Access Services

- 前端 / Frontend: http://localhost:3000
- 后端 API / Backend API: http://localhost:8080
- FreeSWITCH ESL: localhost:8021

---

## 方式 2: 独立部署 / Standalone Deployment

### 后端部署 / Backend Deployment

#### 1. 编译 / Build

```bash
cd backend
go build -o catiSip
```

#### 2. 配置环境变量 / Configure environment

```bash
export SIP_HOST=your-freeswitch-host
export SIP_PORT=5060
export SIP_USERNAME=1000
export SIP_PASSWORD=your-password
export SIP_DOMAIN=your-domain
export PORT=8080
```

或使用配置文件 / Or use config file:

```bash
cp config.example.json config.json
# 编辑 config.json / Edit config.json
export CONFIG_FILE=/path/to/config.json
```

#### 3. 运行 / Run

```bash
./catiSip
```

#### 4. 使用 systemd 管理（推荐）/ Manage with systemd (recommended)

创建服务文件 / Create service file: `/etc/systemd/system/catisip.service`

```ini
[Unit]
Description=CatiSip Backend Service
After=network.target

[Service]
Type=simple
User=catisip
WorkingDirectory=/opt/catisip
Environment="SIP_HOST=localhost"
Environment="SIP_PORT=5060"
Environment="PORT=8080"
ExecStart=/opt/catisip/catiSip
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务 / Start service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable catisip
sudo systemctl start catisip
sudo systemctl status catisip
```

### 前端部署 / Frontend Deployment

#### 1. 构建 / Build

```bash
cd frontend
npm install
npm run build
```

#### 2. 使用 Nginx 部署 / Deploy with Nginx

安装 Nginx / Install Nginx:

```bash
sudo apt-get install nginx
```

配置 Nginx / Configure Nginx: `/etc/nginx/sites-available/catisip`

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /var/www/catisip;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理 / API reverse proxy
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

复制构建文件并启用站点 / Copy build files and enable site:

```bash
sudo cp -r build/* /var/www/catisip/
sudo ln -s /etc/nginx/sites-available/catisip /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 方式 3: Kubernetes 部署 / Kubernetes Deployment

### 创建部署文件 / Create deployment files

#### backend-deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catisip-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: catisip-backend
  template:
    metadata:
      labels:
        app: catisip-backend
    spec:
      containers:
      - name: backend
        image: your-registry/catisip-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: SIP_HOST
          value: "freeswitch-service"
        - name: PORT
          value: "8080"
---
apiVersion: v1
kind: Service
metadata:
  name: catisip-backend-service
spec:
  selector:
    app: catisip-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

#### frontend-deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catisip-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: catisip-frontend
  template:
    metadata:
      labels:
        app: catisip-frontend
    spec:
      containers:
      - name: frontend
        image: your-registry/catisip-frontend:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: catisip-frontend-service
spec:
  selector:
    app: catisip-frontend
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
```

### 部署到 Kubernetes / Deploy to Kubernetes

```bash
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml
```

---

## 生产环境建议 / Production Recommendations

### 安全性 / Security

1. **使用 HTTPS**
   - 配置 SSL/TLS 证书
   - 使用 Let's Encrypt 获取免费证书

2. **防火墙配置**
   ```bash
   # 允许必要的端口
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw allow 5060/udp
   ```

3. **更改默认密码**
   - FreeSWITCH ESL 密码
   - SIP 账户密码

### 性能优化 / Performance

1. **启用 Gzip 压缩**（Nginx）
   ```nginx
   gzip on;
   gzip_types text/plain text/css application/json application/javascript;
   ```

2. **配置缓存**
   ```nginx
   location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
       expires 1y;
       add_header Cache-Control "public, immutable";
   }
   ```

3. **Go 后端优化**
   - 设置 GOMAXPROCS
   - 使用连接池

### 监控和日志 / Monitoring and Logging

1. **日志收集**
   - 使用 ELK Stack (Elasticsearch, Logstash, Kibana)
   - 或使用 Grafana Loki

2. **性能监控**
   - Prometheus + Grafana
   - 监控 CPU、内存、网络

3. **告警配置**
   - 配置告警规则
   - 集成 Slack/Email 通知

### 备份和恢复 / Backup and Recovery

1. **数据库备份**（如果使用）
   ```bash
   # 定时备份脚本
   0 2 * * * /usr/local/bin/backup.sh
   ```

2. **配置备份**
   - FreeSWITCH 配置文件
   - 应用配置文件

---

## 故障排除 / Troubleshooting

### 后端无法启动 / Backend Won't Start

```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 查看日志
journalctl -u catisip -f
```

### 前端无法访问 / Frontend Not Accessible

```bash
# 检查 Nginx 状态
sudo systemctl status nginx

# 查看 Nginx 日志
sudo tail -f /var/log/nginx/error.log
```

### 无法连接 FreeSWITCH / Can't Connect to FreeSWITCH

```bash
# 检查 FreeSWITCH 状态
fs_cli -x "status"

# 检查防火墙
sudo ufw status
```

---

## 更新部署 / Update Deployment

### Docker Compose

```bash
git pull
docker-compose down
docker-compose build
docker-compose up -d
```

### 独立部署 / Standalone

```bash
# 后端 / Backend
systemctl stop catisip
cd backend && go build -o catiSip
systemctl start catisip

# 前端 / Frontend
cd frontend && npm run build
sudo cp -r build/* /var/www/catisip/
```

---

## 扩展性 / Scalability

### 水平扩展 / Horizontal Scaling

1. **负载均衡器**
   - Nginx
   - HAProxy
   - AWS ELB/ALB

2. **多实例部署**
   ```bash
   # 启动多个后端实例
   PORT=8080 ./catiSip &
   PORT=8081 ./catiSip &
   PORT=8082 ./catiSip &
   ```

3. **会话管理**
   - 使用 Redis 存储会话
   - 启用粘性会话（Sticky Sessions）

---

## 支持 / Support

如遇问题，请：

If you encounter issues:

1. 查看文档 / Check documentation
2. 提交 Issue / Submit an issue
3. 联系开发团队 / Contact development team
