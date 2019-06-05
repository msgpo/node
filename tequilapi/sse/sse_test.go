/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package sse

import (
	"testing"

	"github.com/stretchr/testify/assert"

	nodeEvent "github.com/mysteriumnetwork/node/core/node/event"
	"github.com/mysteriumnetwork/node/core/service"
	natEvent "github.com/mysteriumnetwork/node/nat/event"
)

func TestHandler_ConsumeNATEvent(t *testing.T) {
	h := NewHandler()
	me := natEvent.Event{
		Stage: "somestage",
	}
	h.ConsumeNATEvent(me)

	assert.Equal(t, `{"payload":{"stage":"somestage","successful":false},"type":"nat"}`, <-h.messages)
}

func TestHandler_ConsumeServiceStateEvent(t *testing.T) {
	h := NewHandler()

	me := service.EventPayload{
		Status: "somestage",
	}
	h.ConsumeServiceStateEvent(me)

	assert.Equal(t, `{"payload":{"id":"","providerId":"","type":"","status":"somestage"},"type":"service-status"}`, <-h.messages)
}

func TestHandler_Stops(t *testing.T) {
	h := NewHandler()

	wait := make(chan struct{})
	go func() {
		h.serve()
		wait <- struct{}{}
	}()

	h.stop()
	<-wait
}

func TestHandler_ConsumeNodeEvent_Stops(t *testing.T) {
	h := NewHandler()
	me := nodeEvent.Payload{
		Status: nodeEvent.StatusStopped,
	}
	h.ConsumeNodeEvent(me)
	h.serve()
}

func TestHandler_ConsumeNodeEvent_Starts(t *testing.T) {
	h := NewHandler()
	me := nodeEvent.Payload{
		Status: nodeEvent.StatusStarted,
	}

	h.ConsumeNodeEvent(me)

	// without starting, this would block forever
	h.newClients <- make(chan string)
	h.newClients <- make(chan string)

	h.stop()
}