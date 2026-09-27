package collection

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const duration = time.Millisecond * 50

func TestNewRollingWindow(t *testing.T) {
	assert.NotNil(t, NewRollingWindow[int64, *Bucket[int64]](func() *Bucket[int64] {
		return new(Bucket[int64])
	}, 10, time.Second))
	assert.Panics(t, func() {
		NewRollingWindow[int64, *Bucket[int64]](func() *Bucket[int64] {
			return new(Bucket[int64])
		}, 0, time.Second)
	})
}

func TestRollingWindowAdd(t *testing.T) {
	const size = 3
	r, clock := newTestRollingWindow(size, duration)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket[float64]) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	clock.advance(duration)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	clock.advance(duration)
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
	clock.advance(duration)
	r.Add(7)
	assert.Equal(t, []float64{5, 15, 7}, listBuckets())
}

func TestRollingWindowReset(t *testing.T) {
	const size = 3
	r, clock := newTestRollingWindow(size, duration, IgnoreCurrentBucket[float64, *Bucket[float64]]())
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket[float64]) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	r.Add(1)
	clock.advance(duration)
	assert.Equal(t, []float64{0, 1}, listBuckets())
	clock.advance(duration)
	assert.Equal(t, []float64{1}, listBuckets())
	clock.advance(duration)
	assert.Nil(t, listBuckets())

	// cross window
	r.Add(1)
	clock.advance(duration * 10)
	assert.Nil(t, listBuckets())
}

func TestRollingWindowReduce(t *testing.T) {
	const size = 4
	tests := []struct {
		name          string
		ignoreCurrent bool
		expect        float64
	}{
		{
			name:   "all buckets",
			expect: 10,
		},
		{
			name:          "ignore current bucket",
			ignoreCurrent: true,
			expect:        4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var opts []RollingWindowOption[float64, *Bucket[float64]]
			if test.ignoreCurrent {
				opts = append(opts, IgnoreCurrentBucket[float64, *Bucket[float64]]())
			}
			r, clock := newTestRollingWindow(size, duration, opts...)
			for x := 0; x < size; x++ {
				for i := 0; i <= x; i++ {
					r.Add(float64(i))
				}
				if x < size-1 {
					clock.advance(duration)
				}
			}
			var result float64
			r.Reduce(func(b *Bucket[float64]) {
				result += b.Sum
			})
			assert.Equal(t, test.expect, result)
		})
	}
}

func TestRollingWindowBucketTimeBoundary(t *testing.T) {
	const size = 3
	interval := time.Millisecond * 30
	r, clock := newTestRollingWindow(size, interval)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket[float64]) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	clock.advance(time.Millisecond * 45)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	// Advance less than one interval, but cross the aligned bucket boundary.
	clock.advance(time.Millisecond * 20)
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
	clock.advance(time.Millisecond * 100)
	r.Add(7)
	r.Add(8)
	r.Add(9)
	assert.Equal(t, []float64{0, 0, 24}, listBuckets())
}

func TestRollingWindowDataRace(t *testing.T) {
	const size = 3
	r := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, size, duration)
	stop := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				r.Add(float64(rand.Int63()))
				time.Sleep(duration / 2)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				r.Reduce(func(b *Bucket[float64]) {})
			}
		}
	}()
	time.Sleep(duration * 5)
	close(stop)
	wg.Wait()
}

func TestRollingWindowExactBoundary(t *testing.T) {
	const interval = 30 * time.Millisecond
	tests := []struct {
		name      string
		elapsed   time.Duration
		beforeAdd []float64
		afterAdd  []float64
	}{
		{"before boundary", interval - time.Nanosecond, []float64{0, 0, 1}, []float64{0, 0, 3}},
		{"at boundary", interval, []float64{0, 1}, []float64{0, 1, 2}},
		{"after boundary", interval + time.Nanosecond, []float64{0, 1}, []float64{0, 1, 2}},
		{"two buckets", 2 * interval, []float64{1}, []float64{1, 0, 2}},
		{"full window", 3 * interval, nil, []float64{0, 0, 2}},
		{"multiple windows", 6*interval + time.Nanosecond, nil, []float64{0, 0, 2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, clock := newTestRollingWindow(3, interval)
			listBuckets := func() []float64 {
				var values []float64
				r.Reduce(func(b *Bucket[float64]) { values = append(values, b.Sum) })
				return values
			}
			r.Add(1)
			clock.advance(test.elapsed)
			assert.Equal(t, test.beforeAdd, listBuckets())
			r.Add(2)
			assert.Equal(t, test.afterAdd, listBuckets())
			assert.Equal(t, test.elapsed-test.elapsed%interval, r.lastTime)
		})
	}
}

func TestRollingWindowClocksAreIndependent(t *testing.T) {
	first, firstClock := newTestRollingWindow(3, duration)
	second, _ := newTestRollingWindow(3, duration)
	first.Add(1)
	second.Add(2)
	firstClock.advance(3 * duration)
	var firstSum, secondSum float64
	first.Reduce(func(b *Bucket[float64]) { firstSum += b.Sum })
	second.Reduce(func(b *Bucket[float64]) { secondSum += b.Sum })
	assert.Zero(t, firstSum)
	assert.Equal(t, float64(2), secondSum)
}

type rollingWindowTestClock struct {
	time atomic.Int64
}

func (c *rollingWindowTestClock) now() time.Duration {
	return time.Duration(c.time.Load())
}

func (c *rollingWindowTestClock) advance(d time.Duration) {
	c.time.Add(int64(d))
}

func newTestRollingWindow(size int, interval time.Duration,
	opts ...RollingWindowOption[float64, *Bucket[float64]]) (*RollingWindow[float64, *Bucket[float64]], *rollingWindowTestClock) {
	clock := new(rollingWindowTestClock)
	opts = append(opts, func(r *RollingWindow[float64, *Bucket[float64]]) {
		r.now = clock.now
		r.lastTime = clock.now()
	})
	r := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, size, interval, opts...)
	return r, clock
}
