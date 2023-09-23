package util

import "testing"

func TestMax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 > 0",
			args: args{
				a: 1,
				b: 0,
			},
			want: 1,
		},
		{
			name: "1 > -1255",
			args: args{
				a: 1,
				b: -1255,
			},
			want: 1,
		},
		{
			name: "0 > -1",
			args: args{
				a: 0,
				b: -1,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 > 0",
			args: args{
				a: 1,
				b: 0,
			},
			want: 0,
		},
		{
			name: "1 > -1255",
			args: args{
				a: 1,
				b: -1255,
			},
			want: -1255,
		},
		{
			name: "0 > -1",
			args: args{
				a: 0,
				b: -1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}
