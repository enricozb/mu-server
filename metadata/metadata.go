package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

type Metadata struct {
	ID       string `json:"id"`
	Album    string `json:"album"`
	Artist   string `json:"artist"`
	Duration string `json:"duration"`
	Title    string `json:"title"`

	// optional metadata

	Date  string `json:"year,omitempty"`
	Track string `json:"track,omitempty"`
}

// anyString is a type that unmarshals json strings and numbers to a golang string.
type anyString string

func (s *anyString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case string:
		*s = anyString(v)
	case int:
		*s = anyString(strconv.Itoa(v))
	case float64:
		*s = anyString(fmt.Sprintf("%f", v))
	default:
		return fmt.Errorf("unexpected type: %T", v)
	}

	return nil
}

type metadata struct {
	// required metadata

	ID       anyString `json:"SourceFile"`
	Album    anyString `json:"Album"`
	Artist   anyString `json:"Artist"`
	Duration anyString `json:"Duration"`
	Title    anyString `json:"Title"`

	// optional metadata

	Date  anyString `json:"Year"`
	Track anyString `json:"Track"`

	// substitute fields

	// Product is used in place of Album if it's not present.
	Product     anyString `json:"Product"`
	TrackNumber anyString `json:"TrackNumber"`
}

func (m *metadata) validate() error {
	if m.Album == "" && m.Product != "" {
		m.Album = m.Product
	}

	if m.ID == "" || m.Title == "" || m.Artist == "" || m.Album == "" {
		return fmt.Errorf("missing fields: %s", m.ID)
	}

	return nil
}

func (m *metadata) export() Metadata {
	return Metadata{
		ID:       string(m.ID),
		Album:    string(m.Album),
		Artist:   string(m.Artist),
		Duration: string(m.Duration),
		Title:    string(m.Title),
		Date:     string(m.Date),
		Track:    string(m.Track),
	}
}

func Fetch(root string, files []string) ([]Metadata, error) {
	tmp, err := ioutil.TempFile("", "mu-server-files")
	if err != nil {
		return nil, fmt.Errorf("temp file :%v", err)
	}

	for _, file := range files {
		tmp.Write(append([]byte(file), '\n'))
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	out, err := exec.Command(
		"exiftool",
		"-@", tmp.Name(),
		"-json",

		"-Album",
		"-Artist",
		"-Duration",
		"-Product",
		"-Title",
		"-TrackNumber",
		"-Track",
		"-Year",
	).Output()

	if err != nil {
		return nil, fmt.Errorf("exec exiftool: %v", err)
	}

	var metadatas []metadata
	if err := json.Unmarshal(out, &metadatas); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	var exportedMetadatas []Metadata
	for _, metadata := range metadatas {
		if err := metadata.validate(); err != nil {
			return nil, fmt.Errorf("validate: %v", err)
		}

		exportedMetadatas = append(exportedMetadatas, metadata.export())
	}

	return exportedMetadatas, nil
}
