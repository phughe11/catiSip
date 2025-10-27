# FreeSWITCH 集成指南 / FreeSWITCH Integration Guide

## 概述 / Overview

CatiSip 通过 FreeSWITCH 的 Event Socket Library (ESL) 与 FreeSWITCH 进行集成。

CatiSip integrates with FreeSWITCH through the Event Socket Library (ESL).

## FreeSWITCH 安装 / FreeSWITCH Installation

### Debian/Ubuntu

```bash
wget -O - https://files.freeswitch.org/repo/deb/debian-release/fsstretch-archive-keyring.asc | apt-key add -
echo "deb http://files.freeswitch.org/repo/deb/debian-release/ stretch main" > /etc/apt/sources.list.d/freeswitch.list
apt-get update
apt-get install -y freeswitch-meta-all
```

### Docker

```bash
docker pull freeswitch/freeswitch
docker run -d --name freeswitch \
  -p 5060:5060/udp \
  -p 5060:5060/tcp \
  -p 8021:8021 \
  freeswitch/freeswitch
```

## FreeSWITCH 配置 / FreeSWITCH Configuration

### 1. 启用 Event Socket / Enable Event Socket

编辑 `/etc/freeswitch/autoload_configs/event_socket.conf.xml`:

Edit `/etc/freeswitch/autoload_configs/event_socket.conf.xml`:

```xml
<configuration name="event_socket.conf" description="Socket Client">
  <settings>
    <param name="nat-map" value="false"/>
    <param name="listen-ip" value="0.0.0.0"/>
    <param name="listen-port" value="8021"/>
    <param name="password" value="ClueCon"/>
  </settings>
</configuration>
```

### 2. 配置 SIP 分机 / Configure SIP Extensions

编辑 `/etc/freeswitch/directory/default/1000.xml`:

Edit `/etc/freeswitch/directory/default/1000.xml`:

```xml
<include>
  <user id="1000">
    <params>
      <param name="password" value="1234"/>
      <param name="vm-password" value="1000"/>
    </params>
    <variables>
      <variable name="toll_allow" value="domestic,international,local"/>
      <variable name="accountcode" value="1000"/>
      <variable name="user_context" value="default"/>
      <variable name="effective_caller_id_name" value="Extension 1000"/>
      <variable name="effective_caller_id_number" value="1000"/>
    </variables>
  </user>
</include>
```

创建多个分机 (1001.xml, 1002.xml 等)。

Create multiple extensions (1001.xml, 1002.xml, etc.).

### 3. 重启 FreeSWITCH / Restart FreeSWITCH

```bash
systemctl restart freeswitch
```

## CatiSip 配置 / CatiSip Configuration

更新环境变量以连接到 FreeSWITCH：

Update environment variables to connect to FreeSWITCH:

```bash
export SIP_HOST=localhost
export SIP_PORT=5060
export SIP_USERNAME=1000
export SIP_PASSWORD=1234
export SIP_DOMAIN=localhost
```

## 测试连接 / Testing Connection

### 使用 fs_cli 测试 / Test with fs_cli

```bash
fs_cli -H localhost -P ClueCon

# 检查注册的分机 / Check registered extensions
sofia status profile internal reg

# 发起测试呼叫 / Originate test call
originate user/1000 &echo
```

### 使用 CatiSip API / Test with CatiSip API

```bash
# 健康检查 / Health check
curl http://localhost:8080/api/health

# 发起呼叫 / Make a call
curl -X POST http://localhost:8080/api/call/make \
  -H "Content-Type: application/json" \
  -d '{"from": "1000", "to": "1001"}'
```

## ESL 库集成 / ESL Library Integration

当前实现使用模拟客户端。在生产环境中，建议使用以下 ESL 库之一：

The current implementation uses a simulated client. For production, use one of these ESL libraries:

### Go ESL 库 / Go ESL Libraries

1. **freeswitch-esl-go**
   ```bash
   go get github.com/fiorix/go-eventsocket
   ```

2. **goesl**
   ```bash
   go get github.com/0x19/goesl
   ```

### 集成示例 / Integration Example

```go
package sip

import (
    "github.com/fiorix/go-eventsocket/eventsocket"
)

func NewClient(cfg config.SIPConfig) (*Client, error) {
    conn, err := eventsocket.Dial(
        fmt.Sprintf("%s:%d", cfg.Host, 8021),
        "ClueCon",
    )
    if err != nil {
        return nil, err
    }
    
    // Subscribe to events
    conn.Send("events plain ALL")
    
    return &Client{
        conn: conn,
        config: &cfg,
        calls: make(map[string]*Call),
    }, nil
}
```

## 故障排除 / Troubleshooting

### FreeSWITCH 无法启动 / FreeSWITCH Won't Start

```bash
# 检查日志 / Check logs
tail -f /var/log/freeswitch/freeswitch.log

# 检查端口占用 / Check port usage
netstat -tlnp | grep 5060
```

### 分机无法注册 / Extension Won't Register

1. 检查防火墙设置 / Check firewall settings
2. 验证用户名和密码 / Verify username and password
3. 检查 SIP 配置 / Check SIP configuration

```bash
fs_cli -x "sofia status"
fs_cli -x "sofia status profile internal"
```

### 呼叫无法建立 / Calls Won't Establish

```bash
# 启用详细日志 / Enable verbose logging
fs_cli -x "console loglevel debug"

# 监控呼叫 / Monitor calls
fs_cli -x "show calls"
```

## 安全建议 / Security Recommendations

1. 更改默认 ESL 密码 / Change default ESL password
2. 使用强密码 / Use strong passwords for SIP accounts
3. 限制网络访问 / Restrict network access
4. 启用 TLS / Enable TLS for SIP
5. 配置防火墙规则 / Configure firewall rules

## 参考资源 / References

- [FreeSWITCH 官方文档](https://freeswitch.org/confluence/)
- [Event Socket Library](https://freeswitch.org/confluence/display/FREESWITCH/Event+Socket+Library)
- [SIP RFC 3261](https://tools.ietf.org/html/rfc3261)
