# Consul Service
>监听Consul service register 和 Consul KV. 当Endpoint发生变更时，实时同步到Envoy

[![Build Status](https://travis-ci.com/tio-serverless/consul.svg?branch=master)](https://travis-ci.com/tio-serverless/consul)

### SideCar

### Watcher
> 监控Consul KV变化，生成Envoy数据并推送给Envoy

+ Environment

   - CONSUL_ADDRESS
   - DEBUG
   - TIO_CONSUL_CLUSTER_HTTP
   - TIO_CONSUL_CLUSTER_GRPC
   - TIO_CONSUL_CLUSTER_TCP
   - MY_GRPC_PORT