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
	ImageURL string `json:"image_url"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const tableName = "ProjectUsers"

var dbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, request Request) (Response, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: request.Email},
		},
	}
	result, err := dbClient.GetItem(ctx, getItemInput)
	if err != nil {
		return Response{Success: false, Message: fmt.Sprintf("Error looking up the user: %v", err)}, nil
	}

	if result.Item == nil {
		return Response{Success: false, Message: "User does not exist"}, nil
	}

	err = updateItem(ctx, request.Email, request.ImageURL)
	if err != nil {
		return Response{Success: false, Message: "Couldn't update image url"}, nil
	}

	return Response{Success: true, Message: "Successfully updated image url"}, nil
}
func updateItem(ctx context.Context, email, imageUrl string) error {
	_, err := dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
		UpdateExpression:          aws.String("SET image_url = :newVal"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":newVal": &types.AttributeValueMemberS{Value: imageUrl}},
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	return err
}

func main() {
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var req Request
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"success": false, "message": "Invalid request body"}`}, nil
		}

		resp, _ := handler(ctx, req)
		respBody, _ := json.Marshal(resp)
		code := 400
		if resp.Success {
			code = 200
		}
		return events.APIGatewayProxyResponse{
			StatusCode: code,
			Body:       string(respBody),
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
			},
		}, nil
	})
}
