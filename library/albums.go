package library

import (
	"fmt"
	"io"

	"github.com/enricozb/mu-server/metadata"
)

func (l *Library) Albums() map[string][]metadata.Metadata {
	return l.albums
}

func (l *Library) AlbumCover(id string) (io.Reader, error) {
	songs, ok := l.albums[id]
	if !ok {
		return nil, fmt.Errorf("album does not exist: %s", id)
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("album is empty: %s", id)
	}

	path, err := l.Abs(songs[0].ID)
	if err != nil {
		return nil, fmt.Errorf("abs: %v", err)
	}

	return metadata.Cover(path)
}
