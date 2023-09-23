package util

import (
	"reflect"
	"testing"
)

func TestChunkSlice(t *testing.T) {
	tests := []struct {
		name       string
		slice      []int
		chunkSize  int
		wantChunks [][]int
		wantError  bool
	}{
		{
			name:      "Empty slice",
			slice:     []int{},
			chunkSize: 3,
			wantError: true,
		},
		{
			name:       "Slice with length less than chunk size",
			slice:      []int{1, 2, 3},
			chunkSize:  5,
			wantChunks: [][]int{{1, 2, 3}},
		},
		{
			name:       "Slice with length equal to chunk size",
			slice:      []int{1, 2, 3},
			chunkSize:  3,
			wantChunks: [][]int{{1, 2, 3}},
		},
		{
			name:       "Slice with length greater than chunk size",
			slice:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			chunkSize:  3,
			wantChunks: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChunks, err := ChunkSlice(tt.slice, tt.chunkSize)
			if tt.wantError && err == nil {
				t.Errorf("Error expected, has not returns")
				t.FailNow()
			}
			if !reflect.DeepEqual(gotChunks, tt.wantChunks) {
				t.Errorf("ChunkSlice() = %v, want %v", gotChunks, tt.wantChunks)
			}
		})
	}
}
