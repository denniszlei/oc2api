# oc2api

纯 Go 实现的 OpenCode API 代理服务，支持 SSE 流式响应。

## 配置

编辑 `config.yaml`：

| 配置项          | 默认值            | 说明               |
|--------------|----------------|------------------|
| `port`       | `8080`         | 监听端口             |
| `api-key`    | 空              | API 密钥（不设置则匿名访问） |
| `debug`      | `false`        | 调试日志             |
| `timeout-ms` | `300000` (5分钟) | 上游请求超时时间         |

免费模型筛选由以下环境变量控制（判定规则见[父项目 README 的免费模型限制](https://github.com/zhuweiyou/oc2api#免费模型限制)）：

| 环境变量                  | 默认值              | 说明                                          |
|-----------------------|------------------|---------------------------------------------|
| `EXTRA_MODELS`        | 空                | 逗号分隔，强制纳入模型列表（限免/匿名模型兜底，如 `union-alpha`）      |
| `BLOCK_MODELS`        | 空                | 逗号分隔，强制剔除，优先级最高                             |
| `MODEL_CACHE_TTL_MS`  | `600000`（10 分钟）  | 模型列表与 models.dev 目录的缓存有效期，决定限免模型多久后自动出现      |
| `OC_VERSION`          | 内置 `1.18.31`     | UA 里上报的 OpenCode 版本，上游抬高门槛时覆盖                |

## 本地部署

```bash
go build -ldflags="-s -w" -trimpath -o main main.go
./main
```

## Docker 部署

```bash
docker compose up -d
```

## Serverless 部署

本地/Docker 部署出口 IP 固定，建议部署到阿里云函数计算、腾讯云函数 等 Serverless 环境实现多出口 IP 轮询，规避 IP 限制。

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o main main.go
# 将 main + config.yaml 上传至云函数代码空间
# chmod +x ./main
# 启动命令: ./main
# 监听端口: 8080
```

## API

接口与父项目完全一致，详见 [父项目 README](https://github.com/zhuweiyou/oc2api/#api)。
