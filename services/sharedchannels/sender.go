// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

// sender maintains "connections" with remote servers and sends posts.
type sender struct {
	quit chan struct{} // signals loop to exit
	done chan struct{} // closed when looped has exited
}

func newSender() (*sender, error) {
	s := &sender{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	return s, nil
}

func (s *sender) shutdown() {

}

func (s *sender) loop() {

}
