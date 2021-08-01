package library

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type String string

func (s *String) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case string:
		*s = String(v)
	case int:
		*s = String(strconv.Itoa(v))
	case float64:
		*s = String(fmt.Sprintf("%f", v))
	default:
		return fmt.Errorf("unexpected type: %T", v)
	}

	return nil
}

type MediaMetadata struct {
	// required metadata

	ID       String `json:"SourceFile"`
	Album    String `json:"Album"`
	Artist   String `json:"Artist"`
	Duration String `json:"Duration"`
	Title    String `json:"Title"`

	// optional metadata

	Date  String `json:"Year"`
	Track String `json:"Track"`

	// substitute fields

	// Product is used in place of Album if it's not present.
	Product     String `json:"Product"`
	TrackNumber String `json:"TrackNumber"`
}

func (m *MediaMetadata) validate() error {
	if m.Album == "" && m.Product != "" {
		m.Album = m.Product
	}

	if m.ID == "" || m.Title == "" || m.Artist == "" || m.Album == "" {
		return fmt.Errorf("missing fields: %s", m.ID)
	}
	return nil
}

type Library struct {
	metadata map[string]MediaMetadata
	dir      string
	fs       fs.FS
}

func New(dir string) *Library {
	return &Library{
		metadata: map[string]MediaMetadata{},
		dir:      dir,
		fs:       os.DirFS(dir),
	}
}

func (l *Library) Init() error {
	var files []string

	if err := fs.WalkDir(l.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		f, err := l.fs.Open(path)
		if err != nil {
			return fmt.Errorf("open: %v", err)
		}
		defer f.Close()

		mime, err := mimetype.DetectReader(f)
		if err != nil {
			return fmt.Errorf("detect reader: %v", err)
		}

		if _, supported := supportedMimetypes[mime.String()]; !supported {
			if strings.HasPrefix(mime.String(), "audio/") {
				fmt.Printf("unsupported audio format '%s': %s\n", mime, path)
			}
			return nil
		}

		files = append(files, filepath.Join(l.dir, path))

		return nil
	}); err != nil {
		return fmt.Errorf("walk: %v", err)
	}

	return l.initMetadata(files)
}

func (l *Library) initMetadata(files []string) error {
	tmp, err := ioutil.TempFile("", "mu-server-files")
	if err != nil {
		return fmt.Errorf("temp file :%v", err)
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
		return fmt.Errorf("exec exiftool: %v", err)
	}

	var metadatas []MediaMetadata
	if err := json.Unmarshal(out, &metadatas); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	for _, metadata := range metadatas {
		if err := metadata.validate(); err != nil {
			return fmt.Errorf("validate: %v", err)
		}
		l.metadata[string(metadata.ID)] = metadata
	}

	fmt.Printf("precomputed metadata for %d items\n", len(metadatas))

	return nil
}

var supportedMimetypes = map[string]struct{}{
	"audio/flac":  {},
	"audio/m4a":   {},
	"audio/mp3":   {},
	"audio/mpeg":  {},
	"audio/wav":   {},
	"audio/x-m4a": {},
}
