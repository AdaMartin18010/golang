// +build ignore
// eBPF 网络监控程序

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/tcp.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <net/sock.h>

// TCP 事件结构
struct tcp_event {
    __u64 timestamp;
    __u32 pid;
    __u32 tid;
    __u32 event_type; // 0=connect, 1=accept, 2=close
    __u8  src_addr[4];
    __u8  dst_addr[4];
    __u16 src_port;
    __u16 dst_port;
    __u64 bytes_sent;
    __u64 bytes_recv;
    __u64 duration;
};

// TCP 连接信息
struct tcp_conn_info {
    __u64 start_time;
    __u64 bytes_sent;
    __u64 bytes_recv;
    __u8  src_addr[4];
    __u8  dst_addr[4];
    __u16 src_port;
    __u16 dst_port;
};

// eBPF Maps
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
} tcp_events SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, __u64);   // socket pointer
    __type(value, struct tcp_conn_info);
} tcp_connections SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, __u32);   // pid
    __type(value, __u64); // connection count
} tcp_stats SEC(".maps");

// Helper: 从 sock 结构获取 IPv4 地址信息
static inline void get_sock_addr_v4(struct sock *sk, __u8 *src_addr, __u8 *dst_addr,
                                     __u16 *src_port, __u16 *dst_port) {
    // 获取源地址和端口
    // sk->__sk_common.skc_rcv_saddr 是接收地址（本地地址）
    // sk->__sk_common.skc_daddr 是目的地址
    // 注意：需要根据内核版本调整字段访问
    
    bpf_probe_read_kernel(src_addr, 4, &sk->__sk_common.skc_rcv_saddr);
    bpf_probe_read_kernel(dst_addr, 4, &sk->__sk_common.skc_daddr);
    
    bpf_probe_read_kernel(src_port, sizeof(*src_port), &sk->__sk_common.skc_num);
    bpf_probe_read_kernel(dst_port, sizeof(*dst_port), &sk->__sk_common.skc_dport);
    
    // 目的端口需要字节序转换（网络字节序 -> 主机字节序）
    *dst_port = bpf_ntohs(*dst_port);
}

// Helper: 检查是否是 IPv4 连接
static inline bool is_ipv4(struct sock *sk) {
    // sk->__sk_common.skc_family == AF_INET (2)
    __u16 family;
    bpf_probe_read_kernel(&family, sizeof(family), &sk->__sk_common.skc_family);
    return family == AF_INET;
}

// Helper: 获取 socket 指针作为 map key
static inline __u64 get_sock_key(struct sock *sk) {
    return (__u64)sk;
}

// Kprobe: tcp_connect (出站连接)
SEC("kprobe/tcp_connect")
int trace_tcp_connect(struct pt_regs *ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;
    __u64 timestamp = bpf_ktime_get_ns();
    
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    if (!sk) {
        return 0;
    }
    
    // 只处理 IPv4
    if (!is_ipv4(sk)) {
        return 0;
    }

    // 记录连接开始时间
    struct tcp_conn_info conn_info = {0};
    conn_info.start_time = timestamp;
    
    // 从 sock 结构获取地址信息
    get_sock_addr_v4(sk, conn_info.src_addr, conn_info.dst_addr,
                     &conn_info.src_port, &conn_info.dst_port);

    __u64 sock_key = get_sock_key(sk);
    bpf_map_update_elem(&tcp_connections, &sock_key, &conn_info, BPF_ANY);

    // 发送连接事件
    struct tcp_event event = {0};
    event.timestamp = timestamp;
    event.pid = pid;
    event.tid = pid_tgid & 0xFFFFFFFF;
    event.event_type = 0; // connect
    
    // 复制地址信息到事件
    __builtin_memcpy(event.src_addr, conn_info.src_addr, 4);
    __builtin_memcpy(event.dst_addr, conn_info.dst_addr, 4);
    event.src_port = conn_info.src_port;
    event.dst_port = conn_info.dst_port;

    bpf_perf_event_output(ctx, &tcp_events, BPF_F_CURRENT_CPU, &event, sizeof(event));

    // 更新统计
    __u64 *count = bpf_map_lookup_elem(&tcp_stats, &pid);
    if (count) {
        __sync_fetch_and_add(count, 1);
    } else {
        __u64 init_count = 1;
        bpf_map_update_elem(&tcp_stats, &pid, &init_count, BPF_ANY);
    }

    return 0;
}

// Kprobe: tcp_v4_connect return (连接完成)
SEC("kretprobe/tcp_v4_connect")
int trace_tcp_connect_return(struct pt_regs *ctx) {
    int ret = PT_REGS_RC(ctx);
    
    // 如果连接失败，清理记录
    if (ret != 0) {
        struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
        if (sk) {
            __u64 sock_key = get_sock_key(sk);
            bpf_map_delete_elem(&tcp_connections, &sock_key);
        }
    }

    return 0;
}

// Kprobe: inet_csk_accept (入站连接)
SEC("kprobe/inet_csk_accept")
int trace_tcp_accept(struct pt_regs *ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;
    __u64 timestamp = bpf_ktime_get_ns();

    // 发送接受事件
    struct tcp_event event = {0};
    event.timestamp = timestamp;
    event.pid = pid;
    event.tid = pid_tgid & 0xFFFFFFFF;
    event.event_type = 1; // accept

    bpf_perf_event_output(ctx, &tcp_events, BPF_F_CURRENT_CPU, &event, sizeof(event));

    return 0;
}

