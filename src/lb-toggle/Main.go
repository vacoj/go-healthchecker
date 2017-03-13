package main

import (
	"fmt"
)

func showVersion() {
	name := "LB-Toggle " + VERSION
	fmt.Println(name)
}

func main() {
	SETTINGS.parseSettingsFile()

	// Initialize the wait group so threads don't exit
	WG.Add(1)

	// Monitor application for health status
	STATUS.startMonitor()

	// Start the Web application.
	go startWeb()

	WG.Wait()
}
