package queue

import (
	"context"
	"eshop-events-consumer/pkg/awsconfig"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSManager struct {
	cli       *sqs.Client
	awsConfig awsconfig.Config
	sqsConfig SQSConfig
}

func NewSQSManager(ctx context.Context, ac awsconfig.Config, config SQSConfig) (*SQSManager, error) {
	c, err := setupSQS(ctx, ac, &config)
	if err != nil {
		return nil, err
	}

	return &SQSManager{
		cli:       c,
		awsConfig: ac,
		sqsConfig: config,
	}, nil
}

func (s SQSManager) SendMessage(ctx context.Context, msg []byte) error {
	smo, err := s.cli.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(msg)),
		QueueUrl:    &s.sqsConfig.Url,
	})
	if err != nil {
		return fmt.Errorf("error sending message on SQS: '%s': %s", string(msg), err)
	}

	log.Printf("message sent on SQS: %s", *smo.MessageId)
	return nil
}

func setupSQS(ctx context.Context, ac awsconfig.Config, c *SQSConfig) (*sqs.Client, error) {
	cq, err := createSQSClient(ctx, ac, *c)
	if err != nil {
		return nil, err
	}

	log.Printf("creating queue")
	q, err := cq.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(c.QueueName),
	})
	if err == nil {
		c.Url = *q.QueueUrl
		return cq, nil
	}

	log.Printf("can not create queue '%s', checking if yet existing: %s", c.QueueName, err)
	gqu, err := cq.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(c.QueueName),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting queue '%s': %s", c.QueueName, err)
	}

	log.Printf("queue '%s' found at: '%s'", c.QueueName, *gqu.QueueUrl)
	c.Url = *gqu.QueueUrl
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
