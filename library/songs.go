package library

import "github.com/enricozb/mu-server/metadata"

func (l *Library) Songs() map[string]metadata.Metadata {
	return l.songs
}
