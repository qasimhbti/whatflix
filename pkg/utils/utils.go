// Package utils provide common utilities for main package
package utils

import (
	"log"
	"os"
	"runtime"
)

// InitLog initializes log format
func InitLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

// LogStart logs application start
func LogStart(version, env string) {
	log.Println("Start")
	log.Printf("Version: %s", version)
	log.Printf("Environment: %s", env)
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go Max Procs: %d", runtime.GOMAXPROCS(0))
}
