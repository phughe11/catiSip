# 测试指南 / Testing Guide

## 概述 / Overview

本文档介绍如何测试 CatiSip 系统的各个组件。

This document describes how to test the various components of the CatiSip system.

---

## 后端测试 / Backend Testing

### 运行所有测试 / Run All Tests

```bash
cd backend
go test ./... -v
```

### 运行特定包的测试 / Run Tests for Specific Package

```bash
# 测试 SIP 客户端 / Test SIP client
go test ./sip -v

# 测试 HTTP 处理器 / Test HTTP handlers
go test ./handlers -v
```

### 测试覆盖率 / Test Coverage

```bash
# 生成覆盖率报告 / Generate coverage report
go test ./... -coverprofile=coverage.out

# 查看覆盖率 / View coverage
go tool cover -html=coverage.out
```

### 性能测试 / Benchmark Tests

```bash
go test ./... -bench=. -benchmem
```

---

## 前端测试 / Frontend Testing

### 运行测试 / Run Tests

```bash
cd frontend
npm test
```

### 运行测试（监视模式）/ Run Tests (Watch Mode)

```bash
npm test -- --watch
```

### 测试覆盖率 / Test Coverage

```bash
npm test -- --coverage
```

---

## 集成测试 / Integration Testing

### 手动集成测试 / Manual Integration Testing

#### 1. 启动服务 / Start Services

```bash
# 启动后端 / Start backend
cd backend
./catiSip &

# 等待几秒 / Wait a few seconds
sleep 2

# 启动前端 / Start frontend (in another terminal)
cd frontend
npm start
```

#### 2. 测试 API 端点 / Test API Endpoints

```bash
# 健康检查 / Health check
curl http://localhost:8080/api/health

# 发起呼叫 / Make a call
curl -X POST http://localhost:8080/api/call/make \
  -H "Content-Type: application/json" \
  -d '{"from": "1000", "to": "1001"}'

# 查询呼叫状态 / Get call status
curl "http://localhost:8080/api/call/status?call_id=call-XXXXX"

# 获取分机列表 / List extensions
curl http://localhost:8080/api/extensions

# 挂断呼叫 / Hangup call
curl -X POST http://localhost:8080/api/call/hangup \
  -H "Content-Type: application/json" \
  -d '{"call_id": "call-XXXXX"}'
```

### 自动化集成测试脚本 / Automated Integration Test Script

创建测试脚本 `test-integration.sh`:

Create test script `test-integration.sh`:

```bash
#!/bin/bash

set -e

echo "Starting integration tests..."

# Start backend
cd backend
./catiSip > /tmp/backend-test.log 2>&1 &
BACKEND_PID=$!
cd ..

# Wait for backend to be ready
sleep 3

# Test health endpoint
echo "Testing health endpoint..."
HEALTH=$(curl -s http://localhost:8080/api/health)
if [[ $HEALTH == *"healthy"* ]]; then
    echo "✓ Health check passed"
else
    echo "✗ Health check failed"
    kill $BACKEND_PID
    exit 1
fi

# Test make call endpoint
echo "Testing make call endpoint..."
CALL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/call/make \
    -H "Content-Type: application/json" \
    -d '{"from":"1000","to":"1001"}')

CALL_ID=$(echo $CALL_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [[ ! -z "$CALL_ID" ]]; then
    echo "✓ Make call passed (Call ID: $CALL_ID)"
else
    echo "✗ Make call failed"
    kill $BACKEND_PID
    exit 1
fi

# Test call status endpoint
echo "Testing call status endpoint..."
STATUS_RESPONSE=$(curl -s "http://localhost:8080/api/call/status?call_id=$CALL_ID")
if [[ $STATUS_RESPONSE == *"$CALL_ID"* ]]; then
    echo "✓ Call status passed"
else
    echo "✗ Call status failed"
    kill $BACKEND_PID
    exit 1
fi

# Test extensions endpoint
echo "Testing extensions endpoint..."
EXT_RESPONSE=$(curl -s http://localhost:8080/api/extensions)
if [[ $EXT_RESPONSE == *"1000"* ]]; then
    echo "✓ Extensions list passed"
else
    echo "✗ Extensions list failed"
    kill $BACKEND_PID
    exit 1
fi

# Test hangup endpoint
echo "Testing hangup endpoint..."
HANGUP_RESPONSE=$(curl -s -X POST http://localhost:8080/api/call/hangup \
    -H "Content-Type: application/json" \
    -d "{\"call_id\":\"$CALL_ID\"}")
if [[ $HANGUP_RESPONSE == *"ended"* ]]; then
    echo "✓ Hangup call passed"
else
    echo "✗ Hangup call failed"
    kill $BACKEND_PID
    exit 1
fi

# Cleanup
kill $BACKEND_PID

echo ""
echo "All integration tests passed! ✓"
```

运行集成测试：

Run integration tests:

```bash
chmod +x test-integration.sh
./test-integration.sh
```

---

## 端到端测试 / End-to-End Testing

