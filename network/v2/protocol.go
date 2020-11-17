package networkv2

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// StreamNetProtocol used to transfer data from peer to peer
type StreamNetProtocol struct {
	node *Node
}

// sendMessage will send message to peers
func (procotol *StreamNetProtocol) sendMessage(ps *pubsub.PubSub, msg string) {
	msgId := make([]byte, 10)
	_, err := rand.Read(msgId)
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	if err != nil {
		return
	}
	now := time.Now().Unix()
	req := &Request{
		Type: Request_SEND_MESSAGE.Enum(),
		SendMessage: &SendMessage{
			Id:      msgId,
			Data:    []byte(msg),
			Created: &now,
		},
	}
	msgBytes, err := req.Marshal()
	if err != nil {
		return
	}
	err = ps.Publish(pubsubTopic, msgBytes)
}

// updatePeer means update peer, not use
func (procotol *StreamNetProtocol) updatePeer(ps *pubsub.PubSub, id peer.ID, handle string) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = handle

	req := &Request{
		Type: Request_UPDATE_PEER.Enum(),
		UpdatePeer: &UpdatePeer{
			UserHandle: []byte(handle),
		},
	}
	reqBytes, err := req.Marshal()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	err = ps.Publish(pubsubTopic, reqBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s -> %s\n", oldHandle, handle)
}

// chatInputLoop means how to chat with other peers
func (procotol *StreamNetProtocol) chatInputLoop(ctx context.Context, h host.Host, ps *pubsub.PubSub, donec chan struct{}) {
	for {
		var msgB []byte
		select {
		case msgB = <-procotol.node.SendMessageChan:
		}
		msg := string(msgB)
		if strings.HasPrefix(msg, "/name ") {
			newHandle := strings.TrimPrefix(msg, "/name ")
			newHandle = strings.TrimSpace(newHandle)
			procotol.updatePeer(ps, h.ID(), newHandle)
		} else {
			procotol.sendMessage(ps, msg)
		}
	}
}
