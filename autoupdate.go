// Package autoupdate provides Lantern with special tools to autoupdate itself
// with minimal effort.
package autoupdate

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blang/semver"
	"github.com/getlantern/go-update"
	"github.com/getlantern/go-update/check"
	"github.com/getlantern/golog"
)

var (
	log                  = golog.LoggerFor("autoupdate")
	defaultCheckInterval = time.Hour * 4
)

type Config struct {
	// CurrentVersion: the current version of the program, must be in the form
	// X.Y.Z
	CurrentVersion string

	version semver.Version

	// URL: the url at which to check for updates
	URL string

	// PublicKey: the public key against which to check the signature of any
	// received updates.
	PublicKey []byte

	// CheckInterval: the interval at which to check for updates, defaults to
	// 4 hours.
	CheckInterval time.Duration

	// HTTPClient: (optional), an http.Client to use when checking for updates
	HTTPClient *http.Client

	// Operating system (optional, will be inferred)
	OS string
	// Arch (optional, will be inferred)
	Arch string
	// Channel (optional, defaults to stable)
	Channel string
}

// ApplyNext applies the next available update whenever it is available, blocking
// until the next update has been applied. If ApplyNext returns without an
// error, that means that the current program's executable has been udpated in
// place and you may want to restart. If ApplyNext returns an error, that means
// that an unrecoverable error has occurred and we can't continue checking for
// updates.
func ApplyNext(cfg *Config) (newVersion string, err error) {
	// Parse the semantic version
	cfg.version, err = semver.Parse(cfg.CurrentVersion)
	if err != nil {
		return "", fmt.Errorf("Bad version string: %v", err)
	}
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = defaultCheckInterval
		log.Debugf("Defaulted CheckInterval to %v", cfg.CheckInterval)
	}
	return cfg.loop()
}

func (cfg *Config) loop() (string, error) {
	for {
		res, err := cfg.check()

		if err != nil {
			log.Errorf("Problem checking for update: %v", err)
		} else {
			if res == nil {
				log.Debug("No update available")
			} else if isNewerVersion(cfg.version, res.Version) {
				log.Debugf("Attempting to update to %s.", res.Version)
				err, errRecover := res.Update()
				if errRecover != nil {
					// This should never happen, if this ever happens it means bad news such as
					// a missing executable file.
					return "", fmt.Errorf("Failed to recover from failed update attempt: %v\n", errRecover)
				}
				if err == nil {
					log.Debugf("Patching succeeded!")
					return res.Version, nil
				}
				log.Errorf("Patching failed: %q\n", err)
			} else {
				log.Debug("Already up to date.")
			}
		}

		time.Sleep(cfg.CheckInterval)
	}
}

func isNewerVersion(version semver.Version, newer string) bool {
	nv, err := semver.Parse(newer)
	if err != nil {
		log.Errorf("Bad version string on update: %v", err)
		return false
	}
	return nv.GT(version)
}

// check uses go-update to look for updates.
func (cfg *Config) check() (res *check.Result, err error) {
	params := &check.Params{
		AppVersion: cfg.CurrentVersion,
		OS:         cfg.OS,
		Arch:       cfg.Arch,
		Channel:    cfg.Channel,
	}

	update.SetHttpClient(cfg.HTTPClient)

	up := update.New().ApplyPatch(update.PATCHTYPE_BSDIFF)

	if _, err = up.VerifySignatureWithPEM(cfg.PublicKey); err != nil {
		return nil, fmt.Errorf("Problem verifying signature of update: %v", err)
	}

	if res, err = params.CheckForUpdate(cfg.URL, up); err != nil {
		if err == check.ErrNoUpdateAvailable {
			return nil, nil
		}
		return nil, fmt.Errorf("Problem fetching update: %v", err)
	}

	return res, nil
}
