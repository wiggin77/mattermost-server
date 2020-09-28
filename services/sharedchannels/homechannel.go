// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import (
	"errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

// homeChannel represents a shared channel owned by the local cluster.
type homeChannel struct {
	channel *model.SharedChannel
	remotes map[string]*remote
}

func newHomeChannel(sc *model.SharedChannel) *homeChannel {
	return &homeChannel{
		channel: sc,
		remotes: make(map[string]*remote),
	}
}

func (hc *homeChannel) addRemote(invite *model.SharedChannelInvite) error {
	return errors.New("not implemented yet")
}
