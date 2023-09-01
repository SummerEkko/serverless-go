package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var dynamoClient *dynamodb.Client

func createItem(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)

	var requestBody Item
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	id := int(time.Now().Unix()) // Generate a unique ID
	item := Item{
		ID:   id,
		Name: requestBody.Name,
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Item: map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberN{Value: strconv.Itoa(item.ID)},
			"name": &types.AttributeValueMemberS{Value: item.Name},
		},
	}

	_, err = dynamoClient.PutItem(ctx, putItemInput)
	if err != nil {
		log.Printf("failed to put item to DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating item",
		}, nil
	}

	responseBody, _ := json.Marshal(map[string]int{"id": item.ID})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(createItem)
}
