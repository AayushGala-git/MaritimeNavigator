package common

import (
	"log"
	"os"
)

var (
	// Logger GlobalLogger
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[MaritimeTracking] ", log.LstdFlags)
}
