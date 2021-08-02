package metadata

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Metadata struct {
	ID       string `json:"id"`
	Album    string `json:"album"`
	Artist   string `json:"artist"`
	Duration string `json:"duration"`
	Title    string `json:"title"`

	// optional metadata

	Date     string `json:"year,omitempty"`
	Position string `json:"track,omitempty"`
}

// track is structured to match mediainfo's output
type track struct {
	Type string `json:"@type"`

	// required

	Album    string `json:"Album"`
	Artist   string `json:"Performer"`
	Duration string `json:"Duration"`
	Title    string `json:"Track"`

	// optional

	Date     string `json:"Recorded_Date"`
	Position string `json:"Track_Position"`
}

// metadata is structured to match mediainfo's output
type metadata struct {
	ID string `json:"@ref"`

	Tracks []track `json:"track"`

	GeneralIdx int
}

func (m *metadata) validate(root string) (err error) {
	var t track
	for i, track := range m.Tracks {
		if track.Type == "General" {
			t = track
			m.GeneralIdx = i
			break
		}
	}

	if t == (track{}) {
		return fmt.Errorf("no 'General' track: %s", m.ID)
	}

	if t.Title == "" || t.Artist == "" || t.Album == "" || t.Duration == "" {
		return fmt.Errorf("missing fields: %+v", t)
	}

	if m.ID, err = filepath.Rel(root, string(m.ID)); err != nil {
		return fmt.Errorf("rel: %v", err)
	}

	return nil
}

// export must be called after a call to `validate`.
func (m *metadata) export() Metadata {
	t := m.Tracks[m.GeneralIdx]
	return Metadata{
		ID:       m.ID,
		Album:    t.Album,
		Artist:   t.Artist,
		Duration: t.Duration,
		Title:    t.Title,
		Date:     t.Date,
		Position: t.Position,
	}
}

func Fetch(root string, files []string) ([]Metadata, error) {
	var g errgroup.Group
	var mu sync.Mutex
	ctx := context.Background()
	workers := runtime.GOMAXPROCS(0)
	sem := semaphore.NewWeighted(int64(workers))

	var metadatas []metadata
	fetch := func(f string) error {
		var media struct {
			Metadata metadata `json:"media"`
		}

		out, err := exec.Command("mediainfo", "--Output=JSON", f).Output()
		if err != nil {
			return fmt.Errorf("exec mediainfo: %v", err)
		}

		if err := json.Unmarshal(out, &media); err != nil {
			return fmt.Errorf("unmarshal: %v", err)
		}

		mu.Lock()
		metadatas = append(metadatas, media.Metadata)
		mu.Unlock()

		return nil
	}

	for _, f := range files {
		g.Go(func(f string) func() error {
			return func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return fmt.Errorf("acquire: %v", err)
				}
				defer sem.Release(1)

				if err := fetch(f); err != nil {
					return fmt.Errorf("fetch: %v", err)
				}

				return nil
			}
		}(f))
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("wait: %v", err)
	}

	var exported []Metadata
	for _, metadata := range metadatas {
		if err := metadata.validate(root); err != nil {
			return nil, fmt.Errorf("validate: %v", err)
		}

		exported = append(exported, metadata.export())
	}

	return exported, nil
}

func Cover(path string) (io.Reader, error) {
	// --Cover_Data from: https://sourceforge.net/p/mediainfo/discussion/297610/thread/aeb4222d/#c9a7
	out, err := exec.Command("mediainfo", "--Full", "--Cover_Data=base64", path).Output()
	if err != nil {
		return nil, fmt.Errorf("exec mediainfo: %v: %s", err, out)
	}

	for _, line := range bytes.Split(out, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("Cover_Data")) {
			parts := bytes.SplitN(line, []byte(":"), 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("expected 2 parts but got %d", len(parts))
			} else if parts[1][0] != ' ' {
				return nil, fmt.Errorf("expected byte '%b' but got '%b'", ' ', parts[1][0])
			}

			return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(parts[1][1:])), nil
		}
	}

	return nil, fmt.Errorf("no cover data")
}
