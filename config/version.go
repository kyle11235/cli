package config

import (
	"errors"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	// version
	Version    = "v0.0.1"
	releaseURL = "https://github.com/kyle11235/cli/releases"
	// releaseURL     = "https://github.com/fnproject/fn/releases"
)

func GetCurrentVersion() string {
	return Version
}

func GetLatestVersion() (string, error) {
	redirectURL := ""
	client := http.Client{}
	client.Timeout = time.Second * 3
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectURL = req.URL.String()
		return nil
	}
	lastestURL := releaseURL + "/latest"
	res, err := client.Get(lastestURL)
	if err != nil {
		return "", errors.New("error found when check " + lastestURL)
	}
	defer res.Body.Close()
	if !strings.Contains(redirectURL, releaseURL) {
		return "", errors.New("redirect is incorrect when check " + lastestURL)
	}
	latestVersion := path.Base(redirectURL)
	return latestVersion, nil
}
