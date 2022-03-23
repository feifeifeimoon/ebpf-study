//#include "vmlinux.h"
#include <arpa/inet.h>
#include <sys/types.h>


#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <bpf/bpf_helpers.h>


/*
 * We define u64 as uint64_t for every architecture
 * so that we can print it with "%"PRIx64 without getting warnings.
 *
 * typedef __u64 u64;
 * typedef __s64 s64;
 */
typedef uint64_t u64;
typedef int64_t s64;

typedef __u32 u32;
typedef __s32 s32;

typedef __u16 u16;
typedef __s16 s16;


struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, u32);
    __type(value, u64);
    __uint(map_flags, BPF_F_NO_PREALLOC);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} xdp_blacklist_map SEC(".maps");


#define DEBUG 1
#ifdef  DEBUG
/* Only use this for debug output. Notice output from bpf_trace_printk()
 * end-up in /sys/kernel/debug/tracing/trace_pipe
 */
#define bpf_debug(fmt, ...)                        \
        ({                            \
            char ____fmt[] = fmt;                \
            bpf_trace_printk(____fmt, sizeof(____fmt),    \
                     ##__VA_ARGS__);            \
        })
#else
#define bpf_debug(fmt, ...) { } while (0)
#endif

static __always_inline
int parse_eth(struct ethhdr *eth, void *data_end, u16 *eth_proto, u64 *l3_offset) {
    u16 eth_type;
    u64 offset;

    offset = sizeof(*eth);
    if ((void *) eth + offset > data_end)
        return 0;

    eth_type = eth->h_proto;
    bpf_debug("Debug: eth_type:0x%x\n", ntohs(eth_type));

    /* Skip non 802.3 Ethertypes */
    if (ntohs(eth_type) < ETH_P_802_3_MIN)
        return 0;


    *eth_proto = ntohs(eth_type);
    *l3_offset = offset;
    return 1;
}


static __always_inline
u32 parse_ipv4(struct xdp_md *ctx, u64 l3_offset) {
    void *data_end = (void *) (long) ctx->data_end;
    void *data = (void *) (long) ctx->data;
    struct iphdr *iph = data + l3_offset;
    u64 *value;
    u32 ip_src; /* type need to match map */

    /* Hint: +1 is sizeof(struct iphdr) */
    if (iph + 1 > data_end) {
        bpf_debug("Invalid IPv4 packet: L3off:%llu\n", l3_offset);
        return XDP_ABORTED;
    }
    /* Extract key */
    ip_src = iph->saddr;
    //ip_src = ntohl(ip_src); // ntohl does not work for some reason!?!

    bpf_debug("Valid IPv4 packet: raw saddr:0x%x %d\n", ip_src, ip_src);

    value = bpf_map_lookup_elem(&xdp_blacklist_map, &ip_src);
    if (value) {
        /* Don't need __sync_fetch_and_add(); as percpu map */
        *value += 1; /* Keep a counter for drop matches */
        return XDP_DROP;
    }
//    } else {
//        u64 init = 0;
//        bpf_map_update_elem(&xdp_blacklist_map, &ip_src, &init, BPF_ANY);
//    }

    return XDP_PASS;
}


static __always_inline
u32 handle_eth_protocol(struct xdp_md *ctx, u16 eth_proto, u64 l3_offset) {
    switch (eth_proto) {
        case ETH_P_IP:
            return parse_ipv4(ctx, l3_offset);
        case ETH_P_IPV6: /* Not handler for IPv6 yet*/
        case ETH_P_ARP:  /* Let OS handle ARP */
            /* Fall-through */
        default:
            bpf_debug("Not handling eth_proto:0x%x\n", eth_proto);
            return XDP_PASS;
    }
    return XDP_PASS;
}


SEC("xdp_prog")
int xdp_blacklist_prog(struct xdp_md *ctx) {
    void *data_end = (void *) (long) ctx->data_end;
    void *data = (void *) (long) ctx->data;
    struct ethhdr *eth = data;
    u16 eth_proto = 0;
    u64 l3_offset = 0;

    if (!(parse_eth(eth, data_end, &eth_proto, &l3_offset))) {
        bpf_debug("Cannot parse L2: L3off:%llu proto:0x%x\n", l3_offset, eth_proto);
        return XDP_PASS; /* Skip */
    }


    return handle_eth_protocol(ctx, eth_proto, l3_offset);
}

// cannot call GPL-restricted function from non-GPL compatible program
char __license[] SEC("license") = "GPL";
