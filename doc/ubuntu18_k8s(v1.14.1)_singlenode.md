## 1.准备工作

- 关闭防火墙

```bash
sudo ufw disable
```

- 关闭系统swap

```bash
sudo swapoff -a
```

- 安装docker

```
建议版本: docker-ce-18.0x
```

## 2.准备k8s安装环境

- 下载并添加Kubernetes安装的密钥

```bash
sudo apt update && sudo apt install -y apt-transport-https curl
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -
```

- 配置k8s的源

```bash
sudo touch /etc/apt/sources.list.d/kubernetes.list 
sudo echo "deb https://mirrors.aliyun.com/kubernetes/apt kubernetes-xenial main" >> /etc/apt/sources.list.d/kubernetes.list
```

- 安装kubeadm及kubelet等工具

```bash
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
# 保持版本，取消自动更新
sudo apt-mark hold kubelet kubeadm kubectl
```

## 3.kubeadm初始化集群

- 执行初始化操作

```bash
sudo kubeadm init --image-repository registry.aliyuncs.com/google_containers --kubernetes-version v1.14.1 --pod-network-cidr=10.240.0.0/16
```

成功的话，会输出类似以下的这些内容:
```
// ...
Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 192.168.200.212:6443 --token q1guce.z76o2a2bb65vhd0u \
    --discovery-token-ca-cert-hash sha256:2a57a27853c66d608bc544742b57602a21d47c3d09fe58eef15258946d4341c0 
```

- 如果我们想在非root用户下操作kubectl命令, 可以这样配置:

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

- 查看k8s集群的pod状态

```bash
xiaomo@xiaomo:~$ kubectl get pods --all-namespaces
NAMESPACE     NAME                             READY   STATUS    RESTARTS   AGE
kube-system   coredns-8686dcc4fd-4kjld         0/1     Pending   0          14m
kube-system   coredns-8686dcc4fd-4z6st         0/1     Pending   0          14m
kube-system   etcd-xiaomo                      1/1     Running   0          13m
kube-system   kube-apiserver-xiaomo            1/1     Running   0          13m
kube-system   kube-controller-manager-xiaomo   1/1     Running   0          13m
kube-system   kube-proxy-n7xq9                 1/1     Running   0          14m
kube-system   kube-scheduler-xiaomo            1/1     Running   0          14m
```

可以看到目前`coredns`处于`pending`状态，一般是因为还没安装网络插件(这里选calico)．

-  安装Canal插件

```bash
# 从国内的一个源里下载
kubectl apply -f http://mirror.faasx.com/k8s/calico/v3.3.2/rbac-kdd.yaml
kubectl apply -f http://mirror.faasx.com/k8s/calico/v3.3.2/calico.yaml
```

-  Master隔离解除(主节点也能部署工作任务)

`kubectl taint nodes --all node-role.kubernetes.io/master-`
```bash
# 成功时会输出类似提示:
node/xiaomo untainted
```

## 4.加入工作节点(假设有其他的节点需要加入到集群)

```bash
1)SSH到目标机器
2)切换至root用户，如: sudo su
2)根据上面`kubeadm init`命令得到的提示结果，运行`kubeadm join`：

kubeadm join --token <token> <master-ip>:<master-port> --discovery-token-ca-cert-hash sha256:<hash>

# kubeadm token list (可以获取<token>的值)
# kubeadm token create (可以创建<token>新的值)

# --discovery-token-ca-cert-hash (以下命令可以获取 --discovery-token-ca-cert-hash的值)
# openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```

## 5.查看集群状态

```bash
xiaomo@xiaomo:~$ kubectl get pods -n kube-system
NAME                                       READY   STATUS    RESTARTS   AGE
calico-node-tmcmn                          2/2     Running   0          55m
coredns-8686dcc4fd-4kjld                   1/1     Running   5          79m
coredns-8686dcc4fd-4z6st                   1/1     Running   5          79m
etcd-xiaomo                                1/1     Running   0          78m
kube-apiserver-xiaomo                      1/1     Running   0          79m
kube-controller-manager-xiaomo             1/1     Running   0          78m
kube-proxy-n7xq9                           1/1     Running   0          79m
kube-scheduler-xiaomo                      1/1     Running   0          79m
kubernetes-dashboard-5f7b999d65-f9m7d      1/1     Running   0          6m24s
```