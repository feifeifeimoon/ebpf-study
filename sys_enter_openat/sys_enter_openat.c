#include <stdio.h>
#include <unistd.h>
#include <bpf/libbpf.h>
#include <fcntl.h>

#define DEBUGFS "/sys/kernel/debug/tracing/"

void read_trace_pipe(void) {
    int trace_fd;

    trace_fd = open(DEBUGFS "trace_pipe", O_RDONLY, 0);
    if (trace_fd < 0)
        return;

    while (1) {
        static char buf[4096];
        ssize_t sz;

        sz = read(trace_fd, buf, sizeof(buf) - 1);
        if (sz > 0) {
            buf[sz] = 0;
            puts(buf);
        } else {
            printf("read err, ret=%ld\n", sz);
            return;
        }
    }
}

int main(int argc, char **argv) {
    struct bpf_link *link = NULL;
    struct bpf_program *prog;
    struct bpf_object *obj;
    char filename[256];

    snprintf(filename, sizeof(filename), "sys_enter_openat.bpf.o");
    obj = bpf_object__open_file(filename, NULL);
    if (libbpf_get_error(obj)) {
        fprintf(stderr, "ERROR: opening BPF object file failed\n");
        return 0;
    }

    prog = bpf_object__find_program_by_name(obj, "tracepoint__syscalls__sys_enter_openat");
    if (!prog) {
        fprintf(stderr, "ERROR: finding a prog in obj file failed\n");
        goto cleanup;
    }

    /* load BPF program */
    if (bpf_object__load(obj)) {
        fprintf(stderr, "ERROR: loading BPF object file failed\n");
        goto cleanup;
    }

    link = bpf_program__attach(prog);
    if (libbpf_get_error(link)) {
        fprintf(stderr, "ERROR: bpf_program__attach failed\n");
        link = NULL;
        goto cleanup;
    }


    read_trace_pipe();

    cleanup:
    bpf_link__destroy(link);
    bpf_object__close(obj);
    return 0;
}