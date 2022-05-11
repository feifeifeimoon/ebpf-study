# ebpf-study

![logo](https://github.com/feifeifeimoon/ebpf-study/blob/main/img/logo.png?raw=true)
eBPF - extended Berkeley Packet Filter

记录 **eBPF** 学习过程中的资源和动手实现的一些例子🌰。

## Website

+ [ebpf.io](https://ebpf.io/) : **eBPF** 官网，由  **Linux Foundation** 主持。

### 博客
+ [ARTHURCHIAO'S BLOG](https://arthurchiao.art/articles-zh/) : 携程大佬的博客，里面有很多关于 **eBPF** 文章的翻译，包括的[指导手册](https://arthurchiao.art/blog/cilium-bpf-xdp-reference-guide-zh/) 。
+ [深入浅出 eBPF](https://www.ebpf.top/) : **eBPF** 相关的中文博客。
+ [CFC4N的博客](https://www.cnxct.com/) : 美团C哥的博客 从事反入侵安全防护产品研发 主要关于 **eBPF** 安全相关的文章

### 追踪观测相关
+ [如何用eBPF分析Golang应用](https://blog.huoding.com/2021/12/12/970) : 使用 **eBPF** 分析 **Golang**。

### 网络相关
+ [BPF and XDP Reference Guide](https://docs.cilium.io/en/stable/bpf/) : **cilium** 关于 **eBPF** 和 **XDP** 的指导手册，十分详细。

### 安全相关
+ [浅谈一下，Linux中基于eBPF的恶意利用与检测机制](https://mp.weixin.qq.com/s/-1GiCncNTqtfO_grQT7cGw) : 美团安全应急响应中心关于 **eBPF** 的恶意利用方面的分析。


## Catalogue

 + [environment](https://github.com/feifeifeimoon/ebpf-study/blob/main/ENV.md) : 使用 **Multipass** 搭建 **eBPF** 开发环境的记录。
 + [sys enter openat](https://github.com/feifeifeimoon/ebpf-study/blob/main/sys_enter_openat/README.md) : 使用 **eBPF** 追踪系统中打开文件的调用。
 + [fibonacci](https://github.com/feifeifeimoon/ebpf-study/blob/main/fibonacci/README.md) : 使用 **eBPF** 追踪 **Go** 的 **fibonacci** 程序，并解析 **UProbe** 原理。
 + [xdp drop the world](https://github.com/feifeifeimoon/ebpf-study/blob/main/xdp_drop_the_world/README.md) : **XDP** 丢弃所有包的程序。
 + [xdp blacklist](https://github.com/feifeifeimoon/ebpf-study/blob/main/xdp_drop_the_world/README.md) : 在上一个的基础上 通过 **eBFP** 的 **Map** 实现 **IP** 黑名单。

## TODO
+ **cilium** 相关的详细分析，如何实现的跨主机联通
+ **eBPF** 实现的容器逃逸相关原理