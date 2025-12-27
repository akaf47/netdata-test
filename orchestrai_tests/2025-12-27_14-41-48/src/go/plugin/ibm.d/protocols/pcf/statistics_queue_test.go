package pcf

import (
	"testing"
)

// TestStatisticsQueueInitialization tests the initialization of StatisticsQueue
func TestStatisticsQueueInitialization(t *testing.T) {
	tests := []struct {
		name string
		fn   func() *StatisticsQueue
	}{
		{
			name: "should initialize empty statistics queue",
			fn: func() *StatisticsQueue {
				sq := NewStatisticsQueue()
				if sq == nil {
					t.Fatal("StatisticsQueue should not be nil")
				}
				return sq
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := tt.fn()
			if sq == nil {
				t.Error("expected non-nil StatisticsQueue")
			}
		})
	}
}

// TestStatisticsQueueEnqueue tests adding statistics to the queue
func TestStatisticsQueueEnqueue(t *testing.T) {
	tests := []struct {
		name      string
		enqueueFn func(*StatisticsQueue)
		validate  func(*StatisticsQueue) error
	}{
		{
			name: "should enqueue single statistic",
			enqueueFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("queue1"))
			},
			validate: func(sq *StatisticsQueue) error {
				if sq.Length() != 1 {
					t.Errorf("expected length 1, got %d", sq.Length())
				}
				return nil
			},
		},
		{
			name: "should enqueue multiple statistics",
			enqueueFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("queue1"))
				sq.Enqueue(NewQueueStatistic("queue2"))
				sq.Enqueue(NewQueueStatistic("queue3"))
			},
			validate: func(sq *StatisticsQueue) error {
				if sq.Length() != 3 {
					t.Errorf("expected length 3, got %d", sq.Length())
				}
				return nil
			},
		},
		{
			name: "should enqueue with nil statistic",
			enqueueFn: func(sq *StatisticsQueue) {
				sq.Enqueue(nil)
			},
			validate: func(sq *StatisticsQueue) error {
				if sq.Length() != 1 {
					t.Errorf("expected length 1 for nil enqueue, got %d", sq.Length())
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := NewStatisticsQueue()
			tt.enqueueFn(sq)
			tt.validate(sq)
		})
	}
}

// TestStatisticsQueueDequeue tests removing statistics from the queue
func TestStatisticsQueueDequeue(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func(*StatisticsQueue)
		dequeueFn func(*StatisticsQueue) *QueueStatistic
		validate  func(*StatisticsQueue, *QueueStatistic) bool
	}{
		{
			name: "should dequeue from non-empty queue",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("queue1"))
				sq.Enqueue(NewQueueStatistic("queue2"))
			},
			dequeueFn: func(sq *StatisticsQueue) *QueueStatistic {
				return sq.Dequeue()
			},
			validate: func(sq *StatisticsQueue, qs *QueueStatistic) bool {
				if qs == nil {
					t.Error("expected non-nil QueueStatistic")
					return false
				}
				if sq.Length() != 1 {
					t.Errorf("expected length 1 after dequeue, got %d", sq.Length())
					return false
				}
				return true
			},
		},
		{
			name: "should return nil when dequeuing from empty queue",
			setupFn: func(sq *StatisticsQueue) {
				// do nothing - queue is empty
			},
			dequeueFn: func(sq *StatisticsQueue) *QueueStatistic {
				return sq.Dequeue()
			},
			validate: func(sq *StatisticsQueue, qs *QueueStatistic) bool {
				if qs != nil {
					t.Error("expected nil from dequeuing empty queue")
					return false
				}
				return true
			},
		},
		{
			name: "should maintain FIFO order",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("first"))
				sq.Enqueue(NewQueueStatistic("second"))
				sq.Enqueue(NewQueueStatistic("third"))
			},
			dequeueFn: func(sq *StatisticsQueue) *QueueStatistic {
				first := sq.Dequeue()
				second := sq.Dequeue()
				if first == nil || second == nil {
					t.Error("unexpected nil dequeue")
					return nil
				}
				if first.QueueName != "first" || second.QueueName != "second" {
					t.Error("FIFO order not maintained")
				}
				return sq.Dequeue() // return third
			},
			validate: func(sq *StatisticsQueue, qs *QueueStatistic) bool {
				if qs == nil {
					t.Error("expected non-nil third element")
					return false
				}
				if qs.QueueName != "third" {
					t.Errorf("expected 'third', got %s", qs.QueueName)
					return false
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := NewStatisticsQueue()
			tt.setupFn(sq)
			qs := tt.dequeueFn(sq)
			tt.validate(sq, qs)
		})
	}
}

