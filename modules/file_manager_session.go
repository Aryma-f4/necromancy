package modules

import (
	"bufio"
	"encoding/base64"
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

	// Send command
	_, err := fms.Write([]byte(cmd + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to send command: %v", err)
	}

	// Collect output
	var output strings.Builder
	deadline := time.After(timeout)

	for {
		select {
		case line := <-fms.commandOutput:
			output.WriteString(line + "\n")
		case <-deadline:
			return output.String(), nil
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
	// PowerShell command to get detailed file information
	return fmt.Sprintf(`powershell -Command "try { $files = Get-ChildItem -Path '%s' -Force -ErrorAction Stop; foreach ($file in $files) { Write-Host \"$($file.Name)|$($file.Length)|$($file.Mode)|$($file.LastWriteTime.ToString('yyyy-MM-dd HH:mm'))|$($file.PSIsContainer)|$($file.Attributes -band [System.IO.FileAttributes]::ReparsePoint)\" } } catch { Write-Host \"ERROR: $($_.Exception.Message)\" }"`, path)
}

// getLinuxFileListCommand generates Linux-specific file listing
func (fmc *FileManagerCommands) getLinuxFileListCommand(path string) string {
	// Enhanced ls command with proper parsing
	return fmt.Sprintf(`ls -la --time-style=long-iso "%s" 2>/dev/null | awk 'NR>1 {if (NF>=9) {name=""; for(i=9;i<=NF;i++) name=name (i>9?" ":"") $i; printf "%%s|%%s|%%s|%%s|%%s|%%s\n", name, $5, $1, $6" "$7, ($1 ~ /^d/ ? "true" : "false"), ($1 ~ /^l/ ? "true" : "false")}}'`, path)
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

// parseWindowsFileList parses Windows PowerShell output
func (fmc *FileManagerCommands) parseWindowsFileList(output string) []FileInfo {
	var files []FileInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ERROR:") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 6 {
			size := int64(0)
			fmt.Sscanf(parts[1], "%d", &size)

			files = append(files, FileInfo{
				Name:    parts[0],
				Size:    size,
				Mode:    parts[2],
				ModTime: parts[3],
				IsDir:   parts[4] == "true",
				IsLink:  parts[5] == "True",
			})
		}
	}

	return files
}

// parseLinuxFileList parses Linux ls output
func (fmc *FileManagerCommands) parseLinuxFileList(output string) []FileInfo {
	var files []FileInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 6 {
			size := int64(0)
			fmt.Sscanf(parts[1], "%d", &size)

			files = append(files, FileInfo{
				Name:    parts[0],
				Size:    size,
				Mode:    parts[2],
				ModTime: parts[3],
				IsDir:   parts[4] == "true",
				IsLink:  parts[5] == "true",
			})
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
		cmd = `powershell -Command "(Get-Location).Path"`
	default:
		cmd = `pwd`
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

// DownloadFile downloads a file from the target (base64 encoded)
func (fmc *FileManagerCommands) DownloadFile(remotePath, localPath string) error {
	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		cmd = fmt.Sprintf(`powershell -Command "[Convert]::ToBase64String([IO.File]::ReadAllBytes('%s'))"`, remotePath)
	default:
		cmd = fmt.Sprintf(`base64 "%s"`, remotePath)
	}

	output, err := fmc.session.ExecuteCommand(cmd, 30*time.Second)
	if err != nil {
		return err
	}

	// Decode base64 output
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(output))
	if err != nil {
		return fmt.Errorf("failed to decode base64: %v", err)
	}

	// Write to local file
	return ioutil.WriteFile(localPath, decoded, 0644)
}

// UploadFile uploads a file to the target (base64 decode)
func (fmc *FileManagerCommands) UploadFile(localPath, remotePath string) error {
	// Read local file
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	os := fmc.session.DetectedOS()
	var cmd string

	switch os {
	case "windows":
		// Split into chunks for Windows PowerShell
		chunkSize := 1000
		var chunks []string
		for i := 0; i < len(encoded); i += chunkSize {
			end := i + chunkSize
			if end > len(encoded) {
				end = len(encoded)
			}
			chunks = append(chunks, encoded[i:end])
		}

		cmd = fmt.Sprintf(`powershell -Command "$data = ''; %s; [IO.File]::WriteAllBytes('%s', [Convert]::FromBase64String($data))"`,
			buildPowerShellChunks(chunks), remotePath)
	default:
		cmd = fmt.Sprintf(`echo '%s' | base64 -d > "%s"`, encoded, remotePath)
	}

	_, err = fmc.session.ExecuteCommand(cmd, 30*time.Second)
	return err
}

// buildPowerShellChunks builds PowerShell command for large base64 data
func buildPowerShellChunks(chunks []string) string {
	var result strings.Builder
	for i, chunk := range chunks {
		if i == 0 {
			result.WriteString(fmt.Sprintf(`$data = '%s'`, chunk))
		} else {
			result.WriteString(fmt.Sprintf(`; $data += '%s'`, chunk))
		}
	}
	return result.String()
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
