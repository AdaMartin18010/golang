package sampling

import (
	"context"
	"testing"
	"time"
)

func TestAlwaysSampler(t *testing.T) {
	sampler := NewAlwaysSampler()

	if !sampler.ShouldSample(context.Background()) {
		t.Error("AlwaysSampler should always return true")
	}

	if sampler.SampleRate() != 1.0 {
		t.Errorf("Expected sample rate 1.0, got %f", sampler.SampleRate())
	}
}

func TestNeverSampler(t *testing.T) {
	sampler := NewNeverSampler()

	if sampler.ShouldSample(context.Background()) {
		t.Error("NeverSampler should always return false")
	}

	if sampler.SampleRate() != 0.0 {
		t.Errorf("Expected sample rate 0.0, got %f", sampler.SampleRate())
	}
}

func TestProbabilisticSampler(t *testing.T) {
	sampler, err := NewProbabilisticSampler(0.5)
	if err != nil {
		t.Fatalf("Failed to create sampler: %v", err)
	}

	if sampler.SampleRate() != 0.5 {
		t.Errorf("Expected sample rate 0.5, got %f", sampler.SampleRate())
	}

	// 测试多次采样，应该大致符合概率
	samples := 1000
	trueCount := 0
	for i := 0; i < samples; i++ {
		if sampler.ShouldSample(context.Background()) {
			trueCount++
		}
	}

	// 允许一定的误差范围（±10%）
	expectedMin := int(float64(samples) * 0.4)
	expectedMax := int(float64(samples) * 0.6)
	if trueCount < expectedMin || trueCount > expectedMax {
		t.Errorf("Sample rate seems incorrect: got %d/%d (expected ~%d)", trueCount, samples, samples/2)
	}

	// 测试更新采样率
	if err := sampler.UpdateRate(0.8); err != nil {
		t.Errorf("Failed to update rate: %v", err)
	}
	if sampler.SampleRate() != 0.8 {
		t.Errorf("Expected sample rate 0.8, got %f", sampler.SampleRate())
	}
}

func TestProbabilisticSampler_InvalidRate(t *testing.T) {
	_, err := NewProbabilisticSampler(1.5)
	if err == nil {
		t.Error("Expected error for invalid rate > 1.0")
	}

	_, err = NewProbabilisticSampler(-0.1)
	if err == nil {
		t.Error("Expected error for invalid rate < 0.0")
	}
}

func TestRateLimitingSampler(t *testing.T) {
	sampler, err := NewRateLimitingSampler(10.0) // 每秒最多 10 次
	if err != nil {
		t.Fatalf("Failed to create sampler: %v", err)
	}

	// 在短时间内应该允许多次
	allowed := 0
	for i := 0; i < 15; i++ {
		if sampler.ShouldSample(context.Background()) {
			allowed++
		}
		time.Sleep(10 * time.Millisecond)
	}

	// 应该允许至少 10 次
	if allowed < 10 {
		t.Errorf("Expected at least 10 allowed, got %d", allowed)
	}
}

func TestAdaptiveSampler(t *testing.T) {
	sampler, err := NewAdaptiveSampler(0.5, 0.1, 1.0)
	if err != nil {
		t.Fatalf("Failed to create sampler: %v", err)
	}

	if sampler.SampleRate() != 0.5 {
		t.Errorf("Expected initial rate 0.5, got %f", sampler.SampleRate())
	}

	// 测试根据负载调整
	adaptiveSampler := sampler.(*AdaptiveSampler)
	adaptiveSampler.AdjustForLoad(0.9) // 高负载
	if adaptiveSampler.SampleRate() != 0.1 {
		t.Errorf("Expected rate 0.1 for high load, got %f", adaptiveSampler.SampleRate())
	}

	adaptiveSampler.AdjustForLoad(0.2) // 低负载
	if adaptiveSampler.SampleRate() != 1.0 {
		t.Errorf("Expected rate 1.0 for low load, got %f", adaptiveSampler.SampleRate())
	}
}

