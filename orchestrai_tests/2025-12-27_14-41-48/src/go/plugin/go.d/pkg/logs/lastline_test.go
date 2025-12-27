package logs

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestLastLineReader_ReadString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "single line without newline",
			input:   "hello",
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "single line with newline",
			input:   "hello\n",
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "multiple lines returns last line",
			input:   "line1\nline2\nline3\n",
			want:    "line3",
			wantErr: false,
		},
		{
			name:    "multiple lines without final newline",
			input:   "line1\nline2\nline3",
			want:    "line3",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "only newline",
			input:   "\n",
			want:    "",
			wantErr: false,
		},
		{
			name:    "multiple empty lines",
			input:   "\n\n\n",
			want:    "",
			wantErr: false,
		},
		{
			name:    "line with spaces",
			input:   "  hello world  \n",
			want:    "  hello world  ",
			wantErr: false,
		},
		{
			name:    "line with special characters",
			input:   "special!@#$%^&*()\n",
			want:    "special!@#$%^&*()",
			wantErr: false,
		},
		{
			name:    "very long line",
			input:   bytes.Repeat([]byte("x"), 10000) + "\n",
			want:    string(bytes.Repeat([]byte("x"), 10000)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLastLineReader(bytes.NewReader([]byte(tt.input)))
			got, err := r.ReadString('\n')

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ReadString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLastLineReader_ReadString_IOErrors(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{
			name:    "error from reader",
			reader:  &errorReader{},
			wantErr: true,
		},
		{
			name:    "EOF on first read",
			reader:  bytes.NewReader([]byte{}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLastLineReader(tt.reader)
			_, err := r.ReadString('\n')

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLastLineReader_WithCustomDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter byte
		want      string
		wantErr   bool
	}{
		{
			name:      "custom delimiter semicolon",
			input:     "line1;line2;line3;",
			delimiter: ';',
			want:      "line3",
			wantErr:   false,
		},
		{
			name:      "custom delimiter comma",
			input:     "line1,line2,line3",
			delimiter: ',',
			want:      "line3",
			wantErr:   false,
		},
		{
			name:      "custom delimiter zero byte",
			input:     "line1\x00line2\x00line3\x00",
			delimiter: '\x00',
			want:      "line3",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLastLineReader(bytes.NewReader([]byte(tt.input)))
			got, err := r.ReadString(tt.delimiter)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ReadString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLastLineReader_BufferBoundary(t *testing.T) {
	// Test reading at exact buffer boundaries
	input := bytes.Repeat([]byte("a"), bufferSize) + "\n" + bytes.Repeat([]byte("b"), bufferSize) + "\n"
	r := NewLastLineReader(bytes.NewReader(input))

	got, err := r.ReadString('\n')
	if err != nil {
		t.Fatalf("ReadString() error = %v", err)
	}

	want := string(bytes.Repeat([]byte("b"), bufferSize))
	if got != want {
		t.Errorf("ReadString() = %q, want %q", got, want)
	}
}

func TestLastLineReader_ConsecutiveReads(t *testing.T) {
	// Test that multiple reads on the same reader work correctly
	input := "line1\nline2\nline3\n"
	r := NewLastLineReader(bytes.NewReader([]byte(input)))

	// First read should return last line
	got1, err1 := r.ReadString('\n')
	if err1 != nil {
		t.Fatalf("First ReadString() error = %v", err1)
	}

	// Second read on same reader should return EOF
	got2, err2 := r.ReadString('\n')
	if err2 == nil {
		t.Errorf("Second ReadString() should return error, got: %q", got2)
	}
}

// errorReader is a helper that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// Verify NewLastLineReader creates proper instance
func TestNewLastLineReader(t *testing.T) {
	input := "test"
	reader := bytes.NewReader([]byte(input))
	
	r := NewLastLineReader(reader)
	if r == nil {
		t.Error("NewLastLineReader() returned nil")
	}
}

func TestLastLineReader_LargeInput(t *testing.T) {
	// Test with input larger than buffer size
	lines := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		lines[i] = []byte("line " + string(rune(i)))
	}
	input := bytes.Join(lines, []byte("\n"))
	input = append(input, '\n')

	r := NewLastLineReader(bytes.NewReader(input))
	got, err := r.ReadString('\n')
	if err != nil {
		t.Fatalf("ReadString() error = %v", err)
	}

	if got != "line 99" {
		t.Errorf("ReadString() = %q, want 'line 99'", got)
	}
}

func TestLastLineReader_WindowsLineEndings(t *testing.T) {
	input := "line1\r\nline2\r\nline3\r\n"
	r := NewLastLineReader(bytes.NewReader([]byte(input)))

	got, err := r.ReadString('\n')
	if err != nil {
		t.Fatalf("ReadString() error = %v", err)
	}

	// Should return everything up to the delimiter
	expected := "line3\r"
	if got != expected {
		t.Errorf("ReadString() = %q, want %q", got, expected)
	}
}