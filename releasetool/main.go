package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const (
	configFile = "equinox.yaml"
)

var (
	flagConfigFile = flag.String("config", configFile, "Configuration file.")
	flagSource     = flag.String("source", "", "Source binary file.")
	flagVersion    = flag.Int("version", -1, "Version number.")
)

func main() {
	var fp *os.File
	var cfg Config
	var err error
	var buf []byte

	// Parsing flags
	flag.Parse()

	// Opening config file.
	if fp, err = os.Open(*flagConfigFile); err != nil {
		log.Fatal(fmt.Errorf("Could not open config file: %q", err))
	}
	defer fp.Close()

	if buf, err = ioutil.ReadAll(fp); err != nil {
		log.Fatal(fmt.Errorf("Could not read config file: %q", err))
	}

	// Parsing YAML.
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		log.Fatal(fmt.Errorf("Could not parse config file: %q", err))
	}

	// Creating client
	var res *AssetResponse

	cli := NewClient(cfg)
	if res, err = cli.PostRelease(*flagSource); err != nil {
		log.Fatal(fmt.Errorf("Could not upload release: %q", err))
	}

	fmt.Printf("AssetResponse: %v\n", res)
}
