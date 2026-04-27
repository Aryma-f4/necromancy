package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/Aryma-f4/necromancy/pty"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ansiSequencePattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

type shellTab struct {
	id           int
	title        string
	transcript   *tview.TextView
	inputHistory []string
	historyIndex int
}

type sessionShellWorkspace struct {
	app       *App
	session   *core.Session
	sessionID int
	pageName  string

	root     *tview.Flex
	tabsBar  *tview.TextView
	tabPages *tview.Pages
	input    *tview.InputField
	status   *tview.TextView

	tabs      []*shellTab
	activeTab int
	nextTabID int
	outputSub chan []byte
}

func (a *App) openSessionShellTabs(id int) {
	session, exists := a.sessions.Get(id)
	if !exists {
		return
	}

	if workspace, ok := a.shellWorkspaces[id]; ok {
		a.pages.SwitchToPage(workspace.pageName)
		a.tviewApp.SetFocus(workspace.input)
		workspace.setStatus("Returned to the tabbed shell workspace")
		return
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
	a.pages.SwitchToPage(workspace.pageName)
	a.tviewApp.SetFocus(workspace.input)
}

func newSessionShellWorkspace(app *App, session *core.Session) *sessionShellWorkspace {
	ws := &sessionShellWorkspace{
		app:       app,
		session:   session,
		sessionID: session.ID,
		pageName:  fmt.Sprintf("shell_tabs_%d", session.ID),
	}

	ws.tabsBar = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false)
	ws.tabsBar.SetBackgroundColor(colorPanel)

	ws.tabPages = tview.NewPages()

	ws.input = tview.NewInputField().
		SetFieldBackgroundColor(colorPanelAlt).
		SetFieldTextColor(colorPrimary).
		SetLabelColor(colorSecondary)

	ws.status = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	ws.status.SetBackgroundColor(colorPanel)

	ws.root = tview.NewFlex().SetDirection(tview.FlexRow)
	ws.root.AddItem(ws.tabsBar, 1, 0, false)
	ws.root.AddItem(ws.tabPages, 0, 1, false)
	ws.root.AddItem(ws.input, 1, 0, true)
	ws.root.AddItem(ws.status, 1, 0, false)

	ws.newTab()
	ws.bindKeys()
	ws.attachOutput()
	ws.maybeUpgradePTY()
	ws.setStatus("Tabbed shell ready")

	return ws
}

func (ws *sessionShellWorkspace) bindKeys() {
	ws.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			ws.submitCurrentInput()
		}
	})

	handler := func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ws.app.pages.SwitchToPage("session_actions")
			return nil
		case tcell.KeyTAB:
			ws.switchTab(1)
			return nil
		case tcell.KeyBacktab:
			ws.switchTab(-1)
			return nil
		case tcell.KeyCtrlN:
			ws.newTab()
			return nil
		case tcell.KeyCtrlW:
			ws.closeCurrentTab()
			return nil
		case tcell.KeyUp:
			ws.historyUp()
			return nil
		case tcell.KeyDown:
			ws.historyDown()
			return nil
		case tcell.KeyEnter:
			ws.submitCurrentInput()
			return nil
		}
		return event
	}

	ws.root.SetInputCapture(handler)
	ws.input.SetInputCapture(handler)
}

func (ws *sessionShellWorkspace) attachOutput() {
	ws.outputSub = ws.session.Subscribe()

	go func() {
		for data := range ws.outputSub {
			payload := data
			ws.app.tviewApp.QueueUpdateDraw(func() {
				if payload == nil {
					ws.setStatus("Remote session closed")
					ws.appendSystemMessage("Remote host closed the session.")
					return
				}
				ws.appendRemoteOutput(string(payload))
			})
		}
	}()
}

func (ws *sessionShellWorkspace) maybeUpgradePTY() {
	if ws.session.Type == "PTY" || strings.EqualFold(ws.session.DetectedOS(), "windows") {
		return
	}

	ws.appendSystemMessage("Attempting PTY upgrade for better shell behavior...")
	go func() {
		err := pty.AutoUpgrade(ws.session)
		ws.app.tviewApp.QueueUpdateDraw(func() {
			if err != nil {
				ws.appendSystemMessage(fmt.Sprintf("PTY upgrade skipped: %v", err))
				return
			}
			ws.appendSystemMessage("PTY upgrade completed.")
		})
	}()
}

func (ws *sessionShellWorkspace) newTab() {
	ws.nextTabID++
	title := fmt.Sprintf("Tab %d", ws.nextTabID)

	view := tview.NewTextView().
		SetDynamicColors(false).
		SetScrollable(true).
		SetWrap(false)
	styleTextView(view, fmt.Sprintf(" %s ", title))

	tab := &shellTab{
		id:           ws.nextTabID,
		title:        title,
		transcript:   view,
		historyIndex: 0,
	}

	tabPageName := ws.tabPageName(tab.id)
	ws.tabPages.AddPage(tabPageName, view, true, true)
	ws.tabs = append(ws.tabs, tab)
	ws.activeTab = len(ws.tabs) - 1
	history := sanitizeTerminalOutput(string(ws.session.SnapshotHistory()))
	if strings.TrimSpace(history) != "" {
		ws.appendToTab(tab, history)
		if !strings.HasSuffix(history, "\n") {
			ws.appendToTab(tab, "\n")
		}
	}
	ws.appendToTab(tab, fmt.Sprintf("[local] %s created. This tab shares the same remote stream as the other tabs.\n\n", title))
	ws.refreshChrome()
	ws.tabPages.SwitchToPage(tabPageName)
	ws.input.SetLabel(fmt.Sprintf("%s > ", tab.title))
	ws.input.SetText("")
	ws.app.tviewApp.SetFocus(ws.input)
	ws.setStatus(fmt.Sprintf("%s created", title))
}

