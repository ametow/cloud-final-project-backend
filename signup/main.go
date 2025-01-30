package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const tableName = "final-project-users"

var dbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, request Request) (Response, error) {
	// Check if email exists in DynamoDB
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: request.Email},
		},
	}
	result, err := dbClient.GetItem(ctx, getItemInput)
	if err != nil {
		return Response{Success: false, Message: fmt.Sprintf("Error checking email: %v", err)}, nil
	}

	if result.Item != nil {
		return Response{Success: false, Message: "Email already exists"}, nil
	}

	// Insert new user if email is unique
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"email":    &types.AttributeValueMemberS{Value: request.Email},
			"password": &types.AttributeValueMemberS{Value: request.Password},
			"name":     &types.AttributeValueMemberS{Value: request.Name},
		},
	}
	_, err = dbClient.PutItem(ctx, putItemInput)
	if err != nil {
		return Response{Success: false, Message: fmt.Sprintf("Error inserting user: %v", err)}, nil
	}

	return Response{Success: true, Message: "User registered successfully"}, nil
}

func main() {
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var req Request
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"success": false, "message": "Invalid request body"}`}, nil
		}

		resp, _ := handler(ctx, req)
		respBody, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(respBody)}, nil
	})
}
