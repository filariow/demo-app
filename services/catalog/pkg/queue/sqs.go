package queue

import (
	"context"
	"encoding/json"
	"eshop-catalog/pkg/awsconfig"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSManager struct {
	cli       *sqs.Client
	QueueName string
	QueueUrl  string
}

func NewSQSManager(ctx context.Context, ac awsconfig.Config, config SQSConfig) (*SQSManager, error) {
	c, err := setupSQS(ctx, ac, &config)
	if err != nil {
		return nil, err
	}

	return &SQSManager{
		cli:       c,
		QueueName: config.QueueName,
		QueueUrl:  config.Url,
	}, nil
}

func (s SQSManager) IncomingMessages(ctx context.Context) (<-chan OrderCreatedSQSMessage, <-chan error, context.CancelFunc) {
	cmsg := make(chan OrderCreatedSQSMessage)
	cerr := make(chan error)
	ctx2, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx2.Done():
				return
			case <-time.After(5 * time.Second):
				msgs, err := s.receiveMessages(ctx2)
				if err != nil {
					cerr <- err
					continue
				}

				for _, msg := range msgs {
					cmsg <- msg
				}
			}
		}
	}()

	return cmsg, cerr, cancel
}

func (s SQSManager) receiveMessages(ctx context.Context) ([]OrderCreatedSQSMessage, error) {
	for {
		rmo, err := s.cli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: &s.QueueUrl,
		})
		if err != nil {
			return nil, err
		}

		msgs := []OrderCreatedSQSMessage{}
		for _, m := range rmo.Messages {
			var ocm OrderCreatedSQSBody
			if err := json.Unmarshal([]byte(*m.Body), &ocm); err != nil {
				return nil, fmt.Errorf("error unmarshalling message: %s: %w", *m.Body, err)
			}

			log.Printf("processing message: %v", ocm)
			for _, r := range ocm.Records {
				log.Printf("processing record: %v", r)
				msgs = append(msgs, OrderCreatedSQSMessage{Value: r, receiptHandle: m.ReceiptHandle})
			}
		}
		return msgs, nil
	}
}

func (s SQSManager) CompleteMessage(ctx context.Context, msg OrderCreatedSQSMessage) error {
	if _, err := s.cli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.QueueUrl,
		ReceiptHandle: msg.receiptHandle,
	}); err != nil {
		return fmt.Errorf("could not delete message from queue: %w", err)
	}
	return nil
}

func setupSQS(ctx context.Context, ac awsconfig.Config, c *SQSConfig) (*sqs.Client, error) {
	cq, err := createSQSClient(ctx, ac, *c)
	if err != nil {
		return nil, err
	}

	log.Printf("creating queue")
	_, err = cq.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(c.QueueName),
	})
	if err == nil {
		log.Printf("can not create queue '%s', checking if yet existing: %s", c.QueueName, err)
	}

	qui, err := cq.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(c.QueueName),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting queue '%s': %s", c.QueueName, err)
	}
	c.Url = *qui.QueueUrl
	return cq, nil
}

func createSQSClient(ctx context.Context, ac awsconfig.Config, c SQSConfig) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.Region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: c.Url, SigningRegion: c.Region}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     ac.AccessKeyID,
				SecretAccessKey: ac.SecretAccessKey,
				// Source:          ac.Source,
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return sqs.NewFromConfig(cfg), nil
}

// TODO: remove that
func createLocalSQSClient(ctx context.Context) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("elasticmq"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://sqs:9324"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "x",
				SecretAccessKey: "x",
				Source:          "Hard-coded credentials; values are irrelevant for local SQS",
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return sqs.NewFromConfig(cfg), nil
}
