package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// EtcdClient etcd 存储客户端
type EtcdClient struct {
	client *clientv3.Client
}

// NewEtcdClient 创建 etcd 客户端
func NewEtcdClient(endpoints []string) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}
	
	return &EtcdClient{client: cli}, nil
}

// Close 关闭连接
func (e *EtcdClient) Close() error {
	return e.client.Close()
}

// AcquireLeadership 获取领导权
func (e *EtcdClient) AcquireLeadership(ctx context.Context, nodeID string, ttl time.Duration) (bool, error) {
	// 创建租约
	lease, err := e.client.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return false, err
	}
	
	// 尝试获取锁 (使用事务)
	key := "/scheduler/leader"
	
	txn := e.client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, nodeID, clientv3.WithLease(lease.ID))).
		Else(clientv3.OpGet(key))
	
	txnResp, err := txn.Commit()
	if err != nil {
		return false, err
	}
	
	if txnResp.Succeeded {
		// 自动续租
		keepAliveCh, err := e.client.KeepAlive(ctx, lease.ID)
		if err != nil {
			return false, err
		}
		
		// 启动续租 goroutine
		go func() {
			for range keepAliveCh {
				// 续租成功
			}
		}()
		
		return true, nil
	}
	
	return false, nil
}

// ReleaseLeadership 释放领导权
func (e *EtcdClient) ReleaseLeadership(ctx context.Context, nodeID string) error {
	key := "/scheduler/leader"
	_, err := e.client.Delete(ctx, key)
	return err
}

// RegisterWorker 注册工作节点
func (e *EtcdClient) RegisterWorker(ctx context.Context, worker *scheduler.Worker) error {
	key := fmt.Sprintf("/scheduler/workers/%s", worker.ID)
	data, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	
	lease, err := e.client.Grant(ctx, 30) // 30秒 TTL
	if err != nil {
		return err
	}
	
	_, err = e.client.Put(ctx, key, string(data), clientv3.WithLease(lease.ID))
	return err
}

// GetWorkers 获取所有工作节点
func (e *EtcdClient) GetWorkers(ctx context.Context) ([]scheduler.Worker, error) {
	resp, err := e.client.Get(ctx, "/scheduler/workers/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	
	var workers []scheduler.Worker
	for _, kv := range resp.Kvs {
		var worker scheduler.Worker
		if err := json.Unmarshal(kv.Value, &worker); err != nil {
			continue
		}
		workers = append(workers, worker)
	}
	
	return workers, nil
}

// SaveTask 保存任务
func (e *EtcdClient) SaveTask(ctx context.Context, task scheduler.Task) error {
	key := fmt.Sprintf("/scheduler/tasks/%s", task.ID)
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	
	_, err = e.client.Put(ctx, key, string(data))
	return err
}

// GetTasks 获取任务列表
func (e *EtcdClient) GetTasks(ctx context.Context, status string, limit int) ([]scheduler.Task, error) {
	resp, err := e.client.Get(ctx, "/scheduler/tasks/", clientv3.WithPrefix(), clientv3.WithLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	
	var tasks []scheduler.Task
	for _, kv := range resp.Kvs {
		var task scheduler.Task
		if err := json.Unmarshal(kv.Value, &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}
