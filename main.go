package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"snmp-wrapper/helpers"
	"time"
)

func main() {

	logPath := os.Getenv("LOGPATH_SNMP_WRAPPER")
	if logPath != "" {
		helpers.LogPath = fmt.Sprintf("%s/snmp-wrapper.log",logPath)
	}
	fmt.Printf("set log path: %s", helpers.LogPath)
	// log config
	logFile, err_log := os.OpenFile(helpers.LogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err_log != nil {
		log.Fatalln(err_log)
	}

	// defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// pars flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\tExample: snmp-wrapper -l 127.0.0.1 -p 1611 -RHost 192.168.1.1:8081 -U mikrotik -P mikrotik -i 10\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	listenIP := flag.String("l", "127.0.0.1", "UDP host IP address to listen on for the SNMP server")
	listenPort := flag.Int("p", 1161, "UDP port to listen on for the SNMP server")
	mikrotikHost := flag.String("RHost", "", "Mikrotik API server URL (host:port)")
	mikrotikUsername := flag.String("U", "", "Username for Mikrotik API authentication")
	mikrotikPassword := flag.String("P", "", "Password for Mikrotik API authentication")
	VersionAPP := flag.Bool("v", false, "Display the application version information")
	interval := flag.Int("i", 10, "Set the monitoring interval in seconds.")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if *VersionAPP {
		fmt.Printf("Version: %s\n", helpers.VersionAPP)
		os.Exit(0)
	}

	if *mikrotikHost == "" {
		log.Println("Error: Mikrotik API server URL (RHost) is required.")
		flag.Usage()
		os.Exit(1)
	}

	if *mikrotikUsername == "" || *mikrotikPassword == "" {
		log.Println("Warning: Username/Password is empty.")
		os.Exit(2)
	}

	if *interval > 0 {
		helpers.RepeatInterval = time.Duration(*interval) * time.Second
		log.Printf("Monitoring at an interval of %d seconds.", *interval)
	} else {
		log.Println("Invalid interval: must be greater than 0.")
	}

	p := helpers.NewInfoProcess(*listenIP, *mikrotikHost, *mikrotikUsername, *mikrotikPassword, *listenPort)
	p.Process()

}
