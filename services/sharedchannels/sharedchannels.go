// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import (
	"errors"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	v5_store "github.com/mattermost/mattermost-server/v5/store"
)

type ServerIface interface {
	Config() *model.Config
	IsLeader() bool
	AddClusterLeaderChangedListener(listener func()) string
	RemoveClusterLeaderChangedListener(id string)
	GetStore() v5_store.Store
}

type SyncService struct {
	server           ServerIface
	leaderListenerId string

	store *store

	mux          sync.Mutex
	active       bool
	clusterGroup []clusterGroup
}

// NewService creates a service that synchronizes posts with another cluster (or stand-alone server).
func NewService(server ServerIface, pstore PostStore) (*SyncService, error) {
	service := &SyncService{
		server: server,
		store:  newStore(server.GetStore().Post(), server.GetStore().Channel()),
	}
	service.leaderListenerId = server.AddClusterLeaderChangedListener(service.onClusterLeaderChange)

	var err error
	if service.clusterGroup, err = service.buildClusterGroup(); err != nil {
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

// buildClusterGroup creates a clusterGroup array containing all
// the remote clusters this server shares a channel with regardless of
// where the channel is homed.
func (ss *SyncService) buildClusterGroup() ([]clusterGroup, error) {
	return nil, errors.New("not implemented yet")
}
