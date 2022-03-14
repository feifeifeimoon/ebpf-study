# vmlinux.h

## **vmlinux**
自己编译过内核的人应该都知道，当我们编译 Linux 内核时，会输出一个称作 **vmlinux** 的文件组件，其实是一个 **ELF** 的二进制文件，包含了编译好的可启动内核。推荐[老司机带你探索内核编译系统](https://richardweiyang-2.gitbook.io/kernel-exploring/00_index)。
**vmlinux** 文件通常也会被打包在主要的 **Linux** 发行版中 `/sys/kenel/btf/vmlinux`。

## vmlinux.h 从哪来的

先说 **vmlinux.h** 是从哪来的。**vmlinux.h** 其实是使用工具生成的代码，而生成的来源其实就是 **vmlinux** 文件。使用 **bpftool** 可以快速生成 **vmlinux.h**

```bash
$ bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
$ wc -l vmlinux.h 
157682 vmlinux.h
```
> 如果第一个命令执行失败，可能因为当前内核编译时没有开始 `CONFIG_DEBUG_INFO_BTF=y` 和 `CONFIG_DEBUG_INFO=y` 两个选项，需要重新编译一下内核。[Linux - 内核编译安装](https://blog.hyperzsb.tech/linux-kernel-compile-install/)

通过 **wc** 命令可以看到生成的 **vmlinux.h** 文件足足有 157682 行。查看其中内容，发现 **vmlinux.h** 会包含运行内核中所使用的每一个类型定义，比如 `struct task_struct`。
有了 **vmlinux.h**文件，只要在程序中 `#include "vmlinux.h"`，就意味着我们的程序可以使用内核中使用的所有数据类型定义，因此 **BPF** 程序在读取相关的内存时，就可以映射成对应的类型结构按照字段进行读取。


## CO-RE (Compile Once – Run Everywhere)

由于 **vmlinux.h** 文件是由当前运行内核生成的，如果你试图将编译好的 **eBPF** 程序在另一台运行不同内核版本的机器上运行，可能会面临崩溃的窘境。这主要是因为在不同的版本中，对应数据类型的定义可能会在 **Linux** 源代码中发生变化。

但是，通过使用 **libbpf** 库提供的功能可以实现 “**CO:RE**”（一次编译，到处运行）。**libbpf** 库定义了部分宏（比如 **BPF_CORE_READ**），其可分析 **eBPF** 程序试图访问 **vmlinux.h** 中定义的类型中的哪些字段。如果访问的字段在当前内核定义的结构中发生了移动，宏 / 辅助函数会协助自动找到对应字段。因此，我们可以使用当前内核中生成的 **vmlinux.h** 头文件来编译 **eBPF** 程序，然后在不同的内核上运行它。

[BPF CO-RE reference guide](https://nakryiko.com/posts/bpf-core-reference-guide/)