// TestStatisticsQueueLength tests the length operation
func TestStatisticsQueueLength(t *testing.T) {
	tests := []struct {
		name     string
		setupFn  func(*StatisticsQueue)
		expected int
	}{
		{
			name:     "should return 0 for empty queue",
			setupFn:  func(sq *StatisticsQueue) {},
			expected: 0,
		},
		{
			name: "should return correct length after enqueue",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("q1"))
				sq.Enqueue(NewQueueStatistic("q2"))
				sq.Enqueue(NewQueueStatistic("q3"))
			},
			expected: 3,
		},
		{
			name: "should return correct length after enqueue and dequeue",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("q1"))
				sq.Enqueue(NewQueueStatistic("q2"))
				sq.Enqueue(NewQueueStatistic("q3"))
				sq.Dequeue()
			},
			expected: 2,
		},
		{
			name: "should return 0 after dequeueing all elements",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("q1"))
				sq.Enqueue(NewQueueStatistic("q2"))
				sq.Dequeue()
				sq.Dequeue()
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := NewStatisticsQueue()
			tt.setupFn(sq)
			if length := sq.Length(); length != tt.expected {
				t.Errorf("expected length %d, got %d", tt.expected, length)
			}
		})
	}
}

// TestStatisticsQueuePeek tests peeking at the front of the queue
func TestStatisticsQueuePeek(t *testing.T) {
	tests := []struct {
		name     string
		setupFn  func(*StatisticsQueue)
		validate func(*StatisticsQueue, *QueueStatistic) bool
	}{
		{
			name:    "should peek at front element without removing",
			setupFn: func(sq *StatisticsQueue) { sq.Enqueue(NewQueueStatistic("front")) },
			validate: func(sq *StatisticsQueue, qs *QueueStatistic) bool {
				if qs == nil {
					t.Error("expected non-nil peek result")
					return false
				}
				if sq.Length() != 1 {
					t.Error("peek should not remove element")
					return false
				}
				return true
			},
		},
		{
			name:    "should return nil when peeking empty queue",
			setupFn: func(sq *StatisticsQueue) {},
			validate: func(sq *StatisticsQueue, qs *QueueStatistic) bool {
				if qs != nil {
					t.Error("expected nil for empty queue peek")
					return false
				}
				return true
			},
		},
		{
			name: "should return same element as first dequeue would return",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("first"))
				sq.Enqueue(NewQueueStatistic("second"))
			},
			validate: func(sq *StatisticsQueue, peek *QueueStatistic) bool {
				dequeue := sq.Dequeue()
				if peek == nil || dequeue == nil {
					t.Error("unexpected nil values")
					return false
				}
				if peek.QueueName != dequeue.QueueName {
					t.Error("peek and dequeue should return same front element")
					return false
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := NewStatisticsQueue()
			tt.setupFn(sq)
			qs := sq.Peek()
			tt.validate(sq, qs)
		})
	}
}

// TestStatisticsQueueClear tests clearing the queue
func TestStatisticsQueueClear(t *testing.T) {
	tests := []struct {
		name    string
		setupFn func(*StatisticsQueue)
	}{
		{
			name: "should clear non-empty queue",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("q1"))
				sq.Enqueue(NewQueueStatistic("q2"))
				sq.Enqueue(NewQueueStatistic("q3"))
			},
		},
		{
			name: "should clear already empty queue",
			setupFn: func(sq *StatisticsQueue) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := NewStatisticsQueue()
			tt.setupFn(sq)
			sq.Clear()
			if sq.Length() != 0 {
				t.Errorf("expected length 0 after clear, got %d", sq.Length())
			}
		})
	}
}

// TestStatisticsQueueIsEmpty tests the IsEmpty check
func TestStatisticsQueueIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		setupFn  func(*StatisticsQueue)
		expected bool
	}{
		{
			name:     "should return true for empty queue",
			setupFn:  func(sq *StatisticsQueue) {},
			expected: true,
		},
		{
			name: "should return false for non-empty queue",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("q1"))
			},
			expected: false,
		},
		{
			name: "should return true after clearing",
			setupFn: func(sq *StatisticsQueue) {
				sq.Enqueue(NewQueueStatistic("q1"))
				sq.Clear()
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := NewStatisticsQueue()
			tt.setupFn(sq)
			if sq.IsEmpty() != tt.expected {
				t.Errorf("expected IsEmpty %v, got %v", tt.expected, sq.IsEmpty())
			}
		})
	}
}

// TestQueueStatisticCreation tests QueueStatistic creation
func TestQueueStatisticCreation(t *testing.T) {
	tests := []struct {
		name      string
		queueName string
	}{
		{
			name:      "should create statistic with simple queue name",
			queueName: "QUEUE1",
		},
		{
			name:      "should create statistic with complex queue name",
			queueName: "QUEUE.WITH.DOTS",
		},
		{
			name:      "should create statistic with empty queue name",
			queueName: "",
		},
		{
			name:      "should create statistic with special characters",
			queueName: "QUEUE-WITH_SPECIAL.CHARS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := NewQueueStatistic(tt.queueName)
			if qs == nil {
				t.Fatal("expected non-nil QueueStatistic")
			}
			if qs.QueueName != tt.queueName {
				t.Errorf("expected queue name %s, got %s", tt.queueName, qs.QueueName)
			}
		})
	}
}

