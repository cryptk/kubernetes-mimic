package utils

import (
	"reflect"
	"testing"
)

func TestBoolAddr(t *testing.T) {
	type args struct {
		b bool
	}

	tests := []struct {
		name string
		args args
		want *bool
	}{
		{"Test positive bool", args{b: true}, func() *bool { b := true; return &b }()},
		{"Test negative bool", args{b: false}, func() *bool { b := false; return &b }()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolAddr(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
