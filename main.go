package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DiamondDrakeVentures/patchver/executor"
	"github.com/alexflint/go-arg"
)

const UserAgent = "PatchVer/0.1"

func main() {
	err := arg.Parse(&args)
	if err != nil {
		fmt.Printf("Cannot parse arguments!\n%v\n", err)
		os.Exit(1)
	}

	conf, err := parseConf(args.Config)
	if err != nil {
		fmt.Printf("Cannot read config.json: %v\n", err)
		os.Exit(1)
	}

	exec := executor.New(log.Default())

	err = os.MkdirAll(args.TargetDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Cannot create target directory %s!\n%v\n", args.TargetDir, err)
		os.Exit(1)
	}

	for _, mod := range conf.Mods {
		err := registerFabricTasks(exec, &mod, args.TargetDir, UserAgent)
		if err != nil {
			fmt.Printf("Cannot register tasks for %s: %v\n", mod.ID, err)
			os.Exit(1)
		}
	}

	err = exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
