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
			name:       "Empty slice",
			slice:      []int{},
			chunkSize:  3,
			wantError:  true,
			wantChunks: [][]int{},
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
			gotChunks, err := chunkSliceBySize(tt.slice, tt.chunkSize)

			if tt.wantError {
				if err == nil {
					t.Errorf("Error expected, has not returns")
					t.FailNow()
				}

				return
			}

			if !reflect.DeepEqual(gotChunks, tt.wantChunks) {
				t.Errorf("chunkSliceBySize() = %v, want %v", gotChunks, tt.wantChunks)
			}
		})
	}
}
func TestCalculateChunkSize(t *testing.T) {
	tests := []struct {
		name     string
		divider  int
		total    int
		expected int
	}{
		{
			name:     "Positive divider and total",
			divider:  3,
			total:    9,
			expected: 3,
		},
		{
			name:     "Negative divider",
			divider:  -2,
			total:    10,
			expected: 5,
		},
		{
			name:     "Total < divider",
			divider:  5,
			total:    3,
			expected: 1,
		},
		{
			name:     "Zero divider",
			divider:  0,
			total:    8,
			expected: 8,
		},
		{
			name:     "All zero",
			divider:  0,
			total:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateChunkSize(tt.total, tt.divider)
			if got != tt.expected {
				t.Errorf("calculateChunkSize(%d, %d) = %d; want %d", tt.divider, tt.total, got, tt.expected)
			}
		})
	}
}
func TestSplitSlice(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		chunks   int
		expected [][]int
		errMsg   string
	}{
		{
			name:     "Splitting a slice",
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			chunks:   3,
			expected: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			errMsg:   "",
		},
		{
			name:     "Splitting a slice with odd number",
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			chunks:   3,
			expected: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}},
			errMsg:   "",
		},
		{
			name:     "Splitting an empty slice",
			slice:    []int{},
			chunks:   5,
			expected: [][]int{},
			errMsg:   "empty v passed",
		},
		{
			name:     "Splitting with zero chunks",
			slice:    []int{1, 2, 3, 4, 5},
			chunks:   0,
			expected: [][]int{{1, 2, 3, 4, 5}},
			errMsg:   "empty size passed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SplitSlice(tc.slice, tc.chunks)
			if err != nil && err.Error() != tc.errMsg {
				t.Errorf("Unexpected error. Expected: %v, Got: %v", tc.errMsg, err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Unexpected result. Expected: %v, Got: %v", tc.expected, result)
			}
		})
	}
}
