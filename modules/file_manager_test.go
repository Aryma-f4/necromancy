package modules

import (
	"testing"

	"github.com/Aryma-f4/necromancy/core"
)

// MockSession creates a mock session for testing
type MockSession struct {
	*core.Session
	output string
}

func TestFileManagerModule(t *testing.T) {
	module := &FileManagerModule{}

	// Test module interface compliance
	if module.Name() != "filemanager" {
		t.Errorf("Expected module name 'filemanager', got '%s'", module.Name())
	}

	if module.Description() != "Interactive file manager with btop-like UI for target system" {
		t.Errorf("Unexpected module description: %s", module.Description())
	}
}

func TestFileInfo(t *testing.T) {
	file := FileInfo{
		Name:    "test.txt",
		Size:    1024,
		Mode:    "-rw-r--r--",
		ModTime: "2024-01-01 12:00",
		IsDir:   false,
		IsLink:  false,
	}

	if file.Name != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", file.Name)
	}

	if file.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", file.Size)
	}
}

func TestFileManagerCommands(t *testing.T) {
	// Create a mock session
	mockSession := &core.Session{
		ID:         1,
		RemoteAddr: "127.0.0.1:1234",
		Type:       "Raw",
	}

	fmc := NewFileManagerCommands(mockSession)

	// Test Windows commands
	windowsCmd := fmc.getWindowsFileListCommand("C:\\test")
	if !contains(windowsCmd, "powershell") {
		t.Error("Windows file list command should contain PowerShell")
	}

	// Test Linux commands
	linuxCmd := fmc.getLinuxFileListCommand("/tmp")
	if !contains(linuxCmd, "ls -la") {
		t.Error("Linux file list command should contain 'ls -la'")
	}

	// Test path sanitization
	sanitized := fmc.SanitizePath("/tmp/test; rm -rf /")
	if contains(sanitized, ";") {
		t.Error("Sanitized path should not contain dangerous characters")
	}

	// Test path validation
	if !fmc.IsValidPath("/tmp/test") {
		t.Error("Valid path should be recognized as valid")
	}

	if fmc.IsValidPath("/tmp/test; rm -rf /") {
		t.Error("Invalid path with dangerous characters should be rejected")
	}
}

func TestFileListParsing(t *testing.T) {
	mockSession := &core.Session{
		ID:         1,
		RemoteAddr: "127.0.0.1:1234",
		Type:       "Raw",
	}

	fmc := NewFileManagerCommands(mockSession)

	// Test Linux file list parsing
	linuxOutput := `README.md|2048|-rw-r--r--|2024-01-01 12:00|false|false
config|4096|drwxr-xr-x|2024-01-01 11:45|true|false
script.sh|1024|-rwxr-xr-x|2024-01-01 10:30|false|false`

	files := fmc.parseLinuxFileList(linuxOutput)

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}

	// Check first file
	if files[0].Name != "README.md" {
		t.Errorf("Expected first file name 'README.md', got '%s'", files[0].Name)
	}

	if files[0].Size != 2048 {
		t.Errorf("Expected first file size 2048, got %d", files[0].Size)
	}

	if files[0].IsDir {
		t.Error("README.md should not be a directory")
	}

	// Check directory
	if !files[1].IsDir {
		t.Error("config should be a directory")
	}

	// Check executable file
	if files[2].Name != "script.sh" {
		t.Errorf("Expected third file name 'script.sh', got '%s'", files[2].Name)
	}
}

func TestFileSizeFormatting(t *testing.T) {
	mockSession := &core.Session{
		ID:         1,
		RemoteAddr: "127.0.0.1:1234",
		Type:       "Raw",
	}

	fmc := NewFileManagerCommands(mockSession)

	// Test size formatting
	testCases := []struct {
		size     int64
		expected string
	}{
		{0, ""},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tc := range testCases {
		result := fmc.FormatSize(tc.size)
		if result != tc.expected {
			t.Errorf("Expected size %d to format as '%s', got '%s'", tc.size, tc.expected, result)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