// TestQueueStatisticFields tests QueueStatistic field operations
func TestQueueStatisticFields(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func(*QueueStatistic)
		fieldName string
		validate  func(*QueueStatistic) bool
	}{
		{
			name: "should store depth metric",
			setupFn: func(qs *QueueStatistic) {
				qs.SetMetric("depth", int64(100))
			},
			fieldName: "depth",
			validate: func(qs *QueueStatistic) bool {
				depth := qs.GetMetric("depth")
				if depth != int64(100) {
					t.Errorf("expected depth 100, got %v", depth)
					return false
				}
				return true
			},
		},
		{
			name: "should handle negative depth values",
			setupFn: func(qs *QueueStatistic) {
				qs.SetMetric("depth", int64(-1))
			},
			fieldName: "depth",
			validate: func(qs *QueueStatistic) bool {
				depth := qs.GetMetric("depth")
				if depth != int64(-1) {
					t.Errorf("expected depth -1, got %v", depth)
					return false
				}
				return true
			},
		},
		{
			name: "should handle zero values",
			setupFn: func(qs *QueueStatistic) {
				qs.SetMetric("depth", int64(0))
			},
			fieldName: "depth",
			validate: func(qs *QueueStatistic) bool {
				depth := qs.GetMetric("depth")
				if depth != int64(0) {
					t.Errorf("expected depth 0, got %v", depth)
					return false
				}
				return true
			},
		},
		{
			name: "should handle large values",
			setupFn: func(qs *QueueStatistic) {
				qs.SetMetric("depth", int64(9223372036854775807)) // MaxInt64
			},
			fieldName: "depth",
			validate: func(qs *QueueStatistic) bool {
				depth := qs.GetMetric("depth")
				if depth != int64(9223372036854775807) {
					t.Errorf("expected MaxInt64, got %v", depth)
					return false
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := NewQueueStatistic("test_queue")
			tt.setupFn(qs)
			tt.validate(qs)
		})
	}
}

// TestStatisticsQueueConcurrentOperations tests concurrent access patterns
func TestStatisticsQueueConcurrentOperations(t *testing.T) {
	t.Run("should handle alternating enqueue and dequeue", func(t *testing.T) {
		sq := NewStatisticsQueue()
		sq.Enqueue(NewQueueStatistic("q1"))
		sq.Enqueue(NewQueueStatistic("q2"))
		sq.Dequeue()
		sq.Enqueue(NewQueueStatistic("q3"))
		
		if sq.Length() != 2 {
			t.Errorf("expected length 2, got %d", sq.Length())
		}
		
		qs := sq.Dequeue()
		if qs == nil || qs.QueueName != "q2" {
			t.Error("expected 'q2' from FIFO order")
		}
	})

	t.Run("should handle multiple clears", func(t *testing.T) {
		sq := NewStatisticsQueue()
		sq.Enqueue(NewQueueStatistic("q1"))
		sq.Clear()
		sq.Enqueue(NewQueueStatistic("q2"))
		sq.Clear()
		sq.Clear()
		
		if sq.Length() != 0 {
			t.Errorf("expected length 0 after multiple clears, got %d", sq.Length())
		}
	})
}

// TestStatisticsQueueBoundaryConditions tests boundary conditions
func TestStatisticsQueueBoundaryConditions(t *testing.T) {
	t.Run("should handle large number of enqueues", func(t *testing.T) {
		sq := NewStatisticsQueue()
		const count = 10000
		
		for i := 0; i < count; i++ {
			sq.Enqueue(NewQueueStatistic("queue"))
		}
		
		if sq.Length() != count {
			t.Errorf("expected length %d, got %d", count, sq.Length())
		}
		
		for i := 0; i < count; i++ {
			if sq.Dequeue() == nil {
				t.Errorf("unexpected nil dequeue at iteration %d", i)
				break
			}
		}
		
		if sq.Length() != 0 {
			t.Errorf("expected empty queue, got length %d", sq.Length())
		}
	})

	t.Run("should handle peek on single element queue", func(t *testing.T) {
		sq := NewStatisticsQueue()
		sq.Enqueue(NewQueueStatistic("only"))
		
		peek := sq.Peek()
		if peek == nil {
			t.Fatal("peek returned nil")
		}
		
		dequeue := sq.Dequeue()
		if dequeue == nil {
			t.Fatal("dequeue returned nil")
		}
		
		if peek.QueueName != dequeue.QueueName {
			t.Error("peek and dequeue mismatch on single element")
		}
	})
}