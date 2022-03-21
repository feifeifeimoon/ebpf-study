#include "vmlinux.h"
#include <bpf/bpf_helpers.h>

SEC("xpd_drop")
int xdp_drop_the_world(struct xdp_md *ctx) {
    return XDP_DROP;
}

char __license[] SEC("license") = "GPL";