func (ws *sessionShellWorkspace) closeCurrentTab() {
	if len(ws.tabs) <= 1 {
		ws.setStatus("At least one shell tab must remain open")
		return
	}

	index := ws.activeTab
	tab := ws.tabs[index]
	ws.tabPages.RemovePage(ws.tabPageName(tab.id))
	ws.tabs = append(ws.tabs[:index], ws.tabs[index+1:]...)
	if index >= len(ws.tabs) {
		index = len(ws.tabs) - 1
	}
	ws.activeTab = index
	ws.refreshChrome()
	ws.tabPages.SwitchToPage(ws.tabPageName(ws.tabs[ws.activeTab].id))
	ws.input.SetLabel(fmt.Sprintf("%s > ", ws.tabs[ws.activeTab].title))
	ws.app.tviewApp.SetFocus(ws.input)
	ws.setStatus(fmt.Sprintf("%s closed", tab.title))
}

func (ws *sessionShellWorkspace) switchTab(direction int) {
	if len(ws.tabs) == 0 {
		return
	}

	ws.activeTab = (ws.activeTab + direction + len(ws.tabs)) % len(ws.tabs)
	active := ws.tabs[ws.activeTab]
	ws.tabPages.SwitchToPage(ws.tabPageName(active.id))
	ws.input.SetLabel(fmt.Sprintf("%s > ", active.title))
	ws.refreshChrome()
	ws.app.tviewApp.SetFocus(ws.input)
	ws.setStatus(fmt.Sprintf("Switched to %s", active.title))
}

func (ws *sessionShellWorkspace) submitCurrentInput() {
	tab := ws.currentTab()
	if tab == nil {
		return
	}

	command := ws.input.GetText()
	if strings.TrimSpace(command) != "" {
		tab.inputHistory = append(tab.inputHistory, command)
		tab.historyIndex = len(tab.inputHistory)
	}

	prompt := "$ "
	if strings.EqualFold(ws.session.DetectedOS(), "windows") {
		prompt = "PS> "
	}
	ws.appendToTab(tab, fmt.Sprintf("%s%s\n", prompt, command))

	if _, err := ws.session.Write([]byte(command + "\n")); err != nil {
		ws.setStatus(fmt.Sprintf("Send failed: %v", err))
		return
	}

	ws.input.SetText("")
	ws.setStatus(fmt.Sprintf("Command sent from %s", tab.title))
}

func (ws *sessionShellWorkspace) historyUp() {
	tab := ws.currentTab()
	if tab == nil || len(tab.inputHistory) == 0 {
		return
	}
	if tab.historyIndex > 0 {
		tab.historyIndex--
	}
	ws.input.SetText(tab.inputHistory[tab.historyIndex])
}

func (ws *sessionShellWorkspace) historyDown() {
	tab := ws.currentTab()
	if tab == nil || len(tab.inputHistory) == 0 {
		return
	}
	if tab.historyIndex < len(tab.inputHistory)-1 {
		tab.historyIndex++
		ws.input.SetText(tab.inputHistory[tab.historyIndex])
		return
	}
	tab.historyIndex = len(tab.inputHistory)
	ws.input.SetText("")
}

func (ws *sessionShellWorkspace) appendRemoteOutput(text string) {
	cleaned := sanitizeTerminalOutput(text)
	if cleaned == "" {
		return
	}
	for _, tab := range ws.tabs {
		ws.appendToTab(tab, cleaned)
	}
}

func (ws *sessionShellWorkspace) appendSystemMessage(message string) {
	if current := ws.currentTab(); current != nil {
		ws.appendToTab(current, fmt.Sprintf("[local] %s\n", message))
	}
	ws.setStatus(message)
}

func (ws *sessionShellWorkspace) appendToTab(tab *shellTab, text string) {
	fmt.Fprint(tab.transcript, text)
	tab.transcript.ScrollToEnd()
}

func (ws *sessionShellWorkspace) setStatus(message string) {
	ws.status.SetText(fmt.Sprintf("[#788195]%s[-]", message))
}

func (ws *sessionShellWorkspace) currentTab() *shellTab {
	if len(ws.tabs) == 0 || ws.activeTab < 0 || ws.activeTab >= len(ws.tabs) {
		return nil
	}
	return ws.tabs[ws.activeTab]
}

func (ws *sessionShellWorkspace) refreshChrome() {
	var parts []string
	for index, tab := range ws.tabs {
		if index == ws.activeTab {
			parts = append(parts, fmt.Sprintf("[#ffd25c] %s [#5c91ff]|[-]", tab.title))
		} else {
			parts = append(parts, fmt.Sprintf("[#98d6bd] %s [#5c91ff]|[-]", tab.title))
		}
	}

	ws.tabsBar.SetText(" " + strings.Join(parts, " "))
	if active := ws.currentTab(); active != nil {
		ws.input.SetLabel(fmt.Sprintf("%s > ", active.title))
	}
}

func (ws *sessionShellWorkspace) tabPageName(id int) string {
	return fmt.Sprintf("shell_tab_%d_%d", ws.sessionID, id)
}

func sanitizeTerminalOutput(text string) string {
	text = ansiSequencePattern.ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}
