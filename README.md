# ebpf-study
eBPF - extended Berkeley Packet Filter


## Environment

为了快速的启动一个 **Ubuntu21.10**，这里使用 [**Multipass**](https://feifeifeimoon.github.io/posts/5ba8a50c.html) 来搭建的 **ebpf** 开发环境
```bash
# 查看可用的镜像
$ multipass find                    
Image                       Aliases           Version          Description
18.04                       bionic            20220302         Ubuntu 18.04 LTS
20.04                       focal,lts         20220207         Ubuntu 20.04 LTS
21.10                       impish            20220201         Ubuntu 21.10
anbox-cloud-appliance                         latest           Anbox Cloud Appliance
charm-dev                                     latest           A development and testing environment for charmers
docker                                        latest           A Docker environment with Portainer and related tools
minikube                                      latest           minikube is local Kubernetes

# 启动一个2核4G的
$ multipass launch 21.10  -n ubuntu -c 2 -m 4G -d 20G 

# 查看运行中的虚拟机
$ multipass list
Name                    State             IPv4             Image
ubuntu                  Running           192.168.64.6     Ubuntu 21.10
```

这样就得到了一个 **Ubuntu21.10** 环境。一般开发的话都习惯使用 **SSH** 进行 **Remote** 开发，**Multipass** 开启**SSH** 需要额外进行一些配置。

```bash
# 进入虚拟机
$ multipass shell ubuntu

# 设置root密码
$ sudo passwd

# 修改sshd配置文件 允许root登陆以及打开password认证
$ vim /etc/ssh/sshd_config

PermitRootLogin yes
PasswordAuthentication yes

# 重启SSH服务
$ sudo service ssh restart 
```


### 安装 **ebpf** 需要用到的工具

+ **Clang/LLVM** ：用来将 **C** 代码编译生成 **eBPF** 字节码。(**GCC** 目前也在支持，但没有 **Clang/LLVM** 完善) [BPF in GCC](https://lwn.net/Articles/831402/)
+ **libbpf** ：[**libbpf**](https://github.com/libbpf/libbpf) 是可以在用户空间和 **eBPF** 程序中导入的库。它为开发人员提供了一个用于加载 **eBPF** 程序并与之交互的 **API**。
+ **bpftool** ：内核代码提供的 **eBPF** 程序管理工具

```bash
# 安装编译工具链
$ sudo apt-get install clang llvm libelf-dev libbpf-dev libbfd-dev libreadline-dev bison flex

# 查看当前的内核版本
$ uname -nr
ubuntu 5.13.0-28-generic

# 安装当前版本的内核源码 libbpf和bpftool都在内核源码中有一份
$ apt-cache search linux-source
linux-source - Linux kernel source with Ubuntu patches
linux-source-5.13.0 - Linux kernel source for version 5.13.0 with Ubuntu patches
$ apt install linux-source-5.13.0

# 会安装到 /usr/src下
$ cd /usr/src
$ tar -jxvf linux-source-5.11.0.tar.bz2

# libbpf
$ cd /usr/src/linux-source-5.13.0/tools/bpf
$ make && make install prefix=/usr/local

# bpftool
$ cd /usr/src/linux-source-5.13.0/tools/bpf/bpftool
$ make && make install 
# 确认bpftool安装成功
$ bpftool version -p 
{
    "version": "5.13.19",
    "features": {
        "libbfd": true,
        "skeletons": true
    }
}
```

