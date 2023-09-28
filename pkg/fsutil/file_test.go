package fsutil

import (
	"os"
	"testing"
)

func TestFileSwap(t *testing.T) {
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

			// Call the FileSwap function
			err = FileSwap(tc.inputA, tc.inputB)

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

func TestFileOpenOrCreate(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Create new zip file",
			path:    "test.zip",
			wantErr: false,
		},
		{
			name:    "Ignore existing zip file",
			path:    "test.zip",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				file, err := os.Create(tt.path)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer os.Remove(tt.path)
				defer file.Close()
			}

			f, err := FileOpenOrCreate(tt.path)
			f.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileOpenOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	exists := FileExists(file.Name())
	if !exists {
		t.Errorf("Expected file %s to exist, but it doesn't", file.Name())
	}

	exists = FileExists("nonexistentfile")
	if exists {
		t.Errorf("Expected file nonexistentfile to not exist, but it does")
	}

	exists = FileExists(".")
	if exists {
		t.Errorf("Expected directory . to not be a file, but it is")
	}
}
