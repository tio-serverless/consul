# Consul Service
>监听Consul service register 和 Consul KV. 当Endpoint发生变更时，实时同步到Envoy

[![Build Status](https://travis-ci.com/tio-serverless/consul.svg?branch=master)](https://travis-ci.com/tio-serverless/consul) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=tio-serverless_consul&metric=alert_status)](https://sonarcloud.io/dashboard?id=tio-serverless_consul)

## 流程图

![](./consul.svg)

### 组件描述

- SideCar
> 作为服务的Sidecar，当服务启动后，将服务数据注册到Consul中

- Watcher
> 监控Consul KV变化，生成Envoy数据并推送给Envoy

   + 启动参数
      - DEBUG
      - CONSUL_ADDRESS
      - MY_GRPC_PORT

#### 细节

+ Watcher

![](./consul-watcher.svg)

`Sidecar`在启动时会将自身服务数据注册到`Consul`中(先注册实例数据Health Check,然后注册KV数据(服务属性))。 `Watcher`会通过监听`Health Check`和`KV`接受到这两个事件。

而后通过读取`KV`中的服务属性数据获取到服务元数据，然后判断其`Health Check`是否成功，若实例健康检查成功，则生成最新的服务路由数据(包括cluster和router)，并通过chan通知`envoy watcher`。

`envoy watcher`分别通过`StreamCluster`和`StreamRouter`分别将生成的路由数据下发到`Envoy`实例中，从而完成服务注册，服务发现和服务路由.


+ Envoy default bootstrap

```yaml
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }

dynamic_resources:
  cds_config:
    api_config_source:
      api_type: GRPC
      grpc_services:
        envoy_grpc:
          cluster_name: xds_cluster

static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 80 }
      filter_chains:
        - filters:
            - name: envoy.http_connection_manager
              config:
                stat_prefix: ingress_http
                codec_type: AUTO
                access_log:
                  name: envoy.file_access_log
                  config:
                    path: /dev/stdout
                rds:
                  route_config_name: tio
                  config_source:
                    api_config_source:
                      api_type: GRPC
                      grpc_services:
                        envoy_grpc:
                          cluster_name: xds_cluster
                http_filters:
                  - name: envoy.grpc_web
                  - name: envoy.cors
                  - name: envoy.router
  clusters:
    - name: xds_cluster
      connect_timeout: 60s
      type: STATIC
      lb_policy: ROUND_ROBIN
      http2_protocol_options: {}
      hosts: [{ socket_address: { address: 0.0.0.0, port_value: 8000 }}]
    - name: http_proxy_cluster
      connect_timeout: 60s
      type: STRICT_DNS
      lb_policy: ROUND_ROBIN
      http2_protocol_options: {}
      hosts: [{ socket_address: { address: proxy-http.tio.svc.cluster.local, port_value: 80 }}]
```