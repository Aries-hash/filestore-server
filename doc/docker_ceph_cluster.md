通过docker可以快速部署小规模Ceph集群的流程，可用于开发测试。
以下的安装流程是通过linux shell来执行的；假设你只有一台机器，装了linux(如Ubuntu)系统和docker环境，那么可以参考以下步骤安装Ceph:

```bash
# 要用root用户创建, 或有sudo权限
# 注: 建议使用这个docker镜像源:https://registry.docker-cn.com
# 1. 修改docker镜像源
cat > /etc/docker/daemon.json << EOF
{
  "registry-mirrors": [
    "https://registry.docker-cn.com"
  ]
}
EOF
# 重启docker
systemctl restart docker
# 2. 创建Ceph专用网络
docker network create --driver bridge --subnet 172.20.0.0/16 ceph-network
docker network inspect ceph-network
# 3. 删除旧的ceph相关容器
docker rm -f $(docker ps -a | grep ceph | awk '{print $1}')
# 4. 清理旧的ceph相关目录文件，加入有的话
rm -rf /www/ceph /var/lib/ceph/  /www/osd/
# 5. 创建相关目录及修改权限，用于挂载volume
mkdir -p /www/ceph /var/lib/ceph/osd /www/osd/
chown -R 64045:64045 /var/lib/ceph/osd/
chown -R 64045:64045 /www/osd/
# 6. 创建monitor节点
docker run -itd --name monnode --network ceph-network --ip 172.20.0.10 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph ceph/mon
# 7. 在monitor节点上标识3个osd节点
docker exec monnode ceph osd create
docker exec monnode ceph osd create
docker exec monnode ceph osd create
# 8. 创建OSD节点
docker run -itd --name osdnode0 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph -v /www/osd/0:/var/lib/ceph/osd/ceph-0 ceph/osd 
docker run -itd --name osdnode1 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph -v /www/osd/1:/var/lib/ceph/osd/ceph-1 ceph/osd
docker run -itd --name osdnode2 --network ceph-network -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /www/ceph:/etc/ceph -v /www/osd/2:/var/lib/ceph/osd/ceph-2 ceph/osd
# 9. 增加monitor节点，组件成集群
docker run -itd --name monnode_1 --network ceph-network --ip 172.20.0.11 -e MON_NAME=monnode_1 -e MON_IP=172.20.0.11 -v /www/ceph:/etc/ceph ceph/mon
docker run -itd --name monnode_2 --network ceph-network --ip 172.20.0.12 -e MON_NAME=monnode_2 -e MON_IP=172.20.0.12 -v /www/ceph:/etc/ceph ceph/mon
# 10. 创建gateway节点
docker run -itd --name gwnode --network ceph-network --ip 172.20.0.9 -p 9080:80 -e RGW_NAME=gwnode -v /www/ceph:/etc/ceph ceph/radosgw
# 11. 查看ceph集群状态
sleep 10 && docker exec monnode ceph -s
# 12. 创建用户
docker exec -it gwnode radosgw-admin user create --uid=user1 --display-name=user1
```