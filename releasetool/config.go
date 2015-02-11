package main

import (
	"encoding/base64"
	"fmt"
)

const (
	endpointPrefix   = "https://api.equinox.io"
	endpointAssets   = endpointPrefix + "/1/Applications/%s/Assets"
	endpointReleases = endpointPrefix + "/1/Applications/%s/Releases"
)

type Config struct {
	AccountID     string `yaml:"account_id"`
	SecretKey     string `yaml:"secret_key"`
	ApplicationID string `yaml:"application_id"`
	Channel       string `yaml:"channel"`
	PrivateKey    string `yaml:"private_key"`
}

func (c Config) AuthHeader() string {
	s := fmt.Sprintf("%s:%s", c.AccountID, c.SecretKey)
	return base64.StdEncoding.EncodeToString([]byte(s))
}