package main

import (
	"testing"
)

func Test_getDegree(t *testing.T) {
	type args struct {
		x uint64
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"x = 0", args{0}, -1},
		{"x = 1", args{1}, 0},
		{"x = 0xFFFFFFFFFFFFFFFF", args{0xFFFFFFFFFFFFFFFF}, 63},
		{"x = 0x7FFFFFFF", args{0x7FFFFFFF}, 30},
		{"x = 0x17FC1", args{0x17FC1}, 16},
		{"x = 0x27E6C463AA52", args{0x27E6C463AA52}, 45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDegree(tt.args.x); got != tt.want {
				t.Errorf("getDegree() = %v, want %v", got, tt.want)
			}
		})
	}
}
