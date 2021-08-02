package library

import (
	"fmt"
	"io"

	"github.com/enricozb/mu-server/metadata"
)

func (l *Library) Songs() map[string]metadata.Metadata {
	return l.songs
}

func (l *Library) SongCover(id string) (io.Reader, error) {
	path, err := l.Abs(id)
	if err != nil {
		return nil, fmt.Errorf("abs: %v", err)
	}

	return metadata.Cover(path)
}
