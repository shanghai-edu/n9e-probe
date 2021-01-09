package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/freedomkk-qfeng/n9e-probe/config"
	"github.com/freedomkk-qfeng/n9e-probe/http"
	"github.com/freedomkk-qfeng/n9e-probe/probe"
	"github.com/freedomkk-qfeng/n9e-probe/probe/core"

	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"
)

var (
	vers *bool
	help *bool
	conf *string
)

func init() {
	vers = flag.Bool("v", false, "display the version.")
	help = flag.Bool("h", false, "print this help.")
	conf = flag.String("f", "", "specify configuration file.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", config.Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	aconf()
	pconf()
	start()

	cfg := config.Get()
	config.InitLog(cfg.Logger)

	core.InitRpcClients()
	probe.BuildMappers()
	probe.Collect()

	http.Start()
	ending()
}

// auto detect configuration file
func aconf() {
	if *conf != "" && file.IsExist(*conf) {
		return
	}

	*conf = "etc/probe.local.yml"
	if file.IsExist(*conf) {
		return
	}

	*conf = "etc/probe.yml"
	if file.IsExist(*conf) {
		return
	}

	fmt.Println("no configuration file for probe")
	os.Exit(1)
}

// parse configuration file
func pconf() {
	if err := config.Parse(*conf); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	}
}

func ending() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		fmt.Printf("stop signal caught, stopping... pid=%d\n", os.Getpid())
	}

	logger.Close()
	http.Shutdown()
	fmt.Println("probe stopped successfully")
}

func start() {
	runner.Init()
	fmt.Println("probe start, use configuration file:", *conf)
	fmt.Println("runner.cwd:", runner.Cwd)
}
