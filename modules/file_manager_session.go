package modules

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/Aryma-f4/necromancy/core"
)

// FileManagerSession extends session with file manager capabilities
type FileManagerSession struct {
	*core.Session
	commandOutput chan string
	isWaiting     bool
}

// NewFileManagerSession creates a file manager session wrapper
func NewFileManagerSession(session *core.Session) *FileManagerSession {
	fms := &FileManagerSession{
		Session:       session,
		commandOutput: make(chan string, 100),
		isWaiting:     false,
	}

	// Only start monitoring if session has valid history buffer
	if session != nil && session.History != nil {
		go fms.monitorOutput()
	}

	return fms
}

// monitorOutput monitors session output for command responses
func (fms *FileManagerSession) monitorOutput() {
	if fms.Session == nil || fms.Session.History == nil {
		return
	}

	scanner := bufio.NewScanner(fms.Session.History)
	for scanner.Scan() {
		line := scanner.Text()
		if fms.isWaiting {
			select {
			case fms.commandOutput <- line:
			default:
				// Channel full, skip
			}
		}
	}
}

// ExecuteCommand executes a command and returns output
func (fms *FileManagerSession) ExecuteCommand(cmd string, timeout time.Duration) (string, error) {
	fms.isWaiting = true
	defer func() { fms.isWaiting = false }()

	// Clear output channel
	for len(fms.commandOutput) > 0 {
		<-fms.commandOutput
	}

	// Send command with proper formatting
	_, err := fms.Write([]byte(cmd + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send command: %v", err)
	}

	// Collect output with timeout
	var output strings.Builder
	deadline := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case line := <-fms.commandOutput:
			if line != "" {
				output.WriteString(line + "\n")
			}
		case <-ticker.C:
			// Check if we have enough output
			if output.Len() > 0 {
				// Give it a bit more time
				select {
				case line := <-fms.commandOutput:
					if line != "" {
						output.WriteString(line + "\n")
					}
				case <-time.After(200 * time.Millisecond):
					// No more output, return what we have
					return strings.TrimSpace(output.String()), nil
				}
			}
		case <-deadline:
			return strings.TrimSpace(output.String()), nil
		}
	}
}

// FileManagerCommands provides command generation for file operations
type FileManagerCommands struct {
	session *FileManagerSession
}

// NewFileManagerCommands creates a new file manager command generator
func NewFileManagerCommands(session *core.Session) *FileManagerCommands {
	return &FileManagerCommands{session: NewFileManagerSession(session)}
}

// GetFileListCommand generates command to list files based on OS
func (fmc *FileManagerCommands) GetFileListCommand(path string) string {
	os := fmc.session.DetectedOS()

	switch os {
	case "windows":
		return fmc.getWindowsFileListCommand(path)
	case "linux":
		return fmc.getLinuxFileListCommand(path)
	default:
		return fmc.getUnixFileListCommand(path)
	}
}

// getWindowsFileListCommand generates Windows-specific file listing
func (fmc *FileManagerCommands) getWindowsFileListCommand(path string) string {
	// Simple dir command for Windows
	return fmt.Sprintf(`dir /a "%s" 2>nul`, path)
}

// getLinuxFileListCommand generates Linux-specific file listing
func (fmc *FileManagerCommands) getLinuxFileListCommand(path string) string {
	// Simple ls command for Linux
	return fmt.Sprintf(`ls -la "%s" 2>/dev/null`, path)
}

// getUnixFileListCommand generates generic Unix file listing
func (fmc *FileManagerCommands) getUnixFileListCommand(path string) string {
	return fmc.getLinuxFileListCommand(path) // Default to Linux format
}

