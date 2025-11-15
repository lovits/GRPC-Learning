// https://www.liwenzhou.com/posts/Go/consul/
// 
// http://127.0.0.1:8500/
//
// git clone https://github.com/hashicorp/learn-consul-docker.git
//
// cd datacenter-deploy-service-discovery
//
// docker-compose up -d
//
//使用浏览器打开 http:127.0.0.1:8500 就可以看到Consul的管理界面了

postman 
Get http://localhost:8500/v1/agent/services
Put http://localhost:8500/v1/agent/services
{
    "id":"goods-1",
    "name":"order",
    "tags":[
        "cheng",
        "goods"
    ],
    "address":"127.0.0.1",
    "port":60091
}

{
    "id":"order-1",
    "name":"order",
    "tags":[
        "cheng",
        "order"
    ],
    "address":"127.0.0.1",
    "port":60081
}

Get http://localhost:8500/v1/agent/services?filter="cheng" in Tags

GRPC实现

# 多实例 gRPC + Consul 示例（无健康检查）

## 功能说明

- 启动 **3 个 gRPC 服务实例**（端口 50051、50052、50053）
- 将服务注册到本地 **Consul**，不使用健康检查
- 客户端使用 **grpc-consul-resolver** 自动发现所有实例
- 采用 **round-robin 轮询**方式调用多个实例

## 目录结构

grpc-consul-demo/
├── proto/
│   └── greeter.proto       # gRPC Proto 文件
├── server/
│   └── main.go             # 多实例服务端
├── client/
│   └── main.go             # 客户端调用示例
├── docker-compose.yml      # Consul 容器

## 环境依赖

- Go >= 1.20
- Docker Desktop（用于 Consul）
- protoc + protoc-gen-go + protoc-gen-go-grpc
- grpc-consul-resolver

## 使用步骤

1. 启动 Consul：

```bash
docker compose up -d
打开http:127.0.0.1:8500

启动服务
go run main.go

客户端启动
cd client
go run main.go