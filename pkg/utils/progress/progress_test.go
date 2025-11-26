package progress

import (
	"testing"
	"time"
)

func TestNewProgressBar(t *testing.T) {
	pb := NewProgressBar(100)
	if pb.Total() != 100 {
		t.Errorf("Expected total 100, got %d", pb.Total())
	}
}

func TestProgressBarAdd(t *testing.T) {
	pb := NewProgressBar(100)
	pb.Add(10)
	if pb.Current() != 10 {
		t.Errorf("Expected current 10, got %d", pb.Current())
	}
}

func TestProgressBarSet(t *testing.T) {
	pb := NewProgressBar(100)
	pb.Set(50)
	if pb.Current() != 50 {
		t.Errorf("Expected current 50, got %d", pb.Current())
	}
}

func TestProgressBarIncrement(t *testing.T) {
	pb := NewProgressBar(100)
	pb.Increment()
	if pb.Current() != 1 {
		t.Errorf("Expected current 1, got %d", pb.Current())
	}
}

func TestProgressBarPercent(t *testing.T) {
	pb := NewProgressBar(100)
	pb.Set(50)
	percent := pb.Percent()
	if percent != 50.0 {
		t.Errorf("Expected percent 50.0, got %f", percent)
	}
}

func TestProgressBarFinish(t *testing.T) {
	pb := NewProgressBar(100)
	pb.Finish()
	if pb.Current() != pb.Total() {
		t.Errorf("Expected current to equal total")
	}
}

func TestSimpleProgressBar(t *testing.T) {
	spb := NewSimpleProgressBar(100)
	spb.Add(10)
	if spb.current != 10 {
		t.Errorf("Expected current 10, got %d", spb.current)
	}
}

func TestSpinner(t *testing.T) {
	spinner := NewSpinner("Loading")
	spinner.Start()
	time.Sleep(200 * time.Millisecond)
	spinner.Stop()
}
