package main

type AssetResponse struct {
	URL string `json:"Uri"`
}

type AnnouncementAsset struct {
	URL       string            `json:"url"`
	Checksum  string            `json:"checksum"`
	Signature string            `json:"signature"`
	Tags      map[string]string `json:"tags"`
}

type Announcement struct {
	Version string              `json:"version"`
	Tags    map[string]string   `json:"tags"`
	Active  bool                `json:"active"`
	Assets  []AnnouncementAsset `json:"assets"`
}
