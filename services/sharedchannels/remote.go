// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import "errors"

// remote represents a remote cluster.
type remote struct {
	url  string
	send *sender
	recv *receiver
}

func newRemote(URL string, store *store) (*remote, error) {
	return nil, errors.New("not implemented yet")
}
