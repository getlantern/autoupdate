package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/inconshreveable/go-update"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	configFile = "equinox.yaml"
)

var (
	flagConfigFile = flag.String("config", configFile, "Configuration file.")
	flagSource     = flag.String("source", "", "Source binary file.")
	flagChannel    = flag.String("channel", "stable", "Release channel.")
	flagVersion    = flag.Int("version", -1, "Version number.")
	flagArch       = flag.String("arch", "", "Build architecture. (amd64|386|arm)")
	flagOS         = flag.String("os", "", "Operating system. (linux|windows|darwin)")
)

func main() {
	var fp *os.File
	var cfg Config
	var err error
	var buf []byte

	// Parsing flags
	flag.Parse()

	if *flagVersion < 0 {
		log.Fatal("Version must be a positive integer (-version).")
	}

	if *flagArch == "" {
		log.Fatal("Missing build architecture (-arch).")
	}

	if *flagOS == "" {
		log.Fatal("Missing build target operating system (-os).")
	}

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

	var checksum []byte

	if checksum, err = update.ChecksumForFile(*flagSource); err != nil {
		log.Fatal(fmt.Errorf("Could not create checksum for file: %q", err))
	}

	// Loading private key
	var pb []byte
	var fpk *os.File

	if fpk, err = os.Open(cfg.PrivateKey); err != nil {
		log.Fatal(fmt.Errorf("Could not open private key: %q", err))
	}
	defer fpk.Close()

	if pb, err = ioutil.ReadAll(fpk); err != nil {
		log.Fatal(fmt.Errorf("Could not read private key: %q", err))
	}

	// Decoding PEM key.
	pemBlock, _ := pem.Decode(pb)

	var privateKey *rsa.PrivateKey
	if privateKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err != nil {
		log.Fatal(fmt.Errorf("Could not parse private key: %q", err))
	}

	// Checking message signature.
	var signature []byte
	if signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, checksum); err != nil {
		log.Fatal(fmt.Errorf("Could not create signature for file: %q", err))
	}

	// Preparing message.
	announce := Announcement{
		Version: strconv.Itoa(*flagVersion),
		Tags: map[string]string{
			"channel": *flagChannel,
		},
		Active: true,
		Assets: []AnnouncementAsset{
			AnnouncementAsset{
				URL:       res.URL,
				Checksum:  hex.EncodeToString(checksum),
				Signature: hex.EncodeToString(signature),
				Tags: map[string]string{
					"arch": *flagArch,
					"os":   *flagOS,
				},
			},
		},
	}

	var msg []byte
	if msg, err = json.Marshal(announce); err != nil {
		panic(err)
	}

	fmt.Printf("AssetResponse: %v\n", res)
	fmt.Printf("AssetResponse: %v\n", string(msg))
}
