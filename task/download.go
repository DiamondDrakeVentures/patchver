package task

import (
	"errors"
	"os"

	"github.com/DiamondDrakeVentures/patchver/common"
	"github.com/DiamondDrakeVentures/patchver/task/downloader"
)

type DownloadTask struct {
	id   string
	name string

	url       string
	targetDir string
	userAgent string
}

func (t DownloadTask) ID() string {
	if t.id != "" {
		return t.id
	}

	return common.ID()
}

func (t *DownloadTask) SetID(id string) {
	t.id = id
}

func (t DownloadTask) Type() string {
	return "Download"
}

func (t DownloadTask) Name() string {
	if t.name != "" {
		return t.name
	}

	return t.ID()
}

func (t *DownloadTask) SetName(name string) {
	t.name = name
}

func (t DownloadTask) Execute(logger common.Logger) error {
	if _, err := os.Stat(t.targetDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(t.targetDir, os.ModePerm)
		if err != nil {
			logger.Printf("Cannot create directory %s: %v\n", t.targetDir, err)
			return err
		}
	}

	logger.Printf("Downloading %s to %s\n", t.url, t.targetDir)

	d := downloader.Init(t.userAgent)
	err := d.Download(t.url, t.targetDir)
	if err != nil {
		logger.Printf("Cannot download %s: %v\n", t.url, err)
		return err
	}

	return nil
}

func Download(url, targetDir, userAgent string) Task {
	return &DownloadTask{
		url:       url,
		targetDir: targetDir,
		userAgent: userAgent,
	}
}
