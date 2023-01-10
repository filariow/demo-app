package main

import (
	"context"
	"errors"
	"eshop-orders/pkg/awsconfig"
	"eshop-orders/pkg/persistence"
	"eshop-orders/pkg/seed"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoDB struct {
	av     *dynamodb.Client
	config persistence.DynamoDBConfig
}

type Database interface {
	Init(context.Context) error
}

func createClient(ac awsconfig.Config, c persistence.DynamoDBConfig) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: c.Url, SigningRegion: c.Region}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     ac.AccessKeyID,
				SecretAccessKey: ac.SecretAccessKey,
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func newDynamoDB(ac awsconfig.Config, c persistence.DynamoDBConfig) (Database, error) {
	av, err := createClient(ac, c)

	return &dynamoDB{av: av, config: c}, err
}

func (d *dynamoDB) createTables(ctx context.Context) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName:   aws.String(d.config.TableName),
		BillingMode: types.BillingModePayPerRequest,
		StreamSpecification: &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: types.StreamViewTypeNewAndOldImages,
		},
	}

	if _, err := d.av.CreateTable(ctx, input); err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	waiter := dynamodb.NewTableExistsWaiter(d.av)
	if err := waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(d.config.TableName)}, 5*time.Minute); err != nil {
		return fmt.Errorf("error waiting for table creation: %w", err)
	}

	return nil
}

type OrderedProduct struct {
	ID           int64
	Name         string
	PhotoURL     string
	UnitsOrdered int64
}

type Order struct {
	ID              int64
	Date            time.Time
	OrderedProducts []OrderedProduct
}

func (d *dynamoDB) seedData(ctx context.Context) error {
	for _, o := range seed.Data {
		i, err := attributevalue.MarshalMap(o)
		if err != nil {
			return fmt.Errorf("error marshaling data: %w", err)
		}
		if _, err := d.av.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: &d.config.TableName,
			Item:      i,
		}); err != nil {
			return fmt.Errorf("error putting item: %w", err)
		}
	}

	return nil
}

func (d *dynamoDB) Init(ctx context.Context) error {
	created := false
	t, err := d.av.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &d.config.TableName,
	})
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if !errors.As(err, &nfe) {
			return fmt.Errorf("error retrieving table '%s': %w", d.config.TableName, err)
		}
		if err := d.createTables(ctx); err != nil {
			return err
		}
		created = true
	}

	if created || *t.Table.ItemCount == 0 {
		if err := d.seedData(ctx); err != nil {
			return err
		}
	}
	return nil
}
