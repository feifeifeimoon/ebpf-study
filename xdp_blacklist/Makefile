CLANG ?= clang
CFLAGS += -Wall -O2
CCFlags =  -target bpf -Wall
INC ?= -I../inc -I/usr/include/aarch64-linux-gnu

.PHNOY: all
all: xdp_blacklist.o

xdp_blacklist.o: bpf/xdp_blacklist.bpf.c
	$(CLANG) -target bpf -g $(CFLAGS)  $(INC) -c  -o $@ $<


.PHONY: clean

clean:
	-rm xdp_blacklist.o