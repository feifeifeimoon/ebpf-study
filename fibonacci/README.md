# 使用 eBPF 追踪 Go 程序

主要演示使用 **eBPF UProbe** 来追踪一个简单的 **Golang** 服务。并通过 **gdb** 来揭秘 **UProbe** 底层究竟用了什么 **Magic**。

## UProbe 用法
先来看通过 **UProbe** 可以做什么。下面是一个 **demon** 程序，非常简单核心就是将 **Fibonacci** 通过 **HTTP** 向外提供服务。

```go
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func Fibonacci(n int) int {
	if n < 2 {
		return n
	}
	p, q, r := 0, 0, 1
	for i := 2; i <= n; i++ {
		p = q
		q = r
		r = p + q
	}
	return r
}

func main() {
	r := gin.Default()

	r.GET("/fibonacci", func(c *gin.Context) {
		num := c.DefaultQuery("num", "0")
		x, err := strconv.Atoi(num)
		if err != nil || x <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The input must be an integer greater than zero"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ret": Fibonacci(x)})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}

```
用起来也非常简单，客户端通过 **http** 请求就可以得到 **Fibonacci** 的结果。
![demo](https://github.com/feifeifeimoon/ebpf-study/blob/main/img/fibonacci-demo.gif?raw=true)

接着就是 **eBPF UProbe** 出场的时候了。假如这时候我们想要查看每次调用 **Fibonacci** 函数时的入参和得到结果是什么，但又希望修改代码重新编译。这里使用 **bpftrace** 来实现这个目的

![bpftrace](https://github.com/feifeifeimoon/ebpf-study/blob/main/img/fibonacci-bpftrace.gif?raw=true)

可以看到 我们并没有去修改服务端程序也没有改什么配置，仅仅通过一条命令，就成功的追踪到了 **Fibonacci** 函数的参数和返回值。

## UProbe进一步分析

接着我们来看一下这一行命令做了什么：

```Bash
$ bpftrace -e 'uprobe:/root/ebpf-study/fibonacci/main:main.Fibonacci { printf("%d \n", reg("ax")); }  uretprobe:/root/ebpf-study/fibonacci/main:main.Fibonacci {printf("ret %d\n", retval)}'
```

首先 [**bpftrace**](https://github.com/iovisor/bpftrace) 是一个可以通过命令行就快速使用 **eBPF** 的工具，可以省去编写 **eBPF** 代码和加载 **eBPF** 的代码。
+ **-e** 可以指定要执行的程序。

后面的内容我们可以分为两部分来看：

第一部分：

```Bash
uprobe:/root/ebpf-study/fibonacci/main:main.Fibonacci { printf("%d \n", reg("ax")); }
```

+ 首先指定使用 **uprobe** 探针，**uprobe** 探针会在函数执行前运行。
+ `/root/daily/ebpf-test/fibonacci/main` 是我编译得到的可执行文件的位置。
+ `main.Fibonacci` 是指定要附加 **uprobe** 探针的函数
+ `printf("%d \n", reg("ax"));` 就是在探针触发时要执行的内容，这里就是根据 `Golang` 的调用约定将参数打印出来。（注意 **arm** 环境和 **amd** 环境也不相同 具体可以查看 [**Makefile**](https://github.com/feifeifeimoon/ebpf-study/blob/main/fibonacci/Makefile) ）

关于为什么是 reg("ax")，推荐阅读 [如何用eBPF分析Golang应用](https://blog.huoding.com/2021/12/12/970) 。

第二部分：

```Bash

uretprobe:/root/ebpf-test/fibonacci/main:main.Fibonacci {printf("ret %d\n", retval)}
```

第二部分其实和第一部分类似，首先指定了 **uretprobe** 探针，这个探针不同的是它会在函数返回时运行。`/root/ebpf-test/fibonacci/main:main.Fibonacci` 这里同样是可执行文件位置和要添加探针的函数，`printf("ret %d\n", retval)` 这一句就是将函数的返回值打印出来。具体参照 [bpftrace guide](https://github.com/iovisor/bpftrace/blob/master/docs/reference_guide.md#1-builtins-1)


## UProbe原理


当时的疑惑主要是这个探针是如何插入的，首先在用户态函数互相调用过程中应该不会引起用户态和内核态的切换，其次我们在编译程序时只是添加了` -gcflags '-l'` 取消函数内联而已并没有添加其它编译选项。

经过网上查询发现这样一段话（大概意思是这样）

> 其实当我们要 attach 一个 Probe 时，会将原始方法的入口点替换成断点指令，这个指令是和 CPU 架构相关的，比如在 i386 上就是 int 3 指令，在 ARM 环境上是设置一个未定义的指令。 替换指令后当再次执行到这个函数时就会触发一个 trap，然后再 trap 的处理流程中就会调用注册的 Probe 函数。

为了证实这段话，通过 **GDB** 调试了一下上面的 **Fibonacci** 函数

