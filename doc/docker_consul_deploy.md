## docker部署服务注册中心consul

### 当前版本

docker: docker-ce18.06
consul: 1.4.2

### consul配置参数说明

```
–net=host docker参数, 使得docker容器越过了net namespace的隔离，免去手动指定端口映射的步骤
-server consul支持以server或client的模式运行, server是服务发现模块的核心, client主要用于转发请求
-advertise 将本机私有IP传递到consul
-retry-join 指定要加入的consul节点地址，失败会重试, 可多次指定不同的地址
-client consul绑定在哪个client地址上，这个地址提供HTTP、DNS、RPC等服务，默认是127.0.0.1
-bind 绑定服务器的ip地址；该地址用来在集群内部的通讯，集群内的所有节点到地址都必须是可达的，默认是0.0.0.0
allow_stale 设置为true, 表明可以从consul集群的任一server节点获取dns信息, false则表明每次请求都会经过consul server leade
-bootstrap-expect 数据中心中预期的服务器数。提供后，Consul将等待指定数量的服务器可用，然后启动群集。允许自动选举leader，但不能与传统-bootstrap标志一起使用, 需要在服务端模式下运行。
-data-dir 数据存放位置，用于持久化保存集群状态
-node 群集中此节点的名称，这在群集中必须是唯一的，默认情况下是节点的主机名。
-config-dir 指定配置文件，当这个目录下有 .json 结尾的文件就会被加载，详细可参考https://www.consul.io/docs/agent/options.html#configuration_files
-enable-script-checks 检查服务是否处于活动状态，类似开启心跳
-datacenter 数据中心名称
-ui 开启ui界面
-join 加入到已有的集群中
```

### consul端口用途说明

- 8500 : http 端口，用于 http 接口和 web ui访问
- 8300 : server rpc 端口，同一数据中心 consul server 之间通过该端口通信
- 8301 : serf lan 端口，同一数据中心 consul client 通过该端口通信; 用于处理当前datacenter中LAN的gossip
- 8302 : serf wan 端口，不同数据中心 consul server 通过该端口通信; agent Server使用，处理与其他datacenter的gossip
- 8600 : dns 端口，用于已注册的服务发现


### 启动一个server节点

```shell
docker run --name consul1 -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600 consul agent -server -bootstrap-expect 2 -ui -bind=0.0.0.0 -client=0.0.0.0
```

### 启动第二个server节点，并加入consul1

- 查看第一个server节点的ip地址

```shell
$docker inspect --format '{{ .NetworkSettings.IPAddress }}' consul1
172.17.0.2
```

- 启动第二个server节点

```shell
docker run --name consul2 -d -p 8501:8500 consul agent -server -ui -bind=0.0.0.0 -client=0.0.0.0 -join 172.17.0.2
```

### 启动第三个server节点, 并加入consul

```shell
docker run --name consul3 -d -p 8502:8500 consul agent -server -ui -bind=0.0.0.0 -client=0.0.0.0 -join 172.17.0.2
```

### 查看consul集群成员信息

```shell
$docker exec -it consul1 consul members
Node          Address          Status  Type    Build  Protocol  DC   Segment
0392bb73a784  172.17.0.3:8301  alive   server  1.4.2  2         dc1  <all>
39ce8be84203  172.17.0.4:8301  alive   server  1.4.2  2         dc1  <all>
c8e93300cf75  172.17.0.2:8301  alive   server  1.4.2  2         dc1  <all>
```

### 进入ui界面

通过`http://localhost:8500`可以打开ui界面；(8501或8502也可以)

![consul ui](https://git.imooc.com/coding-323/filestore-server/src/master/doc/consul-ui.png)
