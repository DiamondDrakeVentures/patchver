package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Downloader struct {
	client    *http.Client
	UserAgent string
}

func (d *Downloader) Download(url, targetDir string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", d.UserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error when downloading downloading %s: %v", url, resp.Status)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filepath.Join(targetDir, path.Base(req.URL.Path)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Init(userAgent string) *Downloader {
	d := new(Downloader)
	d.client = &http.Client{}
	d.UserAgent = userAgent

	return d
}
