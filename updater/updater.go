package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

// Current version information
var (
	CurrentVersion = "1.0.0" // Will be overridden by main package
	RepoOwner      = "Aryma-f4"
	RepoName       = "necromancy"
)

// SetVersion allows the main package to set the current version
func SetVersion(version string) {
	CurrentVersion = version
}

// CheckForUpdate checks if a new version is available
func CheckForUpdate() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &release, nil
}

// IsNewerVersion compares versions and returns true if remote is newer
func IsNewerVersion(current, remote string) bool {
	// Simple version comparison - remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	remote = strings.TrimPrefix(remote, "v")

	return remote > current
}

// Asset represents a release asset
type Asset struct {
	Name               string
	BrowserDownloadURL string
	Size               int64
}

// GetAssetForPlatform returns the appropriate asset for the current platform
func GetAssetForPlatform(release *GitHubRelease) *Asset {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go platform/arch to asset names
	assetSuffix := ""
	switch platform {
	case "linux":
		if arch == "amd64" {
			assetSuffix = "linux-amd64"
		} else if arch == "arm64" {
			assetSuffix = "linux-arm64"
		}
	case "darwin":
		if arch == "amd64" {
			assetSuffix = "macos-amd64"
		} else if arch == "arm64" {
			assetSuffix = "macos-arm64"
		}
	case "windows":
		if arch == "amd64" {
			assetSuffix = "windows-amd64.exe"
		}
	}

	if assetSuffix == "" {
		return nil
	}

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, assetSuffix) {
			return &Asset{
				Name:               asset.Name,
				BrowserDownloadURL: asset.BrowserDownloadURL,
				Size:               asset.Size,
			}
		}
	}

	return nil
}

// DownloadUpdate downloads the update binary
func DownloadUpdate(url, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "necromancy-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write update: %v", err)
	}

	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
			return fmt.Errorf("failed to make executable: %v", err)
		}
	}

	// Move temp file to destination
	if err := os.Rename(tmpFile.Name(), destination); err != nil {
		return fmt.Errorf("failed to move update: %v", err)
	}

	return nil
}

// AutoUpdate performs automatic update check and download
func AutoUpdate() error {
	fmt.Println("[+] Checking for updates...")

	release, err := CheckForUpdate()
	if err != nil {
		return fmt.Errorf("update check failed: %v", err)
	}

	if !IsNewerVersion(CurrentVersion, release.TagName) {
		fmt.Printf("[✓] You have the latest version (%s)\n", CurrentVersion)
		return nil
	}

	fmt.Printf("[!] New version available: %s (current: %s)\n", release.TagName, CurrentVersion)
	fmt.Printf("[i] Release notes: %s\n", release.Body)

	asset := GetAssetForPlatform(release)
	if asset == nil {
		return fmt.Errorf("no compatible asset found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Printf("[+] Downloading update: %s (%d bytes)\n", asset.Name, asset.Size)

	// Download to current executable location
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %v", err)
	}

	// Backup current executable
	backupPath := currentExe + ".backup"
	if err := os.Rename(currentExe, backupPath); err != nil {
		return fmt.Errorf("failed to backup current executable: %v", err)
	}

	// Download new version
	if err := DownloadUpdate(asset.BrowserDownloadURL, currentExe); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentExe)
		return fmt.Errorf("failed to download update: %v", err)
	}

	fmt.Println("[✓] Update downloaded successfully!")
	fmt.Println("[i] Please restart the application to use the new version.")

	return nil
}

// CheckAndNotify checks for updates and notifies user
func CheckAndNotify() {
	release, err := CheckForUpdate()
	if err != nil {
		fmt.Printf("[-] Update check failed: %v\n", err)
		return
	}

	if !IsNewerVersion(CurrentVersion, release.TagName) {
		return // No update available
	}

	fmt.Printf("[!] Update available: %s → %s\n", CurrentVersion, release.TagName)
	fmt.Printf("[i] Release: %s\n", release.Name)
	fmt.Println("[i] Run with --update flag to download and install")
}
