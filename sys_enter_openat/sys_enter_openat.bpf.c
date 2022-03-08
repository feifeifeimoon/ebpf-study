#include "vmlinux.h"
#include <bpf/bpf_helpers.h>

SEC("tracepoint/syscalls/sys_enter_openat")
int tracepoint__syscalls__sys_enter_openat(struct trace_event_raw_sys_enter *ctx) {

    const char *fname = (const char *) ctx->args[1];
    int flags = (int) ctx->args[2];
    char fmt[] = "[trace openat] filename='%s' flags='%d";
    bpf_trace_printk(fmt, sizeof(fmt), fname, flags);
    return 0;
}


char _license[] SEC("license") = "GPL";