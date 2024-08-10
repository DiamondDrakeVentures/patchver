package task

import (
	"errors"
	"os"

	"github.com/DiamondDrakeVentures/patchver/common"
)

type CleanupTask struct {
	id   string
	name string

	target string
}

func (t CleanupTask) ID() string {
	if t.id != "" {
		return t.id
	} else {
		return common.ID()
	}
}

func (t *CleanupTask) SetID(id string) {
	t.id = id
}

func (t CleanupTask) Type() string {
	return "Cleanup"
}

func (t CleanupTask) Name() string {
	if t.name != "" {
		return t.name
	}

	return t.ID()
}

func (t *CleanupTask) SetName(name string) {
	t.name = name
}

func (t CleanupTask) Execute(logger common.Logger) error {
	logger.Printf("Removing %s\n", t.target)
	if err := os.RemoveAll(t.target); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Printf("Cannot remove %s: %v\n", t.target, err)
			return err
		}
	}
	logger.Printf("Done removing %s\n", t.target)

	return nil
}

func Cleanup(target string) Task {
	return &CleanupTask{
		target: target,
	}
}
