package ebpf

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollectorStructure 测试 Collector 结构体定义
func TestCollectorStructure(t *testing.T) {
	// 测试 Collector 结构体可以被实例化
	c := &Collector{}
	assert.NotNil(t, c, "Collector 实例不应为 nil")
}

// TestNewCollector 测试 NewCollector 构造函数
func TestNewCollector(t *testing.T) {
	// 当前实现返回未实现错误
	collector, err := NewCollector()

	// 由于当前是占位符实现，应返回错误
	require.Error(t, err, "NewCollector 应返回未实现错误")
	assert.Contains(t, err.Error(), "not implemented", "错误信息应包含 'not implemented'")
	// 注意：当前实现即使返回错误也会返回非 nil 的 collector
	assert.NotNil(t, collector, "当前实现返回错误时 collector 不为 nil")
}

// TestCollector_Start 测试 Start 方法
func TestCollector_Start(t *testing.T) {
	c := &Collector{}
	ctx := context.Background()

	err := c.Start(ctx)
	assert.NoError(t, err, "Start 方法不应返回错误")
}

// TestCollector_Stop 测试 Stop 方法
func TestCollector_Stop(t *testing.T) {
	c := &Collector{}

	err := c.Stop()
	assert.NoError(t, err, "Stop 方法不应返回错误")
}

// TestCollector_Collect 测试 Collect 方法
func TestCollector_Collect(t *testing.T) {
	c := &Collector{}

	data, err := c.Collect()

	// 当前是占位符实现，应返回错误
	require.Error(t, err, "Collect 应返回未实现错误")
	assert.Nil(t, data, "返回错误时 data 应为 nil")
	assert.Contains(t, err.Error(), "not implemented", "错误信息应包含 'not implemented'")
}

// TestCollectorLifecycle 测试收集器完整生命周期
func TestCollectorLifecycle(t *testing.T) {
	c := &Collector{}
	ctx := context.Background()

	// 启动
	err := c.Start(ctx)
	require.NoError(t, err, "启动不应失败")

	// 停止
	err = c.Stop()
	require.NoError(t, err, "停止不应失败")
}

// TestCollector_MethodsAfterStop 测试停止后的方法调用
func TestCollector_MethodsAfterStop(t *testing.T) {
	c := &Collector{}
	ctx := context.Background()

	// 先停止
	_ = c.Stop()

	// 再次调用 Start
	err := c.Start(ctx)
	assert.NoError(t, err, "停止后再次启动不应失败")
}
