CLANG ?= clang
CFLAGS += -Wall -O2
INC ?= -I../inc

.PHNOY: all
all: sys_enter_openat.bpf.o sys_enter_openat

sys_enter_openat.bpf.o: sys_enter_openat.bpf.c
	$(CLANG) -target bpf $(CFLAGS)  $(INC) -c  -o $@ $<

sys_enter_openat: sys_enter_openat.c
	$(CLANG) $(CFLAGS) -o $@ $< -l elf -l bpf

.PHONY: clean

clean:
	-rm sys_enter_openat.bpf.o sys_enter_openat