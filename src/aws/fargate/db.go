package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

type Item struct {
	UUID      string `json:"uuid"`
	Status    string `json:"status"`
	Sentiment string `json:"sentiment"`
}

func GetResult(uuid string) (Item, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(cfg.DynamoDbEndpoint),
		Region:   aws.String(cfg.Region),
	})

	svc := dynamodb.New(sess)
	item := Item{}

	// Query
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(cfg.DynamoDbTable),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},
	})

	// Abort on any error
	if err != nil {
		log.Error(err.Error())
		return item, err
	}

	// Unmarshal and return
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Errorf("Failed to unmarshal: %v", err)
	}
	return item, nil
}

func PutItem(uuid string) (*dynamodb.PutItemOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(cfg.DynamoDbEndpoint),
		Region:   aws.String(cfg.Region),
	})

	svc := dynamodb.New(sess)

	// Marshal map in the AWS accepted format
	item, err := dynamodbattribute.MarshalMap(Item{
		UUID:   uuid,
		Status: "processing",
	})
	if err != nil {
		log.Errorf("Error while marshaling map for %v: %v", uuid, err)
		return nil, err
	}

	// Query
	result, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(cfg.DynamoDbTable),
		Item:      item,
	})
	if err != nil {
		log.Errorf("Error while putting item %v in DynamoDB: %v", uuid, err)
	}

	return result, err
}
