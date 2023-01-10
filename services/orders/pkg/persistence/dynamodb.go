package persistence

import (
	"context"
	"fmt"
	"time"

	"eshop-orders/pkg/awsconfig"
	"eshop-orders/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type dynamoDB struct {
	av     *dynamodb.Client
	config DynamoDBConfig
}

type DynamoDBConfig struct {
	Url       string `sbc-key:"url"`
	Region    string `sbc-key:"region"`
	TableName string `sbc-key:"tableName"`
}

func createClient(ac awsconfig.Config, c DynamoDBConfig) (*dynamodb.Client, error) {
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
				SessionToken:    ac.SessionToken,
				Source:          ac.Source,
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func NewDynamoDB(ac awsconfig.Config, c DynamoDBConfig) (Repository, error) {
	av, err := createClient(ac, c)

	return &dynamoDB{av: av, config: c}, err
}

func (d *dynamoDB) Create(ctx context.Context, o models.Order) (*models.Order, error) {
	o.ID = uuid.New().String()
	o.Date = time.Now()
	i, err := attributevalue.MarshalMap(o)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}

	if _, err := d.av.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.config.TableName,
		Item:      i,
	}); err != nil {
		return nil, fmt.Errorf("error putting item: %w", err)
	}

	return &o, nil
}

func (d *dynamoDB) Read(ctx context.Context, id string) (*models.Order, error) {
	mid, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, fmt.Errorf("error marshaling id: %w", err)
	}

	oi, err := d.av.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.config.TableName,
		Key:       map[string]types.AttributeValue{"ID": mid},
	})
	if err != nil {
		return nil, fmt.Errorf("error reading orders with id '%s': %w", id, err)
	}
	var o models.Order
	if err := attributevalue.UnmarshalMap(oi.Item, &o); err != nil {
		return nil, fmt.Errorf("error unmarshaling order: %w", err)
	}

	return &o, nil
}

func (d *dynamoDB) List(ctx context.Context) ([]models.Order, error) {
	ooi, err := d.av.Scan(ctx, &dynamodb.ScanInput{
		TableName: &d.config.TableName,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}
	var oo []models.Order
	if err := attributevalue.UnmarshalListOfMaps(ooi.Items, &oo); err != nil {
		return nil, fmt.Errorf("error unmarshaling orders: %w", err)
	}

	return oo, nil
}
