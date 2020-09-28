// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

import (
	"sync"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	v5_store "github.com/mattermost/mattermost-server/v5/store"
	"github.com/wiggin77/merror"
)

type ServerIface interface {
	Config() *model.Config
	IsLeader() bool
	AddClusterLeaderChangedListener(listener func()) string
	RemoveClusterLeaderChangedListener(id string)
	GetStore() v5_store.Store
	GetLogger() *mlog.Logger
}

type SyncService struct {
	server           ServerIface
	leaderListenerId string

	store *store

	mux            sync.Mutex
	active         bool
	homeChannels   map[string]*homeChannel
	remoteChannels map[string]*remoteChannel
}

// NewService creates a service that synchronizes posts with another cluster (or stand-alone server).
func NewService(server ServerIface) (*SyncService, error) {
	service := &SyncService{
		server: server,
		store:  newStore(server.GetStore().Post(), server.GetStore().Channel()),
	}
	service.leaderListenerId = server.AddClusterLeaderChangedListener(service.onClusterLeaderChange)

	if err := service.loadChannels(); err != nil {
		return nil, err
	}

	service.onClusterLeaderChange()

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

	ss.server.GetLogger().Debug("Shared Channels Sync server active")
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

	ss.server.GetLogger().Debug("Shared Channels Sync server inactive")
}

// loadChannels loads all home and remote channels into the service.
func (ss *SyncService) loadChannels() error {
	merr := merror.New()

	list, err := ss.store.cstore.GetSharedChannels()
	if err != nil {
		return err
	}
	homeChannels := make(map[string]*homeChannel)
	remoteChannels := make(map[string]*remoteChannel)

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
			homeChannels[channel.Id] = hc
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
			remoteChannels[channel.Id] = rc
		}
	}

	ss.homeChannels = homeChannels
	ss.remoteChannels = remoteChannels

	return merr.ErrorOrNil()
}

// OnEvent is called when a post or reaction is added to the local database.
// This triggers a sync to any remote connections.
func (ss *SyncService) OnEvent(event *model.WebSocketEvent) {
	if !shouldHandleEventType(event) {
		return
	}

	ss.server.GetLogger().Debug("Shared channels event received.", mlog.String("type", event.EventType()))
}

func shouldHandleEventType(event *model.WebSocketEvent) bool {
	switch event.EventType() {
	case model.WEBSOCKET_EVENT_POSTED,
		model.WEBSOCKET_EVENT_POST_EDITED,
		model.WEBSOCKET_EVENT_POST_DELETED,
		model.WEBSOCKET_EVENT_REACTION_ADDED,
		model.WEBSOCKET_EVENT_REACTION_REMOVED:
		return true
	}
	return false
}
