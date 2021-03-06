package DCP

import (
	"encoding/json"
	"github.com/google/uuid"
	"testing"
	"time"
)

type Packet struct {
	Msg string `json:"msg"`
}

func TestChannelTransport_Broadcast(t *testing.T) {
	chT := ChannelTransport{
		DataCh:         make(chan []byte),
		ReachableNodes: make(map[chan []byte]struct{}),
	}

	node1Chan := make(chan []byte)

	chT.ReachableNodes[node1Chan] = struct{}{}

	packet := &Packet{
		Msg: "foobar",
	}

	b, _ := json.Marshal(packet)
	chT.Broadcast(uuid.New(), b, func() {
		return
	})

	received := <-node1Chan

	var packetReceived Packet
	_ = json.Unmarshal(received, &packetReceived)

	if packetReceived.Msg != "foobar" {
		t.Fail()
	}
}

func TestChannelTransport_Listener(t *testing.T) {
	chT := ChannelTransport{
		DataCh:         make(chan []byte),
		ReachableNodes: make(map[chan []byte]struct{}),
	}

	go chT.Listen(uuid.New(), func(i []byte) error {
		var packetReceived Packet
		_ = json.Unmarshal(i, &packetReceived)

		if packetReceived.Msg != "foobar" {
			t.Fail()
		}

		return nil
	})

	packet := &Packet{
		Msg: "foobar",
	}

	b, _ := json.Marshal(packet)

	chT.DataCh <- b
}

func TestChannelTransport_ListenerAndBroadcast(t *testing.T) {
	broadcaster := ChannelTransport{
		DataCh:         make(chan []byte),
		ReachableNodes: make(map[chan []byte]struct{}),
	}

	listener := ChannelTransport{
		DataCh:         make(chan []byte),
		ReachableNodes: make(map[chan []byte]struct{}),
	}

	broadcaster.ReachableNodes[listener.DataCh] = struct{}{}
	listener.ReachableNodes[broadcaster.DataCh] = struct{}{}

	go listener.Listen(uuid.New(), func(i []byte) error {
		return nil
	})

	packet := &Packet{
		Msg: "foobar",
	}

	b, _ := json.Marshal(packet)

	broadcaster.Broadcast(uuid.New(), b, func() {
		return
	})

	time.Sleep(1 * time.Millisecond)
}
