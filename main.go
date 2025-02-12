package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	version = "unknown"
)

type SFNEvent struct {
	Path  string
	Input json.RawMessage
}

type RecordEvent struct {
	Path        string          `json:"path"`
	Input       json.RawMessage `json:"input"`
	MyTaskToken string          `json:"myTaskToken"`
}

// Data represents the nested "data" object containing price and color
type Data struct {
	Price float64 `json:"price"`
	Color string  `json:"color"`
}

// Product represents the structure of the main product object
type ResponseData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data Data   `json:"data"`
}

type RequestBody struct {
	ID string `json:"id"`
}

// https://api.restful-api.dev/objects/4
func GetData(id string) (*ResponseData, error) {

	_url := "https://api.restful-api.dev/objects/" + id

	resp, err := http.Get(_url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res *ResponseData
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteMessage(sqsClient *sqs.Client, ctx context.Context, receiptHandle string) {

	sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(os.Getenv("QUEUE_URL")),
		ReceiptHandle: aws.String(receiptHandle),
	})

}

func TransformRequestPayload(message string) (*RequestBody, string, error) {

	record := new(RecordEvent)
	if err := json.Unmarshal([]byte(message), &record); err != nil {
		return nil, record.MyTaskToken, err
	}

	payload := new(RequestBody)
	if err := json.Unmarshal(record.Input, payload); err != nil {
		return nil, record.MyTaskToken, err
	}

	return payload, record.MyTaskToken, nil
}

func sendTask(client *sfn.Client, ctx context.Context, token string, data []byte, err error) error {

	if err != nil {
		params := sfn.SendTaskFailureInput{
			TaskToken: aws.String(token),
			Error:     aws.String(err.Error()),
		}
		client.SendTaskFailure(ctx, &params)
		return err
	}

	params := sfn.SendTaskSuccessInput{
		TaskToken: aws.String(token),
		Output:    aws.String(string(data)),
	}
	client.SendTaskSuccess(ctx, &params)

	return nil

}

func Handler(ctx context.Context, message string) ([]byte, string, error) {

	evt, token, err := TransformRequestPayload(message)
	if err != nil {
		return nil, token, err
	}

	if evt == nil || evt.ID == "" {
		return nil, token, errors.New("event is nil")
	}

	res, err := GetData(evt.ID)
	if err != nil {
		return nil, token, err
	}

	v, err := json.Marshal(res)
	if err != nil {
		return nil, token, err
	}

	return v, token, nil
}

func event(ctx context.Context, sqsEvent events.SQSEvent) error {

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return err
	}

	// Initialize Step Functions client
	client := sfn.NewFromConfig(cfg)

	// Initialize SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	for _, record := range sqsEvent.Records {

		fmt.Printf("Processing request for %s\n", record.Body)

		defer DeleteMessage(sqsClient, ctx, record.ReceiptHandle)

		data, token, err := Handler(ctx, record.Body)
		if err := sendTask(client, ctx, token, data, err); err != nil {
			return err
		}
	}

	return nil
}

func main() {

	lambda.Start(event)

}
