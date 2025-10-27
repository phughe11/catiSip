# API 文档 / API Documentation

## 基础信息 / Base Information

- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`
- **Authentication**: None (可根据需要添加 / Can be added as needed)

## 端点列表 / Endpoints

### 1. 健康检查 / Health Check

检查服务是否正常运行。

Check if the service is running properly.

**请求 / Request:**
```http
GET /api/health
```

**响应 / Response:**
```json
{
  "status": "healthy",
  "service": "catiSip"
}
```

**状态码 / Status Codes:**
- `200 OK`: 服务正常 / Service is healthy

---

### 2. 发起呼叫 / Make Call

发起一个新的 SIP 呼叫。

Initiate a new SIP call.

**请求 / Request:**
```http
POST /api/call/make
Content-Type: application/json

{
  "from": "1000",
  "to": "1001"
}
```

**请求参数 / Request Parameters:**

| 字段 / Field | 类型 / Type | 必填 / Required | 描述 / Description |
|------------|------------|----------------|-------------------|
| `from` | string | Yes | 主叫号码/分机 / Caller number/extension |
| `to` | string | Yes | 被叫号码/分机 / Callee number/extension |

**响应 / Response:**
```json
{
  "id": "call-1698765432",
  "from": "1000",
  "to": "1001",
  "status": "dialing",
  "start_time": "2024-10-27T10:30:32Z"
}
```

**状态码 / Status Codes:**
- `200 OK`: 呼叫发起成功 / Call initiated successfully
- `400 Bad Request`: 请求参数错误 / Invalid request parameters
- `500 Internal Server Error`: 服务器错误 / Server error

---

### 3. 挂断呼叫 / Hangup Call

终止一个正在进行的呼叫。

Terminate an active call.

**请求 / Request:**
```http
POST /api/call/hangup
Content-Type: application/json

{
  "call_id": "call-1698765432"
}
```

**请求参数 / Request Parameters:**

| 字段 / Field | 类型 / Type | 必填 / Required | 描述 / Description |
|------------|------------|----------------|-------------------|
| `call_id` | string | Yes | 呼叫 ID / Call ID |

**响应 / Response:**
```json
{
  "status": "call ended"
}
```

**状态码 / Status Codes:**
- `200 OK`: 呼叫已挂断 / Call hung up successfully
- `400 Bad Request`: 请求参数错误 / Invalid request parameters
- `404 Not Found`: 呼叫不存在 / Call not found

---

### 4. 查询呼叫状态 / Get Call Status

获取指定呼叫的当前状态。

Get the current status of a specific call.

**请求 / Request:**
```http
GET /api/call/status?call_id=call-1698765432
```

**查询参数 / Query Parameters:**

| 参数 / Parameter | 类型 / Type | 必填 / Required | 描述 / Description |
|----------------|------------|----------------|-------------------|
| `call_id` | string | Yes | 呼叫 ID / Call ID |

**响应 / Response:**
```json
{
  "id": "call-1698765432",
  "from": "1000",
  "to": "1001",
  "status": "answered",
  "start_time": "2024-10-27T10:30:32Z",
  "answer_time": "2024-10-27T10:30:37Z"
}
```

**呼叫状态 / Call Status Values:**
- `dialing`: 拨号中 / Dialing
- `ringing`: 响铃中 / Ringing
- `answered`: 已接听 / Answered
- `ended`: 已结束 / Ended

**状态码 / Status Codes:**
- `200 OK`: 成功获取呼叫状态 / Call status retrieved successfully
- `400 Bad Request`: 缺少 call_id 参数 / Missing call_id parameter
- `404 Not Found`: 呼叫不存在 / Call not found

---

### 5. 获取分机列表 / List Extensions

获取所有可用的 SIP 分机列表。

Get a list of all available SIP extensions.

**请求 / Request:**
```http
GET /api/extensions
```

**响应 / Response:**
```json
[
  {
    "extension": "1000",
    "status": "registered"
  },
  {
    "extension": "1001",
    "status": "registered"
  },
  {
    "extension": "1002",
    "status": "available"
  }
]
```

**分机状态 / Extension Status Values:**
- `registered`: 已注册 / Registered
- `available`: 可用但未注册 / Available but not registered
- `busy`: 忙线 / Busy

**状态码 / Status Codes:**
- `200 OK`: 成功获取分机列表 / Extensions retrieved successfully

---

## 数据模型 / Data Models

### Call Object

```typescript
{
  id: string;          // 呼叫唯一标识 / Unique call identifier
  from: string;        // 主叫号码 / Caller number
  to: string;          // 被叫号码 / Callee number
  status: string;      // 呼叫状态 / Call status
  start_time: string;  // 开始时间 (ISO 8601) / Start time
  answer_time?: string; // 接听时间 (可选) / Answer time (optional)
  end_time?: string;    // 结束时间 (可选) / End time (optional)
}
```

### Extension Object

```typescript
{
  extension: string;   // 分机号码 / Extension number
  status: string;      // 分机状态 / Extension status
}
```

---

## 错误处理 / Error Handling

所有错误响应遵循以下格式：

All error responses follow this format:

```json
{
  "error": "错误描述 / Error description"
}
```

### 常见错误码 / Common Error Codes

| 状态码 / Status Code | 描述 / Description |
|---------------------|-------------------|
| 400 | 错误的请求参数 / Bad request parameters |
| 404 | 资源未找到 / Resource not found |
| 405 | 不支持的请求方法 / Method not allowed |
| 500 | 服务器内部错误 / Internal server error |

---

## 使用示例 / Usage Examples

### cURL 示例 / cURL Examples

```bash
# 健康检查 / Health check
curl http://localhost:8080/api/health

# 发起呼叫 / Make a call
curl -X POST http://localhost:8080/api/call/make \
  -H "Content-Type: application/json" \
  -d '{"from": "1000", "to": "1001"}'

# 查询呼叫状态 / Get call status
curl http://localhost:8080/api/call/status?call_id=call-1698765432

# 挂断呼叫 / Hangup call
curl -X POST http://localhost:8080/api/call/hangup \
  -H "Content-Type: application/json" \
  -d '{"call_id": "call-1698765432"}'

# 获取分机列表 / List extensions
curl http://localhost:8080/api/extensions
```

### JavaScript/Fetch 示例 / JavaScript/Fetch Examples

```javascript
// 发起呼叫 / Make a call
async function makeCall(from, to) {
  const response = await fetch('http://localhost:8080/api/call/make', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ from, to }),
  });
  return await response.json();
}

// 使用 / Usage
makeCall('1000', '1001')
  .then(call => console.log('Call initiated:', call))
  .catch(error => console.error('Error:', error));
```

---

## CORS 支持 / CORS Support

服务器默认启用 CORS，允许所有来源的请求。

CORS is enabled by default, allowing requests from all origins.

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

---

## 限流和速率限制 / Rate Limiting

当前版本暂无限流限制。生产环境建议添加适当的限流机制。

No rate limiting is currently implemented. Consider adding rate limiting for production.

---

## 版本控制 / Versioning

当前 API 版本：v1

API 路径可能在未来版本中包含版本号，如 `/api/v1/...`

Current API version: v1

Future versions may include versioning in the path, e.g., `/api/v1/...`