### 使用 Playwright 或 Selenium

对于前端 E2E 测试，可以使用 Playwright 或 Selenium。

For frontend E2E testing, you can use Playwright or Selenium.

#### 安装 Playwright / Install Playwright

```bash
cd frontend
npm install --save-dev @playwright/test
npx playwright install
```

#### 示例 E2E 测试 / Example E2E Test

创建 `frontend/e2e/call-flow.spec.js`:

Create `frontend/e2e/call-flow.spec.js`:

```javascript
const { test, expect } = require('@playwright/test');

test('make a call flow', async ({ page }) => {
  // Navigate to the app
  await page.goto('http://localhost:3000');

  // Wait for the app to load
  await expect(page.locator('h1')).toContainText('CatiSip');

  // Select extension
  await page.selectOption('select', '1000');

  // Enter destination number
  await page.fill('input[placeholder*="phone number"]', '1001');

  // Click make call button
  await page.click('button:has-text("Make Call")');

  // Verify call was initiated
  await expect(page.locator('.status-message')).toContainText('Call initiated');

  // Verify call appears in active calls
  await expect(page.locator('.call-item')).toBeVisible();
});
```

运行 E2E 测试：

Run E2E tests:

```bash
npx playwright test
```

---

## 性能测试 / Performance Testing

### 使用 Apache Bench

```bash
# 测试健康检查端点 / Test health endpoint
ab -n 1000 -c 10 http://localhost:8080/api/health

# 测试分机列表端点 / Test extensions endpoint
ab -n 1000 -c 10 http://localhost:8080/api/extensions
```

### 使用 wrk

```bash
# 安装 wrk / Install wrk
sudo apt-get install wrk

# 性能测试 / Performance test
wrk -t4 -c100 -d30s http://localhost:8080/api/health
```

---

## 负载测试 / Load Testing

### 使用 k6

```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
};

export default function () {
  // Test health endpoint
  let healthRes = http.get('http://localhost:8080/api/health');
  check(healthRes, {
    'health status is 200': (r) => r.status === 200,
  });

  // Test make call endpoint
  let callRes = http.post(
    'http://localhost:8080/api/call/make',
    JSON.stringify({ from: '1000', to: '1001' }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(callRes, {
    'make call status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
```

运行负载测试：

Run load test:

```bash
k6 run load-test.js
```

---

## 安全测试 / Security Testing

### 使用 OWASP ZAP

```bash
# 启动 ZAP 代理 / Start ZAP proxy
docker run -u zap -p 8090:8090 -i owasp/zap2docker-stable \
  zap-baseline.py -t http://localhost:8080
```

### SQL 注入测试 / SQL Injection Testing

虽然当前实现不使用数据库，但如果添加数据库，请确保：

While the current implementation doesn't use a database, if you add one, ensure:

- 使用参数化查询 / Use parameterized queries
- 验证和清理所有输入 / Validate and sanitize all inputs
- 使用 ORM 或查询构建器 / Use an ORM or query builder

---

## 持续集成 / Continuous Integration

### GitHub Actions 示例 / GitHub Actions Example

创建 `.github/workflows/test.yml`:

Create `.github/workflows/test.yml`:

```yaml
name: Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: |
        cd backend
        go test ./... -v -coverprofile=coverage.out
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./backend/coverage.out

  frontend-test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
    
    - name: Install dependencies
      run: |
        cd frontend
        npm ci
    
    - name: Run tests
      run: |
        cd frontend
        npm test -- --coverage
```

---

## 测试最佳实践 / Testing Best Practices

1. **编写测试先行** / Write Tests First
   - 实践 TDD (测试驱动开发) / Practice TDD

2. **保持测试独立** / Keep Tests Independent
   - 每个测试应该能够独立运行 / Each test should run independently

3. **使用有意义的测试名称** / Use Meaningful Test Names
   - 测试名称应该描述测试的内容 / Test names should describe what is being tested

4. **测试边界情况** / Test Edge Cases
   - 不仅测试正常情况 / Don't just test the happy path

5. **模拟外部依赖** / Mock External Dependencies
   - 使用 mock 来隔离测试 / Use mocks to isolate tests

6. **定期运行测试** / Run Tests Regularly
   - 在 CI/CD 管道中自动运行 / Run automatically in CI/CD pipeline

---

## 故障排除 / Troubleshooting

### 测试失败 / Tests Fail

1. 检查依赖是否安装 / Check if dependencies are installed
2. 验证环境变量 / Verify environment variables
3. 查看详细错误信息 / Review detailed error messages

### 性能测试超时 / Performance Tests Timeout

1. 增加超时时间 / Increase timeout values
2. 减少并发数 / Reduce concurrency
3. 优化代码 / Optimize code

---

## 资源 / Resources

- [Go Testing](https://golang.org/pkg/testing/)
- [Jest Documentation](https://jestjs.io/)
- [Playwright Documentation](https://playwright.dev/)
- [k6 Documentation](https://k6.io/docs/)
