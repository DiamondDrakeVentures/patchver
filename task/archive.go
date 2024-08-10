package task

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/DiamondDrakeVentures/patchver/common"
)

type ArchiveTask struct {
	id   string
	name string

	target    string
	sourceDir string
}

func (t ArchiveTask) ID() string {
	if t.id != "" {
		return t.id
	}

	return common.ID()
}

func (t *ArchiveTask) SetID(id string) {
	t.id = id
}

func (t ArchiveTask) Type() string {
	return "Archive"
}

func (t ArchiveTask) Name() string {
	if t.name != "" {
		return t.name
	}

	return t.ID()
}

func (t *ArchiveTask) SetName(name string) {
	t.name = name
}

func (t ArchiveTask) Execute(logger common.Logger) error {
	arF, err := os.Create(t.target)
	if err != nil {
		logger.Printf("Cannot create archive %s: %v\n", t.target, err)
		return err
	}
	defer arF.Close()

	ar := zip.NewWriter(arF)
	defer ar.Close()

	logger.Printf("Archiving %s\n", t.target)
	filepath.WalkDir(t.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Printf("Error reading path %s: %v\n", path, err)
			return err
		}

		filePath, err := filepath.Rel(t.sourceDir, path)
		if err != nil {
			logger.Printf("Cannot create relative path from %s: %v\n", path, err)
			return err
		}
		logger.Printf("  %s\n", filePath)

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			logger.Printf("Cannot probe file %s: %v\n", path, err)
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			logger.Printf("Cannot create header for file %s: %v\n", path, err)
			return err
		}
		header.Name = filePath

		dst, err := ar.CreateHeader(header)
		if err != nil {
			logger.Printf("Cannot create file %s: %v\n", path, err)
			return err
		}

		src, err := os.Open(path)
		if err != nil {
			logger.Printf("Cannot add file %s: %v\n", path, err)
			return err
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			logger.Printf("Cannot add file %s: %v\n", path, err)
			return err
		}

		return nil
	})
	logger.Printf("Done archiving %s\n", t.target)

	return nil
}

func Archive(target, sourceDir string) Task {
	return &ArchiveTask{
		target:    target,
		sourceDir: sourceDir,
	}
}
