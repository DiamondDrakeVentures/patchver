package main

var args struct {
	Config    string `arg:"-c,--config" default:"config.json" help:"Configuration file"`
	TargetDir string `arg:"positional" default:"."`
}
