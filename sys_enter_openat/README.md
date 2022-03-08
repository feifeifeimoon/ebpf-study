# openat


## Build

```bash
$ make      
clang -target bpf -Wall -O2  -I../inc -c  -o sys_enter_openat.bpf.o sys_enter_openat.bpf.c
clang -Wall -O2 -o sys_enter_openat sys_enter_openat.c -l elf -l bpf
```


## Run
```bash
./sys_enter_openat
libbpf: elf: skipping unrecognized data section(4) .rodata.str1.1
           <...>-8726    [000] d... 133846.134355: bpf_trace_printk: [trace openat] filename='/sys/kernel/debug/tracing/trace_pipe' flags='0
             zsh-8651    [000] d... 133852.156043: bpf_trace_printk: [trace openat] filename='/dev/null' flags='258
             zsh-8651    [000] d... 133852.156839: bpf_trace_printk: [trace openat] filename='/dev/null' flags='258
             zsh-8651    [000] d... 133852.157314: bpf_trace_printk: [trace openat] filename='/dev/null' flags='833
             zsh-8651    [000] d... 133852.157399: bpf_trace_printk: [trace openat] filename='/dev/null' flags='833
             zsh-8651    [000] d... 133852.157499: bpf_trace_printk: [trace openat] filename='/usr/share/locale/C.UTF-8/LC_MESSAGES/libc.mo' flags='0
             zsh-8651    [000] d... 133852.157532: bpf_trace_printk: [trace openat] filename='/usr/share/locale/C.utf8/LC_MESSAGES/libc.mo' flags='0
             zsh-8651    [000] d... 133852.157543: bpf_trace_printk: [trace openat] filename='/usr/share/locale/C/LC_MESSAGES/libc.mo' flags='0
             zsh-8651    [000] d... 133852.157551: bpf_trace_printk: [trace openat] filename='/usr/share/locale-langpack/C.UTF-8/LC_MESSAGES/libc.mo' flags='0
             zsh-8651    [000] d... 133852.157563: bpf_trace_printk: [trace openat] filename='/usr/share/locale-langpack/C.utf8/LC_MESSAGES/libc.mo' flags='0
             zsh-8651    [000] d... 133852.157570: bpf_trace_printk: [trace openat] filename='/usr/share/locale-langpack/C/LC_MESSAGES/libc.mo' flags='0
```


## Note

一个简单的 **eBPF** 程序，用来追踪系统中打开文件的调用。
首先查看系统中关于 **open** 的 **tracepoint** 有哪些
```bash
$ ls /sys/kernel/debug/tracing/events/syscalls | grep sys_enter_open
sys_enter_open_by_handle_at
sys_enter_open_tree
sys_enter_openat
sys_enter_openat2
```
可以看到有四个都是 open 相关的但根据名字就可以把前两个排除，剩下 `sys_enter_openat` 和 `sys_enter_openat2`。 这时候我们不知道会调到哪个函数，那我们就借用 **strace** 工具来具体看一下
```bash
$ strace ls    
...
openat(AT_FDCWD, "/etc/ld.so.cache", O_RDONLY|O_CLOEXEC) = 3
...
+++ exited with 0 +++
```
因为我的环境中装了 **zsh** 和 **oh-my-zsh** 会输出一大堆这里省略了许多，但可以看到其中调用了 `openat`