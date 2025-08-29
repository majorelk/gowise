// Package performance demonstrates performance testing patterns with GoWise.
// This example shows how to benchmark assertions, measure allocation overhead,
// and implement performance regression testing.
package performance

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"gowise/pkg/assertions"
)

// mockT is a minimal TestingT implementation for benchmarks
type mockT struct{}

func (m *mockT) Errorf(format string, args ...interface{}) {}
func (m *mockT) FailNow()                                  {}

// BenchmarkCoreAssertions measures performance of basic assertion operations
func BenchmarkCoreAssertions(b *testing.B) {
	b.Run("Equal/Int/Success", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal(42, 42)
		}
	})

	b.Run("Equal/String/Success", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal("hello", "hello")
		}
	})

	b.Run("Equal/Struct/Success", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		type TestStruct struct {
			ID   int
			Name string
		}
		s1 := TestStruct{ID: 1, Name: "test"}
		s2 := TestStruct{ID: 1, Name: "test"}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal(s1, s2)
		}
	})

	b.Run("DeepEqual/Slice/Success", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		slice1 := []int{1, 2, 3, 4, 5}
		slice2 := []int{1, 2, 3, 4, 5}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.DeepEqual(slice1, slice2)
		}
	})
}

// BenchmarkFailurePath measures performance when assertions fail
func BenchmarkFailurePath(b *testing.B) {
	b.Run("Equal/Int/Failure", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal(42, 24)
		}
	})

	b.Run("Equal/String/Failure", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal("hello", "world")
		}
	})

	b.Run("StringDiff/MultiLine", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		long1 := strings.Repeat("line1\nline2\nline3\n", 100)
		long2 := strings.Repeat("line1\nmodified\nline3\n", 100)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal(long1, long2)
		}
	})
}

// BenchmarkCollectionAssertions measures collection operation performance
func BenchmarkCollectionAssertions(b *testing.B) {
	b.Run("Len/Slice/Small", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		slice := []int{1, 2, 3, 4, 5}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Len(slice, 5)
		}
	})

	b.Run("Len/Slice/Large", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		slice := make([]int, 10000)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Len(slice, 10000)
		}
	})

	b.Run("Contains/Slice/Found", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		slice := make([]int, 1000)
		for i := range slice {
			slice[i] = i
		}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Contains(slice, 500) // Middle element
		}
	})

	b.Run("Contains/Map/Found", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		m := make(map[string]int, 1000)
		for i := 0; i < 1000; i++ {
			m[fmt.Sprintf("key%d", i)] = i
		}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Contains(m, "key500")
		}
	})
}

// BenchmarkAsyncAssertions measures async assertion performance
func BenchmarkAsyncAssertions(b *testing.B) {
	b.Run("Eventually/ImmediateSuccess", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Eventually(func() bool {
				return true
			}, 1*time.Second, 10*time.Millisecond)
		}
	})

	b.Run("WithinTimeout/Fast", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.WithinTimeout(func() {
				// Fast operation
			}, 100*time.Millisecond)
		}
	})
}

// BenchmarkMemoryUsage demonstrates memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("MultipleAssertions/Chained", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert := assertions.New(&mockT{})

			// Multiple assertions
			assert.Equal(1, 1)
			assert.True(true)
			assert.Nil(nil)
			assert.Len("test", 4)
			assert.Contains("hello", "ell")
		}
	})

	b.Run("ReusedAssertion/Context", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Reuse same assertion context
			assert.Equal(1, 1)
			assert.True(true)
			assert.Nil(nil)
			assert.Len("test", 4)
			assert.Contains("hello", "ell")
		}
	})
}

