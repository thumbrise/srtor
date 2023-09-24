package fsutil

import (
	"os"
	"testing"
)

func TestSwapFiles(t *testing.T) {
	testCases := []struct {
		name          string
		inputA        string
		inputB        string
		expectedError bool
	}{
		{
			name:          "Swapping two existing files",
			inputA:        "file1.txt",
			inputB:        "file2.txt",
			expectedError: false,
		},
		{
			name:          "Swapping a non-existent file with an existing file",
			inputA:        "nonexistent.txt",
			inputB:        "file3.txt",
			expectedError: true,
		},
		{
			name:          "Swapping two non-existent files",
			inputA:        "nonexistent1.txt",
			inputB:        "nonexistent2.txt",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary files for testing
			file1, err := os.Create(tc.inputA)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			file1.Close()

			file2, err := os.Create(tc.inputB)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			file2.Close()

			// Call the SwapFiles function
			err = SwapFiles(tc.inputA, tc.inputB)

			// Check if the error matches the expected result
			if err != nil {
				if !tc.expectedError {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Clean up the temporary files
			os.Remove(tc.inputA)
			os.Remove(tc.inputB)
		})
	}
}
