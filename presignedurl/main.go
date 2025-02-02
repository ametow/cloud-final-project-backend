package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Request struct {
	Email    string `json:"email"`
	Filename string `json:"filename"`
}

type Response struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	UploadURL string `json:"upload_url"`
}

const tableName = "ProjectUsers"
const bucketName = "arsinuxbucket"
const expiration = 5 * time.Minute // URL valid for 15 minutes

var dbClient *dynamodb.Client
var s3Client *s3.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)
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

	presigner := s3.NewPresignClient(s3Client, s3.WithPresignExpires(expiration))
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String("project/" + request.Filename),
		ContentType: aws.String("application/jpg"),
	})
	if err != nil {
		return Response{Success: false, Message: "Server error getting presigned url"}, nil
	}

	return Response{Success: true, Message: "Successfully generated url", UploadURL: req.URL}, nil
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
