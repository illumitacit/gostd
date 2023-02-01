package workerstd

// Broker represents configuration options for the message queue broker used to enqueue tasks for the worker.
// This can be embedded in a viper compatible config struct.
type Broker struct {
	// Engine is the engine of the message broker. Must be one of azuresb or rabbitmq.
	Engine string `mapstructure:"engine"`

	// TopicName is the message queue topic where messages are published. This corresponds to the exchange when using
	// rabbitmq.
	TopicName string `mapstructure:"topic"`

	// ConnectionString is the connection string for connecting to the specific broker. For AzureSB, this is the
	// ServiceBus Namespace FQDN (NAMESPACE.servicebus.windows.net), while for RabbitMQ this is the server URL in the
	// format USERNAME:PASSWORD@HOST:PORT.
	ConnectionString string `mapstructure:"connstring"`

	// ServiceBusSubscriptionName is the Azure ServiceBus Topic Subscription that the worker should consume as. If
	// blank, assume that the topic is an Azure ServiceBus Queue instead of a Topic. This is only used with Azure
	// ServiceBus.
	ServiceBusSubscriptionName string `mapstructure:"subscription"`
}
