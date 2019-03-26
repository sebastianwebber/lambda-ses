package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	// fmt.Printf("ctx: %#v\n", ctx)
	// fmt.Printf("snsEvent: %#v\n", snsEvent)
	for _, record := range snsEvent.Records {

		msg, err := parseMessage(record.SNS.Message)

		if err != nil {
			fmt.Printf("ERROR: could not parse message: %v\n", err)
			return
		}

		fmt.Printf("msg: %#v\n", msg)

		err = saveItem(msg)
		fmt.Printf("err: %#v\n", err)

		if err != nil {
			fmt.Printf("ERROR: could not save item on Dynamo: %v\n", err)
		}
	}
}

const dynamoTargetTable = "mailing"

var svc = dynamodb.New(session.New())

func saveItem(msg notificationMessage) error {

	for i := 0; i < len(msg.Bounce.Recipients); i++ {

		input := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: &msg.Bounce.Recipients[i].EmailAddress,
				},
				"notificationType": {
					S: &msg.Bounce.Type,
				},
				"from": {
					S: &msg.Mail.Source,
				},
				"state": {
					S: aws.String("disable"),
				},
				"timestamp": {
					S: &msg.Mail.TimeStamp,
				},
			},
			TableName: aws.String(dynamoTargetTable),
		}

		fmt.Printf("input: %#v\n", input)

		_, err := svc.PutItem(input)

		return err
	}

	return nil
}

type notificationMessage struct {
	Type   string `json:"notificationType"`
	Bounce bounce `json:"bounce"`
	Mail   mail   `json:"mail"`
}

type bounce struct {
	Type       string      `json:"bounceType"`
	Recipients []recipient `json:"bouncedRecipients"`
}

type recipient struct {
	EmailAddress string `json:"emailAddress"`
}

type mail struct {
	MessageID string `json:"messageId"`
	Source    string `json:"source"`
	TimeStamp string `json:"timestamp"`
}

func parseMessage(message string) (notificationMessage, error) {

	var out notificationMessage

	err := json.Unmarshal([]byte(message), &out)

	return out, err
}

func main() {
	lambda.Start(handler)
}
