## go-micro测试grpc通信

### 安装依赖工具包

```
sudo apt-get -y install autoconf automake libtool
```

### 安装protobuf

```
mkdir ./tmp && cd ./tmp
git clone https://github.com/google/protobuf
cd protobuf
./autogen.sh
./configure
make
sudo make install
```

### 安装go的grpc相关包

```
github.com/micro/micro测试grpc通信
go get github.com/micro/go-web
go get -v github.com/micro/protobuf/{proto,protoc-gen-go}
go get -v github.com/micro/protoc-gen-micro

export LD_LIBRARY_PATH=/usr/local/lib
export PATH=$GOPATH/bin:$PATH
```

### 生成go版的proto

```
protoc --proto_path=proto --proto_path=/data/go/work/src --micro_out=proto --go_out=proto proto/hello/hello.proto
```