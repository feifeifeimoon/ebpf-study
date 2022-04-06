# openat

## Introduction

一个简单的 **eBPF** 程序，用来追踪系统中打开文件的调用。
首先查看系统中关于 **open** 的 **tracepoint** 有哪些
```bash
$ ls /sys/kernel/debug/tracing/events/syscalls | grep sys_enter_open
sys_enter_open_by_handle_at
sys_enter_open_tree
sys_enter_openat
sys_enter_openat2
```
可以看到有四个都是 **open** 相关的但根据名字就可以把前两个排除，剩下 `sys_enter_openat` 和 `sys_enter_openat2`。 这时候我们不知道会调到哪个函数，那我们就借用 **strace** 工具来具体看一下
```bash
$ strace ls    
...
openat(AT_FDCWD, "/etc/ld.so.cache", O_RDONLY|O_CLOEXEC) = 3
...
+++ exited with 0 +++
```
因为我的环境中装了 **zsh** 和 **oh-my-zsh** 会输出一大堆这里省略了许多，但可以看到其中调用了 `openat`。

找到了 **tracepoint** 接着就查看可以得到的参数有哪些

```shell
 cat /sys/kernel/debug/tracing/events/syscalls/sys_enter_openat2/format 
name: sys_enter_openat2
ID: 548
format:
        field:unsigned short common_type;       offset:0;       size:2; signed:0;
        field:unsigned char common_flags;       offset:2;       size:1; signed:0;
        field:unsigned char common_preempt_count;       offset:3;       size:1; signed:0;
        field:int common_pid;   offset:4;       size:4; signed:1;

        field:int __syscall_nr; offset:8;       size:4; signed:1;
        field:int dfd;  offset:16;      size:8; signed:0;
        field:const char * filename;    offset:24;      size:8; signed:0;
        field:struct open_how * how;    offset:32;      size:8; signed:0;
        field:size_t usize;     offset:40;      size:8; signed:0;

print fmt: "dfd: 0x%08lx, filename: 0x%08lx, how: 0x%08lx, usize: 0x%08lx", ((unsigned long)(REC->dfd)), ((unsigned long)(REC->filename)), ((unsigned long)(REC->how)), ((unsigned long)(REC->usize))
```

对应到 `struct trace_event_raw_sys_enter`结构体，`args[0]` 就是 `dfd` 的值，  `args[1]` 就是 `filename` 以此类推
## Code

```C
SEC("tracepoint/syscalls/sys_enter_openat")
int tracepoint__syscalls__sys_enter_openat(struct trace_event_raw_sys_enter *ctx) {

    const char *fname = (const char *) ctx->args[1];
    int flags = (int) ctx->args[2];
    char fmt[] = "[trace openat] filename='%s' flags='%d";
    bpf_trace_printk(fmt, sizeof(fmt), fname, flags);
    return 0;
}

```

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
