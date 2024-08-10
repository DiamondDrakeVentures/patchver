package main

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/DiamondDrakeVentures/patchver/executor"
	"github.com/DiamondDrakeVentures/patchver/job"
	"github.com/DiamondDrakeVentures/patchver/task"
)

func registerFabricTasks(exec *executor.Executor, mod *Mod, tgtDir, userAgent string) error {
	var filename string
	if mod.Filename != "" {
		filename = mod.Filename
	} else {
		u, err := url.Parse(mod.URL)
		if err != nil {
			return err
		}
		filename = path.Base(u.Path)
	}

	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	archivePath := filepath.Join(tgtDir, filename)
	unarchiveDir, err := os.MkdirTemp(tgtDir, baseFilename+"*")
	if err != nil {
		return err
	}

	srcMan := filepath.Join(unarchiveDir, "fabric.mod.json")
	bakMan := srcMan + ".original"
	tgtMan := srcMan

	tsks := []task.Task{
		task.Download(mod.URL, tgtDir, userAgent),
		task.Unarchive(archivePath, unarchiveDir),
		task.Modify(srcMan, tgtMan, bakMan, mod.Replacers),
		task.Archive(archivePath, unarchiveDir),
		task.Cleanup(unarchiveDir),
	}

	exec.Register(job.NewJob(mod.ID, tsks))

	return nil
}
