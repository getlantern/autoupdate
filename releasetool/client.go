package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	cfg Config
	cli *http.Client
}

func NewClient(cfg Config) *Client {
	c := Client{
		cfg: cfg,
		cli: &http.Client{},
	}
	return &c
}

func (c *Client) Header() http.Header {
	h := http.Header{}
	h.Add("Authorization", "Basic "+c.cfg.AuthHeader())
	h.Add("Content-Type", "application/octect-stream")
	return h
}

func (c *Client) PostRelease(src string) (*AssetResponse, error) {
	var err error
	var fp *os.File
	var buf *bytes.Buffer
	var res *http.Response
	var uri *url.URL

	if uri, err = url.Parse(fmt.Sprintf(endpointAssets, c.cfg.ApplicationID)); err != nil {
		return nil, err
	}

	buf = bytes.NewBuffer(nil)

	// Opening file.
	if fp, err = os.Open(src); err != nil {
		return nil, err
	}

	defer fp.Close()

	io.Copy(buf, fp)

	req := &http.Request{
		URL:    uri,
		Method: "POST",
		Header: c.Header(),
		Body:   ioutil.NopCloser(buf),
	}

	req.ContentLength = int64(buf.Len())

	log.Println("Uploading binary file...")

	if res, err = c.cli.Do(req); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		// Reading body.
		switch res.Header.Get("Content-Type") {
		case "application/json; charset=UTF-8":
			var msgData []byte
			var msg AssetResponse

			if msgData, err = ioutil.ReadAll(res.Body); err != nil {
				return nil, err
			}

			if err = json.Unmarshal(msgData, &msg); err != nil {
				return nil, err
			}

			return &msg, nil
		default:
			return nil, fmt.Errorf("Expecting application/json response, got %s.", res.Header.Get("Content-Type"))
		}
	}

	return nil, errors.New("Failed to create asset.")
}