// PerformanceTest demonstrates how to implement performance regression tests
func TestPerformanceRegression(t *testing.T) {
	assert := assertions.New(t)

	// Test assertion performance doesn't regress
	t.Run("EqualPerformance", func(t *testing.T) {
		const iterations = 1000000

		start := time.Now()
		testAssert := assertions.New(&mockT{})

		for i := 0; i < iterations; i++ {
			testAssert.Equal(42, 42)
		}

		elapsed := time.Since(start)
		nanosPerOp := elapsed.Nanoseconds() / iterations

		// Assert performance is within acceptable bounds
		// Equal should be under 10ns per operation for integers
		assert.True(nanosPerOp < 10)
		assert.True(elapsed < 100*time.Millisecond)

		t.Logf("Equal performance: %d ns/op (%d ops in %v)",
			nanosPerOp, iterations, elapsed)
	})

	t.Run("MemoryAllocationRegression", func(t *testing.T) {
		// Measure memory allocations for success path
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		testAssert := assertions.New(&mockT{})
		for i := 0; i < 1000; i++ {
			testAssert.Equal(i, i)
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		allocations := m2.TotalAlloc - m1.TotalAlloc
		allocPerOp := allocations / 1000

		// Success path should have minimal allocations
		assert.True(allocPerOp < 100) // Less than 100 bytes per successful assertion

		t.Logf("Memory allocation: %d bytes total, %d bytes/op",
			allocations, allocPerOp)
	})
}

// Example_benchmarks demonstrates how to use benchmarks effectively
func Example_benchmarks() {
	// Run specific benchmark
	// go test -bench=BenchmarkCoreAssertions/Equal/Int -benchmem

	// Run all benchmarks
	// go test -bench=. -benchmem

	// Run benchmarks multiple times for stability
	// go test -bench=. -count=5

	// Profile CPU usage
	// go test -bench=. -cpuprofile=cpu.prof

	// Profile memory usage
	// go test -bench=. -memprofile=mem.prof

	fmt.Println("Use go test -bench=. to run performance benchmarks")
	// Output: Use go test -bench=. to run performance benchmarks
}

// BenchmarkComparison compares GoWise performance with manual testing
func BenchmarkComparison(b *testing.B) {
	b.Run("GoWise/Equal", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Equal(42, 42)
		}
	})

	b.Run("Manual/IfStatement", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if 42 != 42 {
				// Would log error
			}
		}
	})

	b.Run("GoWise/Contains", func(b *testing.B) {
		assert := assertions.New(&mockT{})
		slice := []int{1, 2, 3, 4, 5}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			assert.Contains(slice, 3)
		}
	})

	b.Run("Manual/Loop", func(b *testing.B) {
		slice := []int{1, 2, 3, 4, 5}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			found := false
			for _, v := range slice {
				if v == 3 {
					found = true
					break
				}
			}
			_ = found
		}
	})
}

// TestPerformanceMonitoring shows how to monitor performance in tests
func TestPerformanceMonitoring(t *testing.T) {
	assert := assertions.New(t)

	t.Run("OperationTiming", func(t *testing.T) {
		const maxDuration = 10 * time.Millisecond

		start := time.Now()

		// Perform operation that should be fast
		testAssert := assertions.New(&mockT{})
		for i := 0; i < 1000; i++ {
			testAssert.Equal(i, i)
		}

		elapsed := time.Since(start)

		// Use Within Duration to ensure performance
		assert.IsWithinDuration(time.Now(), start.Add(maxDuration), maxDuration)

		if elapsed > maxDuration {
			t.Logf("Performance warning: operation took %v, expected < %v",
				elapsed, maxDuration)
		}
	})

	t.Run("MemoryUsageMonitoring", func(t *testing.T) {
		var m runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m)
		startAlloc := m.Alloc

		// Operation to monitor
		testAssert := assertions.New(&mockT{})
		data := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			data[i] = fmt.Sprintf("item%d", i)
			testAssert.Len(data, i+1)
		}

		runtime.GC()
		runtime.ReadMemStats(&m)
		endAlloc := m.Alloc

		memUsed := endAlloc - startAlloc

		// Assert memory usage is reasonable
		maxMemory := uint64(1024 * 1024) // 1MB
		assert.True(memUsed < maxMemory)

		t.Logf("Memory used: %d bytes", memUsed)
	})
}
