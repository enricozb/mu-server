package library

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/enricozb/mu-server/metadata"
)

type Library struct {
	dir string
	fs  fs.FS

	albums  map[string][]metadata.Metadata
	artists map[string][]metadata.Metadata
	songs   map[string]metadata.Metadata
}

func New(dir string) *Library {
	return &Library{
		dir: dir,
		fs:  os.DirFS(dir),

		albums:  map[string][]metadata.Metadata{},
		artists: map[string][]metadata.Metadata{},
		songs:   map[string]metadata.Metadata{},
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

	metadata, err := metadata.Fetch(l.dir, files)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	for _, m := range metadata {
		l.albums[m.Album] = append(l.albums[m.Album], m)
		l.artists[m.Artist] = append(l.artists[m.Artist], m)
		l.songs[m.ID] = m
	}

	for album := range l.albums {
		sort.Slice(l.albums[album], func(i, j int) bool {
			return strings.TrimLeft(l.albums[album][i].Position, "0") < strings.TrimLeft(l.albums[album][j].Position, "0")
		})
	}

	return nil
}

func (l *Library) Size() int {
	return len(l.songs)
}

func (l *Library) Abs(song string) (string, error) {
	if filepath.IsAbs(song) {
		return "", fmt.Errorf("is abs: %s", song)
	}

	return filepath.Join(l.dir, song), nil
}

var supportedMimetypes = map[string]struct{}{
	"audio/flac":  {},
	"audio/m4a":   {},
	"audio/mp3":   {},
	"audio/mpeg":  {},
	"audio/wav":   {},
	"audio/x-m4a": {},
}
