package library

import "github.com/enricozb/mu-server/metadata"

func (l *Library) Albums() map[string][]metadata.Metadata {
	return l.albums
}
