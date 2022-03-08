project(ebpf-study C)
cmake_minimum_required(3.1.0)

include_directories(inc)

add_executable(sys_enter_openat.bpf sys_enter_openat/sys_enter_openat.bpf.c)
add_executable(sys_enter_openat sys_enter_openat/sys_enter_openat.c)