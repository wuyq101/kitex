package selfupdate

import (
	"context"
	"os"
	"time"

	"github.com/cloudwego/kitex/tool/pkg/log"

	"code.byted.org/kitex/updater"
)

// Start updates the an executable of the given name and returns an bool to indicate
// whether the update process is done succesfully.
func Start(name, version string, restart bool) (ok bool) {
	done := make(chan bool)
	go func() {
		done <- start(name, version, restart)
	}()
	select {
	case ok = <-done:
		return
	case <-time.After(time.Second * 30):
		return false
	}
}

// Update simply checks versions and updates the specified binary to the target file.
func Update(name, currentVersion, targetFile string) (updated bool, latestVersion string, err error) {
	updated, manifest, err := updater.Update(context.Background(), name, currentVersion, targetFile)
	if err == nil {
		return updated, manifest.LatestVersion, nil
	}
	return false, "", err
}

func start(name, nowVersion string, restart bool) (ok bool) {
	httpProxy := os.Getenv("http_proxy")
	httpsProxy := os.Getenv("http_proxy")
	os.Setenv("http_proxy", "")
	os.Setenv("https_proxy", "")
	defer os.Setenv("http_proxy", httpProxy)
	defer os.Setenv("https_proxy", httpsProxy)

	onUpdate := func(ctx context.Context, manifest updater.Manifest) bool {
		log.Warnf("Updated %s to %v\n", name, manifest.LatestVersion)
		return restart
	}
	_, _, err := updater.SelfUpdateAndRestart(context.Background(), name, nowVersion, onUpdate)
	if err != nil {
		log.Warn("self update failed:", err.Error())
		return false
	}
	return true
}
