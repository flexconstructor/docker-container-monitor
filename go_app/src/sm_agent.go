package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
	"system"
)

// Application for collect information about media server system.
// Such as CPU usage, Memory usage, process info.
// The application runs as daemon.

// Init flags.
var (
	signal = flag.String("s", "", `send signal to the daemon
		quit — graceful shutdown
		stop — fast shutdown
		reload — reloading the configuration file`)
)

// Create new system monitor instance.
var (
	monitor = system.NewSystemMonitor("/go/logs/system_monitor.sock")
)

// Main function.
// Parse flags.
// Init daemon context.
// Start daemon.
func main() {
	flag.Parse()
	daemon.AddCommand(
		daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	cntxt := &daemon.Context{
		PidFileName: "system_monitor.pid",
		PidFilePerm: 0644,
		LogFileName: "/go/logs/system_monitor.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[system monitor]"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("system monitor daemon started")

	go monitor.Run()
	defer monitor.Stop()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("system monitor daemon terminated")
}

// Terminate daemon.
func termHandler(sig os.Signal) error {
	log.Println("terminating system monitor...")
	monitor.Stop()
	return daemon.ErrStop
}
