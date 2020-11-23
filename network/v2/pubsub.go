package networkv2

import (
	"context"
	"fmt"
	"os"

	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var handles = map[string]string{}

const pubsubTopic = "/libp2p/1.0.0"

// PubSubHandler ...
type PubSubHandler struct {
	node *Node
}

func (pubsub *PubSubHandler) pubsubMessageHandler(id peer.ID, msg *SendMessage) {
	handle, ok := handles[id.String()]
	if !ok {
		handle = id.ShortString()
	}
	fmt.Printf("%s: %s\n", handle, msg.Data)
	pubsub.node.Receive(msg.Data)
}

func (pubsub *PubSubHandler) pubsubUpdateHandler(id peer.ID, msg *UpdatePeer) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = string(msg.UserHandle)
	fmt.Printf("%s -> %s\n", oldHandle, msg.UserHandle)
}

func (pubsub *PubSubHandler) pubsubHandler(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		req := &Request{}
		err = req.Unmarshal(msg.Data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		switch *req.Type {
		case Request_SEND_MESSAGE:
			pubsub.pubsubMessageHandler(msg.GetFrom(), req.SendMessage)
		case Request_UPDATE_PEER:
			pubsub.pubsubUpdateHandler(msg.GetFrom(), req.UpdatePeer)
		}
	}
}
