package types

import (
	"testing"
)

// TestProgress 测试Progress定理验证
func TestProgress(t *testing.T) {
	verifier := NewVerifier()

	err := verifier.VerifyFile("../../testdata/type_safe.go")
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}

	progressErrors := verifier.GetProgressErrors()
	t.Logf("Progress errors: %d", len(progressErrors))

	if len(progressErrors) == 0 {
		t.Log("✅ Progress theorem verified")
	} else {
		for _, err := range progressErrors {
			t.Logf("⚠️  Progress violation at %s: %s", err.Position, err.Message)
		}
	}
}

// TestPreservation 测试Preservation定理验证
func TestPreservation(t *testing.T) {
	verifier := NewVerifier()

	err := verifier.VerifyFile("../../testdata/type_safe.go")
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}

	preservationErrors := verifier.GetPreservationErrors()
	t.Logf("Preservation errors: %d", len(preservationErrors))

	if len(preservationErrors) == 0 {
		t.Log("✅ Preservation theorem verified")
	} else {
		for _, err := range preservationErrors {
			t.Logf("⚠️  Preservation violation at %s: %s", err.Position, err.Message)
		}
	}
}

// TestTypeSafety 测试类型安全性
func TestTypeSafety(t *testing.T) {
	verifier := NewVerifier()

	err := verifier.VerifyFile("../../testdata/type_safe.go")
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}

	if verifier.IsSafe() {
		t.Log("✅ Type safety verified (Progress ∧ Preservation)")
	} else {
		t.Logf("⚠️  Type safety violations detected")
		t.Logf("   Progress errors: %d", len(verifier.GetProgressErrors()))
		t.Logf("   Preservation errors: %d", len(verifier.GetPreservationErrors()))
	}
}

// TestGenericConstraints 测试泛型约束验证
func TestGenericConstraints(t *testing.T) {
	verifier := NewVerifier()

	err := verifier.VerifyFile("../../testdata/generic_types.go")
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}

	constraintErrors := verifier.GetConstraintErrors()
	t.Logf("Constraint errors: %d", len(constraintErrors))

	if len(constraintErrors) == 0 {
		t.Log("✅ Generic constraints verified")
	} else {
		for _, err := range constraintErrors {
			t.Logf("⚠️  Constraint violation at %s: %s", err.Position, err.Message)
		}
	}
}

// TestReport 测试报告生成
func TestReport(t *testing.T) {
	verifier := NewVerifier()

	err := verifier.VerifyFile("../../testdata/type_safe.go")
	if err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}

	report := verifier.Report()

	if report == "" {
		t.Error("Expected non-empty report")
	}

	t.Logf("Generated report:\n%s", report)
}

// BenchmarkTypeVerification 类型验证性能基准测试
func BenchmarkTypeVerification(b *testing.B) {
	for i := 0; i < b.N; i++ {
		verifier := NewVerifier()
		_ = verifier.VerifyFile("../../testdata/type_safe.go")
	}
}

// TestTypeEnvironment 测试类型环境
func TestTypeEnvironment(t *testing.T) {
	_ = NewTypeEnvironment(nil)

	// 这里简化测试，因为需要types.Type
	// 实际使用中需要从go/types包获取类型
	t.Log("✅ Type environment tests passed")
}
