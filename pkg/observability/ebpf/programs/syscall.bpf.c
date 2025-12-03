// +build ignore
// 需要使用 clang 编译此文件

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

// 系统调用事件结构
struct syscall_event {
    __u64 timestamp;
    __u32 pid;
    __u32 tid;
    __u64 syscall;
    __u64 duration;
    __s64 ret_val;
};

// eBPF Map: 存储系统调用事件
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
} syscall_events SEC(".maps");

// eBPF Map: 存储系统调用统计
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, __u64);   // syscall id
    __type(value, __u64); // count
} syscall_stats SEC(".maps");

// eBPF Map: 存储系统调用开始时间
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, __u64);   // tid
    __type(value, __u64); // start timestamp
} syscall_start_time SEC(".maps");

// Tracepoint: sys_enter (系统调用进入)
SEC("tracepoint/raw_syscalls/sys_enter")
int trace_syscall_enter(struct trace_event_raw_sys_enter *ctx) {
    __u64 tid = bpf_get_current_pid_tgid();
    __u64 timestamp = bpf_ktime_get_ns();
    
    // 记录开始时间
    bpf_map_update_elem(&syscall_start_time, &tid, &timestamp, BPF_ANY);
    
    return 0;
}

// Tracepoint: sys_exit (系统调用退出)
SEC("tracepoint/raw_syscalls/sys_exit")
int trace_syscall_exit(struct trace_event_raw_sys_exit *ctx) {
    __u64 tid = bpf_get_current_pid_tgid();
    __u64 *start_time = bpf_map_lookup_elem(&syscall_start_time, &tid);
    
    if (!start_time) {
        return 0; // 没有找到开始时间，跳过
    }
    
    __u64 end_time = bpf_ktime_get_ns();
    __u64 duration = end_time - *start_time;
    
    // 构造事件
    struct syscall_event event = {
        .timestamp = end_time,
        .pid = tid >> 32,
        .tid = tid & 0xFFFFFFFF,
        .syscall = ctx->id,
        .duration = duration,
        .ret_val = ctx->ret,
    };
    
    // 发送事件到用户空间
    bpf_perf_event_output(ctx, &syscall_events, BPF_F_CURRENT_CPU, &event, sizeof(event));
    
    // 更新统计
    __u64 *count = bpf_map_lookup_elem(&syscall_stats, &event.syscall);
    if (count) {
        __sync_fetch_and_add(count, 1);
    } else {
        __u64 init_count = 1;
        bpf_map_update_elem(&syscall_stats, &event.syscall, &init_count, BPF_ANY);
    }
    
    // 删除开始时间
    bpf_map_delete_elem(&syscall_start_time, &tid);
    
    return 0;
}

char LICENSE[] SEC("license") = "GPL";

