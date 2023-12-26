package main

import (
	"flag"
	"github.com/aaronchen2k/deeptest/cmd/server/serverServe"
	"github.com/aaronchen2k/deeptest/internal/pkg/queue"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	AppVersion string
	BuildTime  string
	GoVersion  string
	GitHash    string
	flagSet    *flag.FlagSet
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	queue.InitServerMq()
	serverServe.Start()
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
