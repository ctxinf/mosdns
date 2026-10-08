# mosdns

原项目：[IrineSistiana/mosdns](https://github.com/IrineSistiana/mosdns)

## 修改原因

原版 `forward` 的查询、DoH 请求和 TLS 握手超时是固定的。距离较远、网络延迟较高的 DoH/DNS 上游（包括通过代理访问的上游）可能在收到响应前就超时，导致查询频繁失败。此版本允许按实际线路延迟配置这些超时，并允许调整 DNS 入口的查询超时，让上游有足够时间返回结果。

## 超时配置

以下超时单位均为**秒**；不配置或配置为非正数时，沿用 v5.3.4 的默认值。`fallback.threshold` 的单位仍是毫秒。

```yaml
plugins:
  - tag: smartdns
    type: forward
    args:
      concurrent: 3
      query_timeout: 20             # 每个上游查询；默认 5 秒
      tls_handshake_timeout: 15     # DoH TLS 握手；默认 3 秒
      doh_request_timeout: 20       # 整次 DoH HTTP 请求；默认 6 秒
      upstreams:
        - addr: https://dns.example/dns-query
          socks5: 127.0.0.1:1080
          # 可在单个上游覆盖以上三个超时
  - tag: dns_udp
    type: udp_server
    args:
      entry: main_sequence
      listen: 127.0.0.1:5335
      query_timeout: 20             # 整次入口查询；默认 5 秒
  - tag: dns_tcp
    type: tcp_server
    args:
      entry: main_sequence
      listen: 127.0.0.1:5335
      query_timeout: 20             # 整次入口查询；默认 5 秒
```

如果使用 `fallback.threshold: 15000` 等待较慢的上游，入口 `query_timeout` 需要大于 15 秒，`forward.query_timeout` 和 DoH 的 `doh_request_timeout` 也应足够长。TLS 握手超时应根据线路延迟设置，且不应长于 DoH 请求超时。
