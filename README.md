# mosdns

功能概述、配置方式、教程等，详见: [wiki](https://irine-sistiana.gitbook.io/mosdns-wiki/)

下载预编译文件、更新日志，详见: [release](https://github.com/IrineSistiana/mosdns/releases)

docker 镜像: [docker hub](https://hub.docker.com/r/irinesistiana/mosdns)

## 可配置的查询超时

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
          # 可在单个上游覆盖 query_timeout、tls_handshake_timeout 和 doh_request_timeout
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

要让 `threshold: 15000` 有机会等待满 15 秒，入口 `query_timeout` 必须大于 15 秒；forward 的 `query_timeout` 和 DoH 的 `doh_request_timeout` 也应足够长。TLS 握手超时应按代理线路延迟设置，且不应长于 DoH 请求超时。
