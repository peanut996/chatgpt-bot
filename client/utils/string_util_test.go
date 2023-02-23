package utils

import (
	"reflect"
	"testing"
)

func TestSplitMessageByMaxSize(t *testing.T) {
	type args struct {
		msg     string
		maxSize int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		//should split by max size 5
		{args: args{msg: "12345678910", maxSize: 5}, want: []string{"12345", "67891", "0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitMessageByMaxSize(tt.args.msg, tt.args.maxSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitMessageByMaxSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
