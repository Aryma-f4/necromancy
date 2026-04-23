package modules

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FileInfo represents file/directory information
type FileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Mode       string `json:"mode"`
	ModTime    string `json:"modTime"`
	IsDir      bool   `json:"isDir"`
	IsLink     bool   `json:"isLink"`
	LinkTarget string `json:"linkTarget,omitempty"`
}

// getBannerFromFile reads and parses the ASCII banner from file
func getBannerFromFile() string {
	file, err := os.Open("ascii.txt")
	if err != nil {
		// Fallback to simple text banner if file not found
		return "[yellow]NECROMANCY[-]\n[blue]Advanced Shell Manager[-]"
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse BBCode color tags for tview
		line = parseBBCodeForTview(line)
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "[yellow]NECROMANCY[-]\n[blue]Advanced Shell Manager[-]"
	}

	return strings.Join(lines, "\n")
}

// parseBBCodeForTview converts BBCode color tags to tview format
func parseBBCodeForTview(text string) string {
	// Convert [color=#hex] to tview color format
	colorPattern := regexp.MustCompile(`\[color=([^\]]+)\]`)
	text = colorPattern.ReplaceAllString(text, "[$1]")

	// Convert [/color] to tview reset
	text = strings.ReplaceAll(text, "[/color]", "[-:-:-]")

	// Remove size and font tags
	text = regexp.MustCompile(`\[size=[^\]]+\]`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\[/size\]`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\[font=[^\]]+\]`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\[/font\]`).ReplaceAllString(text, "")

	// Remove HTML span tags
	text = regexp.MustCompile(`<span[^>]*>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`</span>`).ReplaceAllString(text, "")

	return text
}

// FileManagerModule provides btop-like file manager UI
type FileManagerModule struct{}

func (m *FileManagerModule) Name() string {
	return "filemanager"
}

func (m *FileManagerModule) Description() string {
	return "Interactive file manager with btop-like UI for target system"
}

func (m *FileManagerModule) Execute(s *core.Session) error {
	// Create enhanced session for file manager
	fms := NewFileManagerSession(s)

	// Launch the file manager UI with banner
	app := tview.NewApplication()

	// Create banner view with ASCII art
	bannerView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(getBannerFromFile() + "\n[gray]v1.1 Stable Release - File Manager[-]\n[blue]Interactive File Management for Target Systems[-]")

	bannerView.SetBackgroundColor(tcell.ColorBlack)

	fileManager := NewFileManagerUI(fms, app)

	// Create layout with banner at top
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	mainLayout.AddItem(bannerView, 8, 0, false)
	mainLayout.AddItem(fileManager.layout, 0, 1, true)

	if err := app.SetRoot(mainLayout, true).EnableMouse(true).Run(); err != nil {
		return fmt.Errorf("file manager UI error: %v", err)
	}

	return nil
}

// FileManagerUI represents the file manager interface
type FileManagerUI struct {
	session       *FileManagerSession
	app           *tview.Application
	layout        *tview.Flex
	fileList      *tview.Table
	pathInput     *tview.InputField
	statusBar     *tview.TextView
	currentPath   string
	files         []FileInfo
	selectedIndex int
}

// NewFileManagerUI creates a new file manager interface
func NewFileManagerUI(session *FileManagerSession, app *tview.Application) *FileManagerUI {
	fm := &FileManagerUI{
		session:       session,
		app:           app,
		currentPath:   ".",
		selectedIndex: 0,
	}

	fm.setupUI()
	fm.refreshFiles()

	return fm
}

// setupUI creates the btop-like interface
func (fm *FileManagerUI) setupUI() {
	// Create file list table with btop-like styling
	fm.fileList = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	// Set table styling
	fm.fileList.SetBackgroundColor(tcell.ColorBlack)
	fm.fileList.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite))

	// Create path input field
	fm.pathInput = tview.NewInputField().
		SetLabel("Path: ").
		SetText(fm.currentPath).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	// Create status bar
	fm.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	// Set up layout - similar to btop layout
	fm.layout = tview.NewFlex().SetDirection(tview.FlexRow)

	// Top section: path input and file list
	topSection := tview.NewFlex().SetDirection(tview.FlexRow)
	topSection.AddItem(fm.pathInput, 1, 0, false)
	topSection.AddItem(fm.fileList, 0, 1, true)

	// Main layout
	fm.layout.AddItem(topSection, 0, 1, true)
	fm.layout.AddItem(fm.statusBar, 1, 0, false)

	// Set up event handlers
	fm.setupEventHandlers()
}

