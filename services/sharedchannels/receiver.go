// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannels

// receiver has incoming messages routed via RestAPI or websocket.
type receiver struct {
	quit chan struct{} // signals loop to exit
	done chan struct{} // closed when looped has exited
}

func newReceiver() (*receiver, error) {
	r := &receiver{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	return r, nil
}

func (r *receiver) shutdown() {

}

func (r *receiver) loop() {

}
