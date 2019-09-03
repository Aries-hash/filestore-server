## centos7安装docker

- 1. yum源使用阿里云的源

```bash
#方法1
cd /etc/yum.repos.d/
#下载阿里云yum源
wget http://mirrors.aliyun.com/repo/Centos-7.repo
mv CentOS-Base.repo CentOS-Base.repo.bak
mv Centos-7.repo CentOS-Base.repo

# 方法2
sudo yum install -y yum-utils device-mapper-persistent-data lvm2
sudo yum-config-manager \
    --add-repo \
    https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
```

- 2. 重置yum源

```bash
yum clean all
yum makecache
```

- 3. 开始安装docker

```bash
#查看阿里云上docker 源信息
yum list docker-ce
#安装docker最新社区版(截止目前最新版本是18.09)
yum -y install docker-ce
#查看docker版本
docker -v
#启动docker
systemctl start docker
#查看docker详细状态信息
docker info
```