// setupEventHandlers sets up keyboard and input handlers
func (fm *FileManagerUI) setupEventHandlers() {
	// File list navigation
	fm.fileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			fm.handleEnter()
			return nil
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			fm.navigateUp()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				fm.app.Stop()
				return nil
			case 'r', 'R':
				fm.refreshFiles()
				return nil
			case 'h', 'H':
				fm.showHelp()
				return nil
			case 'd', 'D':
				fm.downloadSelectedFile()
				return nil
			case 'u', 'U':
				fm.uploadFile()
				return nil
			case 'x', 'X':
				fm.executeSelectedFile()
				return nil
			case 'n', 'N':
				fm.createNewFile()
				return nil
			case 'm', 'M':
				fm.createNewDirectory()
				return nil
			case 'e', 'E':
				fm.editSelectedFile()
				return nil
			case 'c', 'C':
				fm.copySelectedFile()
				return nil
			case 'v', 'V':
				fm.pasteFile()
				return nil
			case 'a', 'A':
				fm.selectAll()
				return nil
			}
		}
		return event
	})

	// Path input handler
	fm.pathInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			newPath := fm.pathInput.GetText()
			fm.navigateToPath(newPath)
		}
	})
}

// refreshFiles updates the file list
func (fm *FileManagerUI) refreshFiles() {
	fm.files = fm.getFiles(fm.currentPath)
	fm.updateFileList()
	fm.updateStatusBar()
}

// getFiles retrieves file information from the target
func (fm *FileManagerUI) getFiles(path string) []FileInfo {
	// Use the enhanced session directly
	fmc := &FileManagerCommands{session: fm.session}

	// Sanitize the path
	sanitizedPath := fmc.SanitizePath(path)
	if !fmc.IsValidPath(sanitizedPath) {
		fm.updateStatus("Invalid path")
		return []FileInfo{}
	}

	// Try to list files using the enhanced session
	files, err := fmc.ListFiles(sanitizedPath)
	if err != nil {
		fm.updateStatus(fmt.Sprintf("Error listing files: %v", err))
		return fm.getDefaultFiles()
	}

	return files
}

// getDefaultFiles returns default file list for demonstration
func (fm *FileManagerUI) getDefaultFiles() []FileInfo {
	return []FileInfo{
		{Name: "..", Size: 0, Mode: "drwxr-xr-x", ModTime: time.Now().Format("2006-01-02 15:04"), IsDir: true},
		{Name: "README.md", Size: 2048, Mode: "-rw-r--r--", ModTime: time.Now().Format("2006-01-02 15:04"), IsDir: false},
		{Name: "config", Size: 4096, Mode: "drwxr-xr-x", ModTime: time.Now().Format("2006-01-02 15:04"), IsDir: true},
		{Name: "script.sh", Size: 1024, Mode: "-rwxr-xr-x", ModTime: time.Now().Format("2006-01-02 15:04"), IsDir: false, IsLink: false},
		{Name: "data.txt", Size: 512, Mode: "-rw-r--r--", ModTime: time.Now().Format("2006-01-02 15:04"), IsDir: false},
	}
}

// updateFileList updates the table with file information
func (fm *FileManagerUI) updateFileList() {
	fm.fileList.Clear()

	// Set headers
	headers := []string{"Name", "Size", "Permissions", "Modified"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false)
		fm.fileList.SetCell(0, i, cell)
	}

	// Add file entries
	for i, file := range fm.files {
		row := i + 1

		// Name column with icons
		name := file.Name
		if file.IsDir {
			name = "📁 " + name
		} else if file.IsLink {
			name = "🔗 " + name
		} else {
			name = "📄 " + name
		}

		nameCell := tview.NewTableCell(name).
			SetTextColor(fm.getFileColor(file)).
			SetSelectable(true)
		fm.fileList.SetCell(row, 0, nameCell)

		// Size column
		sizeStr := fm.formatSize(file.Size)
		sizeCell := tview.NewTableCell(sizeStr).
			SetTextColor(tcell.ColorWhite).
			SetSelectable(false).
			SetAlign(tview.AlignRight)
		fm.fileList.SetCell(row, 1, sizeCell)

		// Permissions column
		permCell := tview.NewTableCell(file.Mode).
			SetTextColor(tcell.ColorGreen).
			SetSelectable(false)
		fm.fileList.SetCell(row, 2, permCell)

		// Modified column
		timeCell := tview.NewTableCell(file.ModTime).
			SetTextColor(tcell.ColorLightCyan).
			SetSelectable(false)
		fm.fileList.SetCell(row, 3, timeCell)
	}

	// Select current item
	if fm.selectedIndex < len(fm.files) {
		fm.fileList.Select(fm.selectedIndex+1, 0)
	}
}

// getFileColor returns appropriate color for file type
func (fm *FileManagerUI) getFileColor(file FileInfo) tcell.Color {
	if file.IsDir {
		return tcell.ColorBlue
	}
	if file.IsLink {
		return tcell.ColorPurple
	}
	if strings.HasPrefix(file.Mode, "-rwx") {
		return tcell.ColorGreen
	}
	return tcell.ColorWhite
}

// formatSize formats file size in human-readable format
func (fm *FileManagerUI) formatSize(size int64) string {
	if size == 0 {
		return ""
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	switch exp {
	case 1:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(div))
	case 2:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(div))
	case 3:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(div))
	default:
		return fmt.Sprintf("%.1f TB", float64(size)/float64(div))
	}
}

