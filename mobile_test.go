// Tests for mobile_test.go

package autoupdate

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/getlantern/golog"
	"github.com/stretchr/testify/assert"
)

var (
	updateURL = "https://update.getlantern.org/update"
)

type TestUpdater struct {
	log golog.Logger
	Updater
}

func (u *TestUpdater) Progress(percentage int) {
	u.log.Debugf("Current progress: %6.02d%%", percentage)
}

func TestCheckUpdateAvailable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping hitting updateserver in short mode.")
	}
	// test with an older version number
	doTestCheckUpdate(t, false, false, "2.2.0")
}

func TestCheckUpdateUnavailable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping hitting updateserver in short mode.")
	}
	// test with a blank version number
	doTestCheckUpdate(t, true, true, "")
	// test with a way newer version
	doTestCheckUpdate(t, true, false, "9.3.3")
}

// urlEmpty and shouldErr are booleans that indicate whether or not
// CheckMobileUpdate should return a blank url or non-nil error
func doTestCheckUpdate(t *testing.T, urlEmpty, shouldErr bool, version string) string {

	config := &Config{
		URL:            updateURL,
		CurrentVersion: version,
		PublicKey:      []byte(PackagePublicKey),
		OS:             "android",
		Arch:           "arm",
		Channel:        "stable",
	}

	url, err := CheckMobileUpdate(config)

	if shouldErr {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}

	if urlEmpty {
		assert.Empty(t, url)
	} else {
		assert.NotEmpty(t, url)
	}

	return url
}

func TestDoUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping hitting updateserver in short mode.")
	}

	url := doTestCheckUpdate(t, false, false, "2.2.0")
	if !assert.NotEmpty(t, url) {
		return
	}

	// create a temporary file to write the update to
	out, err := ioutil.TempFile(os.TempDir(), "update")
	assert.Nil(t, err)

	defer os.Remove(out.Name())

	testUpdater := &TestUpdater{
		log: golog.LoggerFor("update-mobile-test"),
	}

	// check for an invalid apk path destination
	err = UpdateMobile(url, "", testUpdater, nil)
	assert.NotNil(t, err)

	// check for a missing url
	err = doUpdateMobile("", out, testUpdater, nil)
	assert.NotNil(t, err)

	// successful update
	err = doUpdateMobile(url, out, testUpdater, nil)
	assert.Nil(t, err)

}
