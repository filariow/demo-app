package queue

type SQSConfig struct {
	Url       string `sbc-key:"url"`
	Region    string `sbc-key:"region"`
	QueueName string `sbc-key:"queueName"`
}
