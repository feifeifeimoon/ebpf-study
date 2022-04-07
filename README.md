# ebpf-study

![logo](https://github.com/feifeifeimoon/ebpf-study/blob/main/img/logo.png?raw=true)
eBPF - extended Berkeley Packet Filter

记录 **eBPF** 学习过程中的资源和动手实现的一些例子🌰。

## Website

+ [ebpf.io](https://ebpf.io/) : **eBPF** 官网，由  **Linux Foundation** 主持。
+ [BPF and XDP Reference Guide](https://docs.cilium.io/en/stable/bpf/) : **cilium** 关于 **BPF** 和 **XDP** 的指导手册，十分详细。
+ [ARTHURCHIAO'S BLOG](https://arthurchiao.art/articles-zh/) : 携程大佬的博客，里面有很多关于 **eBPF** 文章的翻译，包括上面的[指导手册](https://arthurchiao.art/blog/cilium-bpf-xdp-reference-guide-zh/) 。
+ [深入浅出 eBPF](https://www.ebpf.top/) : **eBPF** 相关的中文博客。


## Catalogue

 + [environment](https://github.com/feifeifeimoon/ebpf-study/blob/main/ENV.md) : 使用 **Multipass** 搭建 **eBPF** 开发环境的记录。
 + [sys enter openat](https://github.com/feifeifeimoon/ebpf-study/blob/main/sys_enter_openat/README.md) : 使用 **eBPF** 追踪系统中打开文件的调用。
 + [fibonacci](https://github.com/feifeifeimoon/ebpf-study/blob/main/fibonacci/README.md) : 使用 **eBPF** 追踪 **Go** 的 **fibonacci** 程序，并解析 **UProbe** 原理。
 + [xdp drop the world](https://github.com/feifeifeimoon/ebpf-study/blob/main/xdp_drop_the_world/README.md) : **XDP** 丢弃所有包的程序。
 + [xdp blacklist](https://github.com/feifeifeimoon/ebpf-study/blob/main/xdp_drop_the_world/README.md) : 在上一个的基础上 通过 **eBFP** 的 **Map** 实现 **IP** 黑名单。