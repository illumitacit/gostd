package workercommons

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/azuresb"
	"gocloud.dev/pubsub/rabbitpubsub"
	"google.golang.org/protobuf/proto"
)

type PubClient struct {
	topic *pubsub.Topic

	ctx        context.Context
	logger     *zap.SugaredLogger
	sender     *azservicebus.Sender
	rabbitConn *amqp.Connection
}

// NewPubClient returns an initialized publisher client for the configured broker from the given application config.
func NewPubClient(logger *zap.SugaredLogger, broker *Broker, ctx context.Context) (*PubClient, error) {
	switch broker.Engine {
	case "azuresb":
		return newAzureSBSenderClient(logger, broker, ctx)
	case "rabbitmq":
		return newRabbitMQPublisherClient(logger, broker, ctx)
	}
	return nil, fmt.Errorf("Unknown engine")
}

// Close will close all the associated connections of the given publisher client.
func (clt *PubClient) Close() error {
	if clt == nil {
		return nil
	}

	if err := clt.topic.Shutdown(clt.ctx); err != nil {
		clt.logger.Errorf("Error shutting down publisher: %s", err)
		return err
	}

	if clt.sender != nil {
		if err := clt.sender.Close(clt.ctx); err != nil {
			clt.logger.Errorf("Error closing Azure PubSub sender: %s", err)
			return err
		}
	}

	if clt.rabbitConn != nil {
		if err := clt.rabbitConn.Close(); err != nil {
			clt.logger.Errorf("Error closing RabbitMQ connection: %s", err)
			return err
		}
	}

	return nil
}

// SendTask will send a protobuf encoded message representing a worker task across the open pubsub topic.
func (clt *PubClient) SendTask(task proto.Message) error {
	taskMsg, err := proto.Marshal(task)
	if err != nil {
		return err
	}
	return clt.topic.Send(clt.ctx, &pubsub.Message{
		Body: taskMsg,
		Metadata: map[string]string{
			// TODO: use enum
			"type": "Task",
		},
	})
}

func newAzureSBSenderClient(
	logger *zap.SugaredLogger, broker *Broker, ctx context.Context) (*PubClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	clt, err := azservicebus.NewClient(broker.ConnectionString, cred, nil)
	if err != nil {
		return nil, err
	}
	sender, err := azuresb.NewSender(clt, broker.TopicName, nil)
	if err != nil {
		return nil, err
	}
	topic, err := azuresb.OpenTopic(ctx, sender, nil)
	if err != nil {
		return nil, err
	}
	return &PubClient{
		topic:  topic,
		logger: logger,
		ctx:    ctx,
		sender: sender,
	}, nil
}

func newRabbitMQPublisherClient(
	logger *zap.SugaredLogger, broker *Broker, ctx context.Context,
) (*PubClient, error) {
	rabbitConn, err := amqp.Dial(fmt.Sprintf("amqp://%s/", broker.ConnectionString))
	if err != nil {
		return nil, err
	}

	topic := rabbitpubsub.OpenTopic(rabbitConn, broker.TopicName, nil)
	return &PubClient{
		topic:      topic,
		logger:     logger,
		ctx:        ctx,
		rabbitConn: rabbitConn,
	}, nil
}
