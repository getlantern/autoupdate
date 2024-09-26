package autoupdate

import (
	"bytes"
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/blang/semver"
)

type Updater interface {
	// PublishProgress: publish percentage of update already downloaded
	Progress(int)
}

// byteCounter wraps an existing io.Reader and keeps track of the byte
// count while downloading the latest update
type byteCounter struct {
	io.Reader // Underlying io.Reader to track bytes transferred
	Updater
	total    int64   // Total bytes transferred
	length   int64   // Expected length
	progress float64 // How much of the update has been downloaded
}

func (pt *byteCounter) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)
		pt.Updater.Progress(int(percentage))
	}
	return n, err
}

// CheckMobileUpdate checks if a new update is available for mobile.
func CheckMobileUpdate(cfg *Config) (string, error) {
	log.Debugf("Checking for new mobile version; current version: %s", cfg.CurrentVersion)

	res, err := cfg.check()
	if err != nil {
		log.Errorf("Error checking for update for mobile: %v", err)
		return "", err
	}

	if res == nil {
		log.Debugf("No new version available!")
		return "", nil
	}

	v, err := semver.Make(cfg.CurrentVersion)
	if err != nil {
		log.Errorf("Error checking for update; could not parse version number: %v", err)
		return "", err
	}

	if isNewerVersion(v, res.Version) {
		log.Debugf("Newer version of Lantern mobile available! %s at %s", res.Version, res.Url)
		return res.Url, nil
	}

	return "", nil
}

// UpdateMobile downloads the latest APK from the given url to file apkPath.
func UpdateMobile(url, apkPath string, updater Updater, httpClient *http.Client) error {
	out, err := os.Create(apkPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer out.Close()
	return doUpdateMobile(url, out, updater, httpClient)
}

var ErrInvalidStatusCode = errors.New("request returned unexpected status code")

func doUpdateMobile(url string, out io.Writer, updater Updater, httpClient *http.Client) error {
	var req *http.Request
	var res *http.Response
	var err error

	log.Debugf("Attempting to download APK from %s", url)

	if req, err = http.NewRequest("GET", url, http.NoBody); err != nil {
		log.Errorf("Error downloading update: %v", err)
		return err
	}

	req.Header.Add("Accept-Encoding", "gzip")

	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if res, err = httpClient.Do(req); err != nil {
		log.Errorf("Error requesting update: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Errorf("failed to download update, status code: ", res.StatusCode)
		return fmt.Errorf("%w: expected status 200 status code but got %d", ErrInvalidStatusCode, res.StatusCode)
	}

	// We use a special byteCounter that storres a reference
	// to the updater interface to make it easy to publish progress
	// for how much of the update has been downloaded already.
	bytespt := &byteCounter{Updater: updater,
		Reader: res.Body, length: res.ContentLength}

	contents, err := io.ReadAll(bytespt)
	if err != nil {
		log.Errorf("Error reading update: %v", err)
		return err
	}

	_, err = io.Copy(out, bzip2.NewReader(bytes.NewReader(contents)))
	if err != nil {
		log.Errorf("Error copying update: %v", err)
		return err
	}

	return nil
}
