# xdp blacklist

## Introduction
**xdp_blacklist** 是使用 **XDP** 和 **eBPF Map** 实现的黑名单程序。可以指定将某些 **IP** 添加到黑名单中，然后通过 **XDP** 拦截黑名单中的 **IP** 。

## Building
**xdp_blacklist** only supports Linux now.

```bash
$ git clone https://github.com/feifeifeimoon/ebpf-study.git
$ cd ebpf-study/xdp_blacklist
$ make
```

## Using

```bash
$ xdp_blacklist                   
xdp_blacklist is a blacklist program implemented through ebpf

Usage:
  xdp_blacklist [command]

Available Commands:
  add         Appends a new IP address to the blacklist
  attach      Attach xdp to network device
  completion  Generate the autocompletion script for the specified shell
  delete      Delete IP from blacklist
  detach      Detach xdp from network device
  help        Help about any command
  list        Show ip in blacklist

Flags:
  -h, --help   help for xdp_blacklist

Use "xdp_blacklist [command] --help" for more information about a command.

```