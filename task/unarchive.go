package task

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/DiamondDrakeVentures/patchver/common"
)

type UnarchiveTask struct {
	id   string
	name string

	filename  string
	targetDir string
}

func (t UnarchiveTask) ID() string {
	if t.id != "" {
		return t.id
	}

	return common.ID()
}

func (t *UnarchiveTask) SetID(id string) {
	t.id = id
}

func (t UnarchiveTask) Type() string {
	return "Unarchive"
}

func (t UnarchiveTask) Name() string {
	if t.name != "" {
		return t.name
	}

	return t.ID()
}

func (t *UnarchiveTask) SetName(name string) {
	t.name = name
}

func (t UnarchiveTask) Execute(logger common.Logger) error {
	ar, err := zip.OpenReader(t.filename)
	if err != nil {
		logger.Printf("Cannot open %s: %v\n", t.filename, err)
		return err
	}
	defer ar.Close()

	logger.Printf("Extracting %s:\n", t.filename)
	for _, file := range ar.File {
		filePath := filepath.Join(t.targetDir, file.Name)

		logger.Printf("  %s\n", file.Name)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				logger.Printf("Cannot create directory %s: %v\n", filePath, err)
				return err
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			logger.Printf("Cannot create directory %s: %v\n", filePath, err)
			return err
		}

		dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			logger.Printf("Cannot create file %s: %v\n", filePath, err)
			return err
		}
		defer dst.Close()

		src, err := file.Open()
		if err != nil {
			logger.Printf("Cannot extract file %s: %v\n", filepath.Join(t.filename, file.Name), err)
			return err
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			logger.Printf("Cannot extract file %s: %v\n", filepath.Join(t.filename, file.Name), err)
			return err
		}
	}
	logger.Printf("Done extracting %s\n", t.filename)

	return nil
}

func Unarchive(filename, targetDir string) Task {
	return &UnarchiveTask{
		filename:  filename,
		targetDir: targetDir,
	}
}