// ListFiles executes file listing command and parses output
func (fmc *FileManagerCommands) ListFiles(path string) ([]FileInfo, error) {
	cmd := fmc.GetFileListCommand(path)
	output, err := fmc.session.ExecuteCommand(cmd, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return fmc.ParseFileListOutput(output), nil
}

// ParseFileListOutput parses the output from file listing commands
func (fmc *FileManagerCommands) ParseFileListOutput(output string) []FileInfo {
	os := fmc.session.DetectedOS()

	switch os {
	case "windows":
		return fmc.parseWindowsFileList(output)
	case "linux":
		return fmc.parseLinuxFileList(output)
	default:
		return fmc.parseUnixFileList(output)
	}
}

// parseWindowsFileList parses Windows dir output
func (fmc *FileManagerCommands) parseWindowsFileList(output string) []FileInfo {
	var files []FileInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Volume") || strings.HasPrefix(line, "Directory of") || strings.Contains(line, "<DIR>") && strings.Contains(line, "bytes") {
			continue
		}

		// Parse Windows dir output format
		// Format: MM/DD/YYYY  HH:MM AM/PM    <DIR>          .
		//         MM/DD/YYYY  HH:MM AM/PM             1,234 filename.ext

		// Skip header lines and summary
		if strings.Contains(line, "File(s)") || strings.Contains(line, "Dir(s)") {
			continue
		}

		// Try to parse file/directory entry
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			date := parts[0]
			time := parts[1]
			ampm := ""
			if len(parts) > 3 && (parts[2] == "AM" || parts[2] == "PM") {
				ampm = parts[2]
			}

			sizeOrDir := parts[len(parts)-2]
			name := parts[len(parts)-1]

			file := FileInfo{
				Name:    name,
				ModTime: date + " " + time + ampm,
				IsDir:   sizeOrDir == "<DIR>",
				IsLink:  false,
			}

			if sizeOrDir != "<DIR>" {
				// Parse file size (remove commas)
				sizeStr := strings.ReplaceAll(sizeOrDir, ",", "")
				fmt.Sscanf(sizeStr, "%d", &file.Size)
			}

			// Simple mode detection
			if file.IsDir {
				file.Mode = "drwxr-xr-x"
			} else {
				file.Mode = "-rw-r--r--"
			}

			files = append(files, file)
		}
	}

	return files
}

// parseLinuxFileList parses Linux ls -la output
func (fmc *FileManagerCommands) parseLinuxFileList(output string) []FileInfo {
	var files []FileInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total ") {
			continue
		}

		// Parse ls -la format: -rw-r--r-- 1 user group 1234 Jan 01 12:00 filename
		parts := strings.Fields(line)
		if len(parts) >= 9 {
			mode := parts[0]
			size := int64(0)
			fmt.Sscanf(parts[4], "%d", &size)

			// Build name from remaining parts (handle spaces in filenames)
			name := strings.Join(parts[8:], " ")

			// Build timestamp from parts 5,6,7
			modTime := parts[5] + " " + parts[6] + " " + parts[7]

			file := FileInfo{
				Name:    name,
				Size:    size,
				Mode:    mode,
				ModTime: modTime,
				IsDir:   strings.HasPrefix(mode, "d"),
				IsLink:  strings.HasPrefix(mode, "l"),
			}

			files = append(files, file)
		}
	}

	return files
}

// parseUnixFileList parses generic Unix output
func (fmc *FileManagerCommands) parseUnixFileList(output string) []FileInfo {
	return fmc.parseLinuxFileList(output)
}

// GetCurrentDirectory gets the current working directory
func (fmc *FileManagerCommands) GetCurrentDirectory() (string, error) {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = `cd` // Simple cd command for Windows
	default:
		cmd = `pwd` // Simple pwd for Unix
	}

	output, err := fmc.session.ExecuteCommand(cmd, 3*time.Second)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// GetFileContent gets file content
func (fmc *FileManagerCommands) GetFileContent(filepath string) (string, error) {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`powershell -Command "Get-Content '%s' -Raw"`, filepath)
	default:
		cmd = fmt.Sprintf(`cat "%s"`, filepath)
	}

	return fmc.session.ExecuteCommand(cmd, 10*time.Second)
}

// DownloadFile downloads a file from the target using simple commands
func (fmc *FileManagerCommands) DownloadFile(remotePath, localPath string) error {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		// Use simple type command for Windows
		cmd = fmt.Sprintf(`type "%s"`, remotePath)
	default:
		// Use simple cat command for Unix
		cmd = fmt.Sprintf(`cat "%s"`, remotePath)
	}

	output, err := fmc.session.ExecuteCommand(cmd, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Write output directly to local file
	return ioutil.WriteFile(localPath, []byte(output), 0644)
}