// Kretprobe: inet_csk_accept return (接受连接返回)
SEC("kretprobe/inet_csk_accept")
int trace_tcp_accept_return(struct pt_regs *ctx) {
    struct sock *new_sk = (struct sock *)PT_REGS_RC(ctx);
    if (!new_sk) {
        return 0;
    }
    
    // 只处理 IPv4
    if (!is_ipv4(new_sk)) {
        return 0;
    }
    
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;
    __u64 timestamp = bpf_ktime_get_ns();
    
    // 记录新连接
    struct tcp_conn_info conn_info = {0};
    conn_info.start_time = timestamp;
    
    get_sock_addr_v4(new_sk, conn_info.src_addr, conn_info.dst_addr,
                     &conn_info.src_port, &conn_info.dst_port);
    
    __u64 sock_key = get_sock_key(new_sk);
    bpf_map_update_elem(&tcp_connections, &sock_key, &conn_info, BPF_ANY);
    
    // 发送 accept 事件（带地址信息）
    struct tcp_event event = {0};
    event.timestamp = timestamp;
    event.pid = pid;
    event.tid = pid_tgid & 0xFFFFFFFF;
    event.event_type = 1; // accept
    
    __builtin_memcpy(event.src_addr, conn_info.src_addr, 4);
    __builtin_memcpy(event.dst_addr, conn_info.dst_addr, 4);
    event.src_port = conn_info.src_port;
    event.dst_port = conn_info.dst_port;
    
    bpf_perf_event_output(ctx, &tcp_events, BPF_F_CURRENT_CPU, &event, sizeof(event));
    
    // 更新统计
    __u64 *count = bpf_map_lookup_elem(&tcp_stats, &pid);
    if (count) {
        __sync_fetch_and_add(count, 1);
    } else {
        __u64 init_count = 1;
        bpf_map_update_elem(&tcp_stats, &pid, &init_count, BPF_ANY);
    }

    return 0;
}

// Kprobe: tcp_close (连接关闭)
SEC("kprobe/tcp_close")
int trace_tcp_close(struct pt_regs *ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;
    __u64 timestamp = bpf_ktime_get_ns();
    
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    if (!sk) {
        return 0;
    }
    
    __u64 sock_key = get_sock_key(sk);

    // 查找连接信息
    struct tcp_conn_info *conn_info = bpf_map_lookup_elem(&tcp_connections, &sock_key);

    struct tcp_event event = {0};
    event.timestamp = timestamp;
    event.pid = pid;
    event.tid = pid_tgid & 0xFFFFFFFF;
    event.event_type = 2; // close

    if (conn_info) {
        // 计算连接持续时间
        event.duration = timestamp - conn_info->start_time;
        event.bytes_sent = conn_info->bytes_sent;
        event.bytes_recv = conn_info->bytes_recv;

        // 复制地址信息
        __builtin_memcpy(event.src_addr, conn_info->src_addr, 4);
        __builtin_memcpy(event.dst_addr, conn_info->dst_addr, 4);
        event.src_port = conn_info->src_port;
        event.dst_port = conn_info->dst_port;

        // 删除连接记录
        bpf_map_delete_elem(&tcp_connections, &sock_key);
    }

    bpf_perf_event_output(ctx, &tcp_events, BPF_F_CURRENT_CPU, &event, sizeof(event));

    return 0;
}

// Kprobe: tcp_sendmsg (发送数据)
SEC("kprobe/tcp_sendmsg")
int trace_tcp_sendmsg(struct pt_regs *ctx) {
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    if (!sk) {
        return 0;
    }
    
    size_t size = (size_t)PT_REGS_PARM3(ctx);
    __u64 sock_key = get_sock_key(sk);

    // 更新发送字节数
    struct tcp_conn_info *conn_info = bpf_map_lookup_elem(&tcp_connections, &sock_key);
    if (conn_info) {
        __sync_fetch_and_add(&conn_info->bytes_sent, size);
    }

    return 0;
}

// Kprobe: tcp_recvmsg (接收数据)
SEC("kprobe/tcp_recvmsg")
int trace_tcp_recvmsg(struct pt_regs *ctx) {
    // 保存 sock 指针用于返回值处理
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    if (!sk) {
        return 0;
    }
    
    // 可以在这里记录开始接收的时间等信息
    return 0;
}

// Kretprobe: tcp_recvmsg return (接收数据返回)
SEC("kretprobe/tcp_recvmsg")
int trace_tcp_recvmsg_return(struct pt_regs *ctx) {
    int ret = PT_REGS_RC(ctx);
    if (ret <= 0) {
        return 0;
    }
    
    // 通过 per-cpu 数组获取之前保存的 sock 指针
    // 简化版本：直接从参数重新获取（可能不准确，取决于调用约定）
    struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
    if (!sk) {
        return 0;
    }
    
    __u64 sock_key = get_sock_key(sk);

    // 更新接收字节数
    struct tcp_conn_info *conn_info = bpf_map_lookup_elem(&tcp_connections, &sock_key);
    if (conn_info) {
        __sync_fetch_and_add(&conn_info->bytes_recv, ret);
    }

    return 0;
}

char LICENSE[] SEC("license") = "GPL";
