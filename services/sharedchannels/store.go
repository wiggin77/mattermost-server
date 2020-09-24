// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import "github.com/mattermost/mattermost-server/v5/model"

// ServerStore is a subset of `store.Store`
type ServerStore interface {
	Channel() ChannelStore
	Post() PostStore
}

// ChannelStore is a subset of `store.ChannelStore`
type ChannelStore interface {
	// GetSharedChannels fetches all shared channels across all teams.
	GetSharedChannels() (*model.SharedChannelList, error)
}

// PostStore is a subset of `store.PostStore`
type PostStore interface {
	GetPosts(options model.GetPostsOptions, allowFromCache bool) (*model.PostList, *model.AppError)
}

// store provides storage for shared channel service.
type store struct {
	pstore PostStore
	cstore ChannelStore
}

func newStore(p PostStore, c ChannelStore) *store {
	return &store{
		pstore: p,
		cstore: c,
	}
}
