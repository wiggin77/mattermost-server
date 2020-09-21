// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

// clusterGroup represents two or more clusters which share one or more channels with each other.
type clusterGroup struct {
	members map[string]*cluster
}

// cluster represents a cluster within a cluster group
// A cluster can be an actual cluster or stand-alone server.
type cluster struct {
}
