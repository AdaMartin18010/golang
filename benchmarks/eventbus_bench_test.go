// EventBus 性能基准测试
package benchmarks

import (
	"context"
	"testing"

	"github.com/yourusername/golang/pkg/eventbus"
)

// BenchmarkEventBus_Publish 测试事件发布性能
func BenchmarkEventBus_Publish(b *testing.B) {
	eb := eventbus.NewEventBus(1000)
	defer eb.Stop()

	// 预创建订阅者
	for i := 0; i < 10; i++ {
		eb.Subscribe("test.event", func(ctx context.Context, e eventbus.Event) error {
			return nil
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		eb.Publish(eventbus.NewEvent("test.event", i))
	}
}

// BenchmarkEventBus_PublishParallel 测试并发发布性能
func BenchmarkEventBus_PublishParallel(b *testing.B) {
	eb := eventbus.NewEventBus(1000)
	defer eb.Stop()

	for i := 0; i < 10; i++ {
		eb.Subscribe("test.event", func(ctx context.Context, e eventbus.Event) error {
			return nil
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			eb.Publish(eventbus.NewEvent("test.event", 1))
		}
	})
}
