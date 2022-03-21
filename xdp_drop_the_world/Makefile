CLANG ?= clang
CFLAGS += -Wall -O2
CCFlags =  -target bpf -Wall
INC ?= -I../inc

.PHNOY: all
all: xdp_drop_the_world.bpf.o

xdp_drop_the_world.bpf.o: xdp_drop_the_world.bpf.c
	$(CLANG) -target bpf $(CFLAGS)  $(INC) -c  -o $@ $<


.PHONY: clean

clean:
	-rm xdp_drop_the_world.bpf.o