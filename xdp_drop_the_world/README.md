# xdp drop the world

```bash

$ docker run -d -P --name iamfoo traefik/whoam

$ docker inspect --format '{{ .NetworkSettings.Networks.bridge.IPAddress }}'  iamfoo
172.17.0.3

$ curl http://172.17.0.3             
Hostname: 1370a61208df
IP: 127.0.0.1
IP: 172.17.0.3
RemoteAddr: 172.17.0.1:55848
GET / HTTP/1.1
Host: 172.17.0.3
User-Agent: curl/7.61.1
Accept: */*

$ clang -target bpf -Wall -O2 -I../inc -c xdp_drop_the_world.bpf.c

$ ip link set dev veth055b6a7 xdp obj xdp_drop_the_world.bpf.o sec xpd_drop

# 此时curl会卡住 证明加载成功
$ curl http://172.17.0.3 

# 卸载    
$ ip link set dev veth055b6a7 xdp off

```