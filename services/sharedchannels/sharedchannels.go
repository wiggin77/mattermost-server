// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import (
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/wiggin77/merror"
)

type ServerIface interface {
	Config() *model.Config
	IsLeader() bool
	AddClusterLeaderChangedListener(listener func()) string
	RemoveClusterLeaderChangedListener(id string)
	GetStore() ServerStore
}

type SyncService struct {
	server           ServerIface
	leaderListenerId string

	store *store

	mux            sync.Mutex
	active         bool
	homeChannels   []*homeChannel
	remoteChannels []*remoteChannel
}

// NewService creates a service that synchronizes posts with another cluster (or stand-alone server).
func NewService(server ServerIface, pstore PostStore) (*SyncService, error) {
	service := &SyncService{
		server: server,
		store:  newStore(server.GetStore().Post(), server.GetStore().Channel()),
	}
	service.leaderListenerId = server.AddClusterLeaderChangedListener(service.onClusterLeaderChange)

	if err := service.loadChannels(); err != nil {
		return nil, err
	}

	return service, nil
}

// Shutdown stops the service and frees any resources.
func (ss *SyncService) Shutdown() {
	ss.server.RemoveClusterLeaderChangedListener(ss.leaderListenerId)
}

// onClusterLeaderChange is called whenever the cluster leader may have changed.
func (ss *SyncService) onClusterLeaderChange() {
	if ss.server.IsLeader() {
		ss.start()
	} else {
		ss.stop()
	}
}

// start activates synchronization of posts between clusters (or stand-alone server).
// If sync is already active this call has no effect.
func (ss *SyncService) start() {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	if ss.active {
		return // already started
	}
	ss.active = true

}

// stop deactivates synchronization of posts between clusters.
// If sync is not active this call has no effect.
func (ss *SyncService) stop() {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	if !ss.active {
		return // already stopped
	}
	ss.active = false
}

// loadChannels loads all home and remote channels into the service.
func (ss *SyncService) loadChannels() error {
	merr := merror.New()

	list, err := ss.store.cstore.GetSharedChannels()
	if err != nil {
		return err
	}
	homeChannels := make([]*homeChannel, 0)
	remoteChannels := make([]*remoteChannel, 0)

	ss.mux.Lock()
	defer ss.mux.Unlock()

	for _, channel := range *list {
		if channel.Home {
			r, err := newRemote(channel.URL, ss.store)
			if err != nil {
				merr.Append(err)
				continue
			}
			hc := &homeChannel{
				channel: channel,
				remotes: map[string]*remote{channel.URL: r},
			}
			homeChannels = append(homeChannels, hc)
		} else {
			r, err := newRemote(channel.URL, ss.store)
			if err != nil {
				merr.Append(err)
				continue
			}
			rc := &remoteChannel{
				channel: channel,
				remote:  r,
			}
			remoteChannels = append(remoteChannels, rc)
		}
	}

	ss.homeChannels = homeChannels
	ss.remoteChannels = remoteChannels

	return merr.ErrorOrNil()
}
