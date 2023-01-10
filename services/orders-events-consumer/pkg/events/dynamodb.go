package events

import (
	"context"
	"encoding/json"
	"eshop-events-consumer/pkg/awsconfig"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

type Repository interface {
	ReadEvents(ctx context.Context) (<-chan []byte, <-chan error, error)
}

type dynamoStreams struct {
	av     *dynamodbstreams.Client
	config DynamoStreamsConfig
}

type DynamoStreamsConfig struct {
	Url        string `sbc-key:"streamsUrl"`
	Region     string `sbc-key:"region"`
	TableName  string `sbc-key:"tableName"`
	StreamsArn string `sbc-key:"streamArn"`
}

func createClient(ctx context.Context, ac awsconfig.Config, c DynamoStreamsConfig) (*dynamodbstreams.Client, error) {
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
				SessionToken:    ac.SessionToken,
				Source:          ac.Source,
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return dynamodbstreams.NewFromConfig(cfg), nil
}

func NewDynamoStreams(ctx context.Context, ac awsconfig.Config, c DynamoStreamsConfig) (Repository, error) {
	av, err := createClient(ctx, ac, c)

	return &dynamoStreams{av: av, config: c}, err
}

func (d *dynamoStreams) ReadEvents(ctx context.Context) (<-chan []byte, <-chan error, error) {
	sdd, err := d.getStreamDescriptions(ctx)
	if err != nil {
		return nil, nil, err
	}

	sdsii := map[*types.StreamDescription][]string{}
	for _, sd := range sdd {
		lsii, err := d.getShardsId(ctx, sd)
		if err != nil {
			return nil, nil, err
		}
		sdsii[sd] = lsii
	}

	cr, ce := make(chan []byte, 1), make(chan error, 1)
	for sd, sii := range sdsii {
		for _, si := range sii {
			it, err := d.getShardIterator(ctx, sd, &si)
			if err != nil {
				return nil, nil, err
			}
			go func() {
				defer func() {
					close(cr)
					close(ce)
				}()

				d.receiveMessages(ctx, sd, it, cr, ce)
			}()
		}
	}
	return cr, ce, nil
}

func (d *dynamoStreams) getStreamARN(ctx context.Context) (*string, error) {
	if d.config.StreamsArn != "" {
		return &d.config.StreamsArn, nil
	}

	sdd, err := d.getStreamDescriptions(ctx)
	if err != nil {
		return nil, err
	}

	return sdd[0].StreamArn, nil
}

func (d *dynamoStreams) getStreamDescription(ctx context.Context) (*types.StreamDescription, error) {
	arn, _ := d.getStreamARN(ctx)

	ds, err := d.av.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{
		StreamArn: arn,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving description for stream (Arn: %s, Table: %s): %w",
			*arn, *&d.config.TableName, err)
	}

	return ds.StreamDescription, nil
}

func (d *dynamoStreams) getStreamDescriptions(ctx context.Context) ([]*types.StreamDescription, error) {
	lso, err := d.av.ListStreams(ctx, &dynamodbstreams.ListStreamsInput{
		TableName: &d.config.TableName,
	})
	if err != nil {
		return nil, err
	}

	if len(lso.Streams) == 0 {
		return nil, fmt.Errorf("no streams found")
	}
	jsonStreams, _ := json.Marshal(lso.Streams)
	log.Printf("found %d streams: %s", len(lso.Streams), jsonStreams)

	sdd := make([]*types.StreamDescription, len(lso.Streams))
	for i, sd := range lso.Streams {
		arn := sd.StreamArn

		if err != nil {
			return nil, fmt.Errorf("could not get StreamArn: %w", err)
		}

		ds, err := d.av.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{
			StreamArn: arn,
		})
		if err != nil {
			return nil, fmt.Errorf(
				"error retrieving description for stream (Arn: %s, Table: %s): %w",
				*arn, *&d.config.TableName, err)
		}

		sdd[i] = ds.StreamDescription
	}
	return sdd, nil
}

func (d *dynamoStreams) getShardsId(ctx context.Context, so *types.StreamDescription) ([]string, error) {
	ds, err := d.av.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{
		StreamArn: so.StreamArn,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving description for stream (Arn: %s, Label: %s, Table: %s): %w",
			*so.StreamArn, *so.StreamLabel, *so.TableName, err)
	}

	sdi := []string{}
	for _, s := range ds.StreamDescription.Shards {
		if s.SequenceNumberRange.EndingSequenceNumber != nil ||
			s.SequenceNumberRange.StartingSequenceNumber != s.SequenceNumberRange.EndingSequenceNumber {
			sdi = append(sdi, *s.ShardId)
		}
	}
	return sdi, nil
}

func (d *dynamoStreams) getShardIterator(ctx context.Context, so *types.StreamDescription, shardId *string) (*string, error) {
	si, err := d.av.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
		// TODO: a mechanism to be consistent across application crashes should be implemented
		ShardIteratorType: types.ShardIteratorTypeLatest,
		StreamArn:         so.StreamArn,
		ShardId:           shardId,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving ShardIterator for stream (Arn: %s, Label: %s, Table: %s): %s", *so.StreamArn, *so.StreamLabel, *so.TableName, err)
	}

	return si.ShardIterator, nil
}

func (d *dynamoStreams) receiveMessages(ctx context.Context, so *types.StreamDescription, si *string, ci chan []byte, cerr chan error) {
	shi := si
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("processing shard %s", *shi)
			nshi, err := d.processShard(ctx, shi, ci, cerr)
			if err != nil {
				cerr <- fmt.Errorf(
					"error retrieving records from stream (Arn: %s, Label: %s, Table: %s): %w",
					*so.StreamArn, *so.StreamLabel, *so.TableName, err)
			}

			if nshi != nil && *nshi != *shi {
				shi = nshi
			}
			time.Sleep(60 * time.Second)
		}
	}
}

func (d *dynamoStreams) processShard(ctx context.Context, si *string, ci chan []byte, cerr chan error) (*string, error) {
	rr, err := d.av.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{ShardIterator: si})
	if err != nil {
		return nil, err
	}

	log.Printf("found %d records for shard '%s'", len(rr.Records), *si)
	if len(rr.Records) == 0 {
		return nil, nil
	}

	for _, r := range rr.Records {
		b, err := json.Marshal(r)
		if err != nil {
			cerr <- fmt.Errorf("error marshaling event '%+v': %w", r, err)
			continue
		}
		log.Printf("send message on channel: %s", string(b))
		ci <- b

		// smo, err := cq.SendMessage(ctx, &sqs.SendMessageInput{
		// 	MessageBody:            aws.String(string(b)),
		// 	QueueUrl:               qu,
		// 	MessageDeduplicationId: generateMessageId(10),
		// 	MessageGroupId:         aws.String(MessageGroupId),
		// })
		// if err != nil {
		// 	log.Printf("error sending message on SQS: '%s': %s", string(b), err)
		// 	continue
		// }

		// log.Printf("message sent on SQS: %s", *smo.MessageId)
	}

	return rr.NextShardIterator, nil
}
