package main

import "fmt"

type args struct {
	Config    string `arg:"-c,--config" default:"config.json" help:"Configuration file"`
	TargetDir string `arg:"positional" default:"."`
}

func (args) Version() string {
	return fmt.Sprintf("%s v%s", Name, Version)
}
