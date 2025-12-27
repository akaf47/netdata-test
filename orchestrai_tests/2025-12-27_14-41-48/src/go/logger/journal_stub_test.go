package logger

import (
	"testing"
)

// TestJournalStubInit tests the initialization of journal stub
func TestJournalStubInit(t *testing.T) {
	tests := []struct {
		name string
		want JournalStub
	}{
		{
			name: "should initialize empty journal stub",
			want: JournalStub{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewJournalStub()
			if got == nil {
				t.Errorf("NewJournalStub() = nil, want non-nil")
			}
		})
	}
}

// TestJournalStubWrite tests writing to journal stub
func TestJournalStubWrite(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "should write normal string",
			input:   "test message",
			wantErr: false,
		},
		{
			name:    "should write empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "should write long string",
			input:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " + string(make([]byte, 10000)),
			wantErr: false,
		},
		{
			name:    "should write string with special characters",
			input:   "test\n\t\r special chars !@#$%^&*()",
			wantErr: false,
		},
		{
			name:    "should write unicode string",
			input:   "你好世界 🌍 тест",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := NewJournalStub()
			err := js.Write(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestJournalStubRead tests reading from journal stub
func TestJournalStubRead(t *testing.T) {
	tests := []struct {
		name    string
		writes  []string
		want    []string
		wantErr bool
	}{
		{
			name:    "should read empty when no writes",
			writes:  []string{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "should read single write",
			writes:  []string{"test"},
			want:    []string{"test"},
			wantErr: false,
		},
		{
			name:    "should read multiple writes in order",
			writes:  []string{"first", "second", "third"},
			want:    []string{"first", "second", "third"},
			wantErr: false,
		},
		{
			name:    "should read mixed content",
			writes:  []string{"", "test", "", "another"},
			want:    []string{"", "test", "", "another"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := NewJournalStub()
			for _, w := range tt.writes {
				if err := js.Write(w); err != nil && !tt.wantErr {
					t.Fatalf("Write() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			got := js.Read()
			if len(got) != len(tt.want) {
				t.Errorf("Read() length = %d, want %d", len(got), len(tt.want))
			}
			for i, g := range got {
				if i < len(tt.want) && g != tt.want[i] {
					t.Errorf("Read()[%d] = %s, want %s", i, g, tt.want[i])
				}
			}
		})
	}
}

// TestJournalStubClear tests clearing journal stub
func TestJournalStubClear(t *testing.T) {
	tests := []struct {
		name   string
		writes []string
	}{
		{
			name:   "should clear after single write",
			writes: []string{"test"},
		},
		{
			name:   "should clear after multiple writes",
			writes: []string{"first", "second", "third"},
		},
		{
			name:   "should clear empty journal",
			writes: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := NewJournalStub()
			for _, w := range tt.writes {
				js.Write(w)
			}
			js.Clear()
			got := js.Read()
			if len(got) != 0 {
				t.Errorf("Clear() did not clear all entries, got %d entries", len(got))
			}
		})
	}
}

// TestJournalStubLength tests getting length of journal stub
func TestJournalStubLength(t *testing.T) {
	tests := []struct {
		name   string
		writes []string
		want   int
	}{
		{
			name:   "should return 0 for empty journal",
			writes: []string{},
			want:   0,
		},
		{
			name:   "should return 1 for single write",
			writes: []string{"test"},
			want:   1,
		},
		{
			name:   "should return 5 for five writes",
			writes: []string{"a", "b", "c", "d", "e"},
			want:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := NewJournalStub()
			for _, w := range tt.writes {
				js.Write(w)
			}
			got := js.Length()
			if got != tt.want {
				t.Errorf("Length() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestJournalStubSequential tests sequential operations
func TestJournalStubSequential(t *testing.T) {
	js := NewJournalStub()

	// Write some messages
	messages := []string{"msg1", "msg2", "msg3"}
	for _, msg := range messages {
		if err := js.Write(msg); err != nil {
			t.Fatalf("Write() failed: %v", err)
		}
	}

	// Verify length
	if got := js.Length(); got != 3 {
		t.Errorf("Length() = %d, want 3", got)
	}

	// Read and verify
	if got := js.Read(); len(got) != 3 {
		t.Errorf("Read() length = %d, want 3", len(got))
	}

	// Clear
	js.Clear()
	if got := js.Length(); got != 0 {
		t.Errorf("Length() after Clear() = %d, want 0", got)
	}
}

// TestJournalStubConcurrency tests concurrent access patterns
func TestJournalStubConcurrency(t *testing.T) {
	js := NewJournalStub()

	// Write concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				js.Write("concurrent")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify writes completed
	length := js.Length()
	if length < 100 {
		t.Errorf("Concurrent writes resulted in %d entries, want at least 100", length)
	}
}

// TestJournalStubNilReceiver tests nil receiver handling
func TestJournalStubNilReceiver(t *testing.T) {
	var js *JournalStub
	
	// Test that methods handle nil gracefully or panic as expected
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Nil receiver caused panic as expected: %v", r)
		}
	}()

	if js != nil {
		js.Write("test")
	}
}

// TestJournalStubEdgeCases tests edge cases
func TestJournalStubEdgeCases(t *testing.T) {
	js := NewJournalStub()

	// Test very long string
	longStr := string(make([]byte, 1000000))
	if err := js.Write(longStr); err != nil {
		t.Logf("Write of long string: %v", err)
	}

	// Test multiple clears
	js.Clear()
	js.Clear()
	js.Clear()
	if js.Length() != 0 {
		t.Errorf("Multiple clears failed")
	}

	// Test write after clear
	js.Write("after clear")
	if js.Length() != 1 {
		t.Errorf("Write after clear failed")
	}
}