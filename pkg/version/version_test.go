package version

import (
	"os"
	"path"
	"testing"
)

func TestGetVersion(t *testing.T) {
	testVersion := "1.0.0"
	wd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	if err := os.WriteFile(path.Join(wd, "/VERSION"), []byte(testVersion), 0644); err != nil {
		t.Fatalf("Failed to create test VERSION file: %v", err)
	}

	tests := []struct {
		name    string
		want    *Version
		wantErr bool
	}{
		{
			name: "successful version read",
			want: &Version{
				Version: testVersion,
				Build:   "dev",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVersion()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.Build != tt.want.Build {
				t.Errorf("GetVersion() Build = %v, want %v", got.Build, tt.want.Build)
			}
		})
	}

	// Test error case when VERSION file doesn't exist
	t.Run("file not found", func(t *testing.T) {
		// Rename the VERSION file temporarily
		err := os.Remove(path.Join(wd, "/VERSION"))

		if err != nil {
			t.Fatalf("Failed to remove test VERSION file: %v", err)
		}

		got, err := GetVersion()

		if err == nil {
			t.Error("GetVersion() expected error when VERSION file doesn't exist")
		}

		if got != nil {
			t.Error("GetVersion() expected nil result when VERSION file doesn't exist")
		}
	})
}
