ARCH=$(shell go env GOARCH)
CURRENT_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

main: main.go
	go build -gcflags '-l' $<

.PYONY: bpftrace clean
bpftrace:
ifeq (amd64,${ARCH})
	bpftrace -e 'uprobe:${CURRENT_PATH}/main:main.Fibonacci { printf("%d \n", reg("ax")); }  uretprobe:${CURRENT_PATH}/main:main.Fibonacci {printf("ret %d\n", retval)}'
else
	bpftrace -e 'uprobe:${CURRENT_PATH}/main:main.Fibonacci { printf("%d \n", arg0); }  uretprobe:${CURRENT_PATH}/main:main.Fibonacci {printf("ret %d\n", retval)}'
endif

clean:
	rm main