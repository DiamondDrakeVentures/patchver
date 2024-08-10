package task

import (
	"io"
	"os"

	"github.com/DiamondDrakeVentures/patchver/common"
	"github.com/DiamondDrakeVentures/patchver/fabric"
)

type ModifyTask struct {
	id   string
	name string

	filename   string
	targetFile string
	backupFile string

	replacers map[string]string
}

func (t ModifyTask) ID() string {
	if t.id != "" {
		return t.id
	}

	return common.ID()
}

func (t *ModifyTask) SetID(id string) {
	t.id = id
}

func (t ModifyTask) Type() string {
	return "Modify"
}

func (t ModifyTask) Name() string {
	if t.name != "" {
		return t.name
	}

	return t.ID()
}

func (t *ModifyTask) SetName(name string) {
	t.name = name
}

func (t *ModifyTask) Execute(logger common.Logger) error {
	err := os.Rename(t.filename, t.backupFile)
	if err != nil {
		logger.Printf("Cannot backup file %s to %s: %v\n", t.filename, t.backupFile, err)
		return err
	}

	srcFile, err := os.Open(t.backupFile)
	if err != nil {
		logger.Printf("Cannot read file %s: %v\n", t.backupFile, err)
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		logger.Printf("Cannot probe file %s: %v\n", t.backupFile, err)
		return err
	}

	src, err := io.ReadAll(srcFile)
	if err != nil {
		logger.Printf("Cannot read file %s: %v\n", t.backupFile, err)
		return err
	}

	man, err := fabric.FromJSON(string(src))
	if err != nil {
		logger.Printf("Cannot parse JSON: %v\n", err)
		return err
	}

	for key, val := range t.replacers {
		if man.Depends(key, val) {
			logger.Printf("Replacing %s with %s in .%s", key, val, "depends")
			continue
		}
		if man.Recommends(key, val) {
			logger.Printf("Replacing %s with %s in .%s", key, val, "recommends")
			continue
		}
		if man.Suggests(key, val) {
			logger.Printf("Replacing %s with %s in .%s", key, val, "suggests")
			continue
		}
		if man.Breaks(key, val) {
			logger.Printf("Replacing %s with %s in .%s", key, val, "breaks")
			continue
		}
		if man.Conflicts(key, val) {
			logger.Printf("Replacing %s with %s in .%s", key, val, "conflicts")
			continue
		}
	}

	dst, err := os.OpenFile(t.targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		logger.Printf("Cannot create file %s: %v\n", t.targetFile, err)
		return err
	}
	defer dst.Close()

	j, err := man.JSON()
	if err != nil {
		logger.Printf("Cannot encode to JSON: %v\n", err)
		return err
	}

	_, err = dst.WriteString(j)
	if err != nil {
		logger.Printf("Cannot write to file %s: %v\n", t.targetFile, err)
		return err
	}

	return nil
}

func Modify(filename, targetFile, backupFile string, replacers map[string]string) Task {
	return &ModifyTask{
		filename:   filename,
		targetFile: targetFile,
		backupFile: backupFile,

		replacers: replacers,
	}
}
