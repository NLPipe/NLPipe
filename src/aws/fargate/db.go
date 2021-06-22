package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type Item struct {
	UUID      string `json:"uuid"`
	Status    string `json:"status"`
	Sentiment string `json:"sentiment"`
}

func GetResult(uuid string) (Item, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://127.0.0.1:8000"),
		Region:   aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)
	item := Item{}

	// Query
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("nlpipe-results"),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},
	})

	// Abort on any error
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	// Unmarshal and return
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal, %v", err))
	}
	return item, nil
}

func PutItem(uuid string) (*dynamodb.PutItemOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://127.0.0.1:8000"),
		Region:   aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)

	// Marshal map in the AWS accepted format
	item, err := dynamodbattribute.MarshalMap(Item{
		UUID:   uuid,
		Status: "processing",
	})
	if err != nil {
		log.Fatalf("Error while marshaling map for %v %v", uuid, err)
		return nil, err
	}

	// Query
	result, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("nlpipe-results"),
		Item:      item,
	})
	if err != nil {
		log.Fatalf("Error while putting item %v in DynamoDB %v", uuid, err)
	}

	return result, err
}
