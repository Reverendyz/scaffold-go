package create

import (
	"testing"
)

func TestReadYamlFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		environment string
		isProd      bool
		wantErr     bool
	}{
		{
			name:        "Test with a valid name",
			filename:    "examples/somefile.yaml",
			environment: "test",
			isProd:      false,
			wantErr:     false,
		},
		{
			name:        "Test with an invalid name",
			filename:    "examples/somemockfile.yaml",
			environment: "test",
			isProd:      false,
			wantErr:     true,
		},
		{
			name:        "Invalid Extension",
			filename:    "examples/somefile.yasaml",
			environment: "test",
			isProd:      false,
			wantErr:     true,
		},
		{
			name:        "Invalid with existent file but wrong extension",
			filename:    "examples/somefile.json",
			environment: "test",
			isProd:      false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadYamlFile(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("ReadJsonFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
