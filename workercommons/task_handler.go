package workercommons

import (
	"gocloud.dev/pubsub"
	"google.golang.org/protobuf/proto"
)

// TaskHandler is the interface that task handlers passed to the worker app should implement.
type TaskHandler[T proto.Message] interface {
	HandleTaskMsg(*T, *pubsub.Message) error
}
