// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import "github.com/mattermost/mattermost-server/v5/model"

// remoteChannel represents a shared channel owned by a remote cluster.
type remoteChannel struct {
	channel *model.SharedChannel
	remote  *remote
}