// UploadFile uploads a file to the target using simple echo commands
func (fmc *FileManagerCommands) UploadFile(localPath, remotePath string) error {
	// Read local file
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	// For small files, use simple echo command
	// For larger files, we'll implement chunked upload
	if len(data) < 1000 {
		os := fmc.session.DetectedOS()
		var cmd string

		// Escape special characters for shell
		content := string(data)
		content = strings.ReplaceAll(content, `"`, `\"`)
		content = strings.ReplaceAll(content, `'`, `\'`)
		content = strings.ReplaceAll(content, `$`, `\$`)

		switch os {
		case "windows":
			cmd = fmt.Sprintf(`echo "%s" > "%s"`, content, remotePath)
		default:
			cmd = fmt.Sprintf(`echo '%s' > "%s"`, content, remotePath)
		}

		_, err = fmc.session.ExecuteCommand(cmd, 10*time.Second)
		return err
	}

	// For larger files, return error for now
	return fmt.Errorf("file too large for simple upload (max 1KB), use alternative method")
}

// CreateDirectory creates a new directory
func (fmc *FileManagerCommands) CreateDirectory(path string) error {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`mkdir "%s"`, path)
	default:
		cmd = fmt.Sprintf(`mkdir -p "%s"`, path)
	}

	_, err := fmc.session.ExecuteCommand(cmd, 5*time.Second)
	return err
}

// RemoveFile removes a file or directory
func (fmc *FileManagerCommands) RemoveFile(path string) error {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`Remove-Item -Path "%s" -Recurse -Force`, path)
	default:
		cmd = fmt.Sprintf(`rm -rf "%s"`, path)
	}

	_, err := fmc.session.ExecuteCommand(cmd, 10*time.Second)
	return err
}

// CopyFile copies a file or directory
func (fmc *FileManagerCommands) CopyFile(source, destination string) error {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`Copy-Item -Path "%s" -Destination "%s" -Recurse`, source, destination)
	default:
		cmd = fmt.Sprintf(`cp -r "%s" "%s"`, source, destination)
	}

	_, err := fmc.session.ExecuteCommand(cmd, 15*time.Second)
	return err
}

// MoveFile moves/renames a file
func (fmc *FileManagerCommands) MoveFile(source, destination string) error {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`Move-Item -Path "%s" -Destination "%s"`, source, destination)
	default:
		cmd = fmt.Sprintf(`mv "%s" "%s"`, source, destination)
	}

	_, err := fmc.session.ExecuteCommand(cmd, 10*time.Second)
	return err
}

// GetDiskUsage gets disk usage information
func (fmc *FileManagerCommands) GetDiskUsage(path string) (used, free, percent string, err error) {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`powershell -Command "$drive = Get-PSDrive -Name (Get-Item '%s').PSDrive.Name; Write-Host \"$($drive.Used)|$($drive.Free)|$([math]::Round(($drive.Free/($drive.Used+$drive.Free))*100,2))\""`, path)
	default:
		cmd = fmt.Sprintf(`df -h "%s" 2>/dev/null | awk 'NR==2 {print $3"|"$4"|"$5}'`, path)
	}

	output, err := fmc.session.ExecuteCommand(cmd, 5*time.Second)
	if err != nil {
		return "", "", "", err
	}

	parts := strings.Split(strings.TrimSpace(output), "|")
	if len(parts) >= 3 {
		return parts[0], parts[1], parts[2], nil
	}

	return "", "", "", fmt.Errorf("invalid output format")
}

// FormatSize formats file size in human-readable format
func (fmc *FileManagerCommands) FormatSize(size int64) string {
	if size == 0 {
		return ""
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for size >= div && exp < 3 {
		div *= unit
		exp++
	}

	switch exp {
	case 1:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(unit))
	case 2:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(unit*unit))
	case 3:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(unit*unit*unit))
	default:
		return fmt.Sprintf("%.1f TB", float64(size)/float64(unit*unit*unit*unit))
	}
}

// SanitizePath sanitizes file path for safe usage
func (fmc *FileManagerCommands) SanitizePath(path string) string {
	// Remove any dangerous characters
	dangerous := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\n", "\r"}
	result := path
	for _, char := range dangerous {
		result = strings.ReplaceAll(result, char, "")
	}
	return filepath.Clean(result)
}

// IsValidPath checks if path is valid and safe
func (fmc *FileManagerCommands) IsValidPath(path string) bool {
	if path == "" {
		return false
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") && !strings.HasPrefix(path, "/") {
		return false
	}

	// Check for dangerous characters
	dangerous := []string{";", "&", "|", "`", "$"}
	for _, char := range dangerous {
		if strings.Contains(path, char) {
			return false
		}
	}

	return true
}