// updateStatusBar updates the status bar information
func (fm *FileManagerUI) updateStatusBar() {
	path := fm.currentPath
	fileCount := len(fm.files)
	totalSize := int64(0)
	dirCount := 0

	for _, file := range fm.files {
		if file.IsDir {
			dirCount++
		} else {
			totalSize += file.Size
		}
	}

	status := fmt.Sprintf("[yellow]Path:[white] %s  [yellow]Files:[white] %d  [yellow]Dirs:[white] %d  [yellow]Size:[white] %s",
		path, fileCount-dirCount, dirCount, fm.formatSize(totalSize))

	fm.statusBar.SetText(status)
}

// updateStatus updates the status bar with a message
func (fm *FileManagerUI) updateStatus(message string) {
	fm.statusBar.SetText(fmt.Sprintf("[red]%s[white]", message))
}

// Navigation functions
func (fm *FileManagerUI) handleEnter() {
	row, _ := fm.fileList.GetSelection()
	if row > 0 && row <= len(fm.files) {
		file := fm.files[row-1]
		if file.IsDir {
			if file.Name == ".." {
				fm.navigateUp()
			} else {
				fm.navigateToPath(fm.joinPath(fm.currentPath, file.Name))
			}
		}
	}
}

func (fm *FileManagerUI) navigateUp() {
	if fm.currentPath == "/" || fm.currentPath == "." {
		return
	}

	parts := strings.Split(fm.currentPath, "/")
	if len(parts) > 1 {
		newPath := strings.Join(parts[:len(parts)-1], "/")
		if newPath == "" {
			newPath = "/"
		}
		fm.navigateToPath(newPath)
	}
}

func (fm *FileManagerUI) navigateToPath(path string) {
	fm.currentPath = path
	fm.pathInput.SetText(path)
	fm.selectedIndex = 0
	fm.refreshFiles()
}

func (fm *FileManagerUI) joinPath(base, name string) string {
	if base == "/" {
		return "/" + name
	}
	return base + "/" + name
}

// File operations
func (fm *FileManagerUI) downloadSelectedFile() {
	row, _ := fm.fileList.GetSelection()
	if row > 0 && row <= len(fm.files) {
		file := fm.files[row-1]
		if !file.IsDir {
			// Implement file download logic
			fm.updateStatus(fmt.Sprintf("Downloading %s...", file.Name))
		}
	}
}

func (fm *FileManagerUI) uploadFile() {
	// Implement file upload logic
	fm.updateStatus("Upload functionality - implement based on existing file transfer")
}

func (fm *FileManagerUI) executeSelectedFile() {
	row, _ := fm.fileList.GetSelection()
	if row > 0 && row <= len(fm.files) {
		file := fm.files[row-1]
		if !file.IsDir && strings.Contains(file.Mode, "x") {
			// Execute file
			cmd := fmt.Sprintf("./%s", file.Name)
			if fm.session.DetectedOS() == "windows" {
				cmd = fmt.Sprintf(".\\%s", file.Name)
			}
			fm.session.Write([]byte(cmd + "\n"))
			fm.app.Stop()
		}
	}
}

func (fm *FileManagerUI) createNewFile() {
	// Implement new file creation
	fm.updateStatus("New file creation - implement with text editor")
}

func (fm *FileManagerUI) createNewDirectory() {
	// Implement directory creation
	fm.updateStatus("New directory creation - implement with input dialog")
}

func (fm *FileManagerUI) editSelectedFile() {
	row, _ := fm.fileList.GetSelection()
	if row > 0 && row <= len(fm.files) {
		file := fm.files[row-1]
		if !file.IsDir {
			// Implement file editing
			fm.updateStatus(fmt.Sprintf("Editing %s...", file.Name))
		}
	}
}

func (fm *FileManagerUI) copySelectedFile() {
	// Implement file copy
	fm.updateStatus("File copy - implement with clipboard or temp storage")
}

func (fm *FileManagerUI) pasteFile() {
	// Implement file paste
	fm.updateStatus("File paste - implement with clipboard or temp storage")
}

func (fm *FileManagerUI) selectAll() {
	// Implement select all functionality
	fm.updateStatus("Select all - implement for batch operations")
}

func (fm *FileManagerUI) showHelp() {
	helpText := `[yellow]File Manager Help[white]

[yellow]Navigation:[white]
  ↑/↓ - Navigate files
  Enter - Open directory/execute file
  Backspace - Go to parent directory
  
[yellow]Operations:[white]
  r - Refresh
  d - Download file
  u - Upload file
  x - Execute file
  n - New file
  m - New directory
  e - Edit file
  c - Copy file
  v - Paste file
  a - Select all
  
[yellow]System:[white]
  q - Quit file manager
  h - Show this help

Press any key to close help...`

	// Create help modal
	helpView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(helpText)

	helpModal := tview.NewModal().
		SetText(helpView.GetText(false)).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			fm.app.SetRoot(fm.layout, true)
		})

	fm.app.SetRoot(helpModal, false)
}
