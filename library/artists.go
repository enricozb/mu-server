package library

import "github.com/enricozb/mu-server/metadata"

func (l *Library) Artists() map[string][]metadata.Metadata {
	return l.artists
}
