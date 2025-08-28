package assertions

import (
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkEventually measures performance of basic Eventually assertions.
func BenchmarkEventually(b *testing.B) {
	b.Run("ImmediateSuccess", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			assert.Eventually(func() bool {
				return true
			}, 1*time.Second, 10*time.Millisecond)
		}
	})

	b.Run("ShortDelay", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			var ready int32
			
			go func() {
				time.Sleep(50 * time.Millisecond)
				atomic.StoreInt32(&ready, 1)
			}()

			assert.Eventually(func() bool {
				return atomic.LoadInt32(&ready) == 1
			}, 1*time.Second, 20*time.Millisecond)
		}
	})

	b.Run("MultipleAttempts", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			var counter int32

			assert.Eventually(func() bool {
				return atomic.AddInt32(&counter, 1) >= 5
			}, 1*time.Second, 10*time.Millisecond)
		}
	})
}

// BenchmarkNever measures performance of Never assertions.
func BenchmarkNever(b *testing.B) {
	b.Run("NeverTrue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			assert.Never(func() bool {
				return false
			}, 100*time.Millisecond, 10*time.Millisecond)
		}
	})

	b.Run("CounterCheck", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			var counter int32

			assert.Never(func() bool {
				return atomic.AddInt32(&counter, 1) > 20 // Never exceeds 20
			}, 100*time.Millisecond, 5*time.Millisecond)
		}
	})
}

// BenchmarkEventuallyWith measures performance of EventuallyWith with different configurations.
func BenchmarkEventuallyWith(b *testing.B) {
	b.Run("BasicConfig", func(b *testing.B) {
		config := EventuallyConfig{
			Timeout:  500*time.Millisecond,
			Interval: 25*time.Millisecond,
		}
		
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			assert.EventuallyWith(func() bool {
				return true
			}, config)
		}
	})

	b.Run("ExponentialBackoff", func(b *testing.B) {
		config := EventuallyConfig{
			Timeout:       1*time.Second,
			Interval:      10*time.Millisecond,
			BackoffFactor: 2.0,
			MaxInterval:   100*time.Millisecond,
		}
		
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			var counter int32
			
			assert.EventuallyWith(func() bool {
				return atomic.AddInt32(&counter, 1) >= 3
			}, config)
		}
	})

	b.Run("NoBackoff", func(b *testing.B) {
		config := EventuallyConfig{
			Timeout:       500*time.Millisecond,
			Interval:      20*time.Millisecond,
			BackoffFactor: 1.0, // No backoff
		}
		
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			var counter int32
			
			assert.EventuallyWith(func() bool {
				return atomic.AddInt32(&counter, 1) >= 3
			}, config)
		}
	})
}

// BenchmarkConcurrentAssertions measures performance under concurrent load.
func BenchmarkConcurrentAssertions(b *testing.B) {
	b.Run("ParallelEventually", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				assert := New(&mockT{})
				var ready int32
				
				go func() {
					time.Sleep(30 * time.Millisecond)
					atomic.StoreInt32(&ready, 1)
				}()

				assert.Eventually(func() bool {
					return atomic.LoadInt32(&ready) == 1
				}, 500*time.Millisecond, 15*time.Millisecond)
			}
		})
	})

	b.Run("ParallelNever", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				assert := New(&mockT{})
				var counter int32

				assert.Never(func() bool {
					return atomic.AddInt32(&counter, 1) > 15
				}, 100*time.Millisecond, 10*time.Millisecond)
			}
		})
	})
}

// BenchmarkResourceUsage measures resource efficiency.
func BenchmarkResourceUsage(b *testing.B) {
	b.Run("MemoryEfficiency", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			
			// Test memory efficiency with longer operations
			var counter int32
			assert.Eventually(func() bool {
				return atomic.AddInt32(&counter, 1) >= 10
			}, 500*time.Millisecond, 25*time.Millisecond)
		}
	})

	b.Run("GoroutineCleanup", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			
			// Test that cleanup properly handles goroutines
			var ready int32
			go func() {
				time.Sleep(100 * time.Millisecond)
				atomic.StoreInt32(&ready, 1)
			}()
			
			assert.Eventually(func() bool {
				return atomic.LoadInt32(&ready) == 1
			}, 500*time.Millisecond, 30*time.Millisecond)
		}
	})

	b.Run("TimerEfficiency", func(b *testing.B) {
		b.ReportAllocs()
		config := EventuallyConfig{
			Timeout:       300*time.Millisecond,
			Interval:      20*time.Millisecond,
			BackoffFactor: 1.5,
		}
		
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			var counter int32
			
			assert.EventuallyWith(func() bool {
				return atomic.AddInt32(&counter, 1) >= 5
			}, config)
		}
	})
}

// BenchmarkErrorReporting measures performance of error message generation.
func BenchmarkErrorReporting(b *testing.B) {
	b.Run("TimeoutError", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			
			assert.Eventually(func() bool {
				return false // Always fails
			}, 50*time.Millisecond, 10*time.Millisecond)
			
			// Force error message generation
			_ = assert.Error()
		}
	})

	b.Run("NeverError", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			
			assert.Never(func() bool {
				return true // Always fails immediately
			}, 100*time.Millisecond, 20*time.Millisecond)
			
			// Force error message generation
			_ = assert.Error()
		}
	})
}

// BenchmarkConfigurationOverhead measures overhead of different configurations.
func BenchmarkConfigurationOverhead(b *testing.B) {
	b.Run("MinimalConfig", func(b *testing.B) {
		config := EventuallyConfig{
			Timeout:  200*time.Millisecond,
			Interval: 50*time.Millisecond,
		}
		
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			assert.EventuallyWith(func() bool {
				return true
			}, config)
		}
	})

	b.Run("FullConfig", func(b *testing.B) {
		config := EventuallyConfig{
			Timeout:       500*time.Millisecond,
			Interval:      25*time.Millisecond,
			BackoffFactor: 1.8,
			MaxInterval:   200*time.Millisecond,
		}
		
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			assert.EventuallyWith(func() bool {
				return true
			}, config)
		}
	})

	b.Run("DefaultValidation", func(b *testing.B) {
		config := EventuallyConfig{
			Timeout:       -1*time.Second, // Invalid
			Interval:      -100*time.Millisecond, // Invalid
			BackoffFactor: 0.5, // Invalid
		}
		
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			assert := New(&mockT{})
			assert.EventuallyWith(func() bool {
				return true
			}, config)
		}
	})
}