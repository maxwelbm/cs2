package memory

import (
	"syscall"
	"testing"
)

func Test_uint16ToString(t *testing.T) {
	tests := []struct {
		name string
		args [syscall.MAX_PATH]uint16
		want string
	}{
		{
			name: "case 01",
			args: [syscall.MAX_PATH]uint16{1, 2, 3, 4, 5, 8, 99},
			want: "c",
		},
		{
			name: "case 02",
			args: [syscall.MAX_PATH]uint16{3121, 123, 213, 43, 88},
			want: "ఱ{Õ+X",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxPathToString(tt.args); got != tt.want {
				t.Errorf("uint16ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
