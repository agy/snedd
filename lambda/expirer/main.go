package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

type Message struct {
	InstanceID string `json:"instance-id"`
}

func Handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	msg := new(Message)
	if err := json.Unmarshal(evt, &msg); err != nil {
		return nil, err
	}

	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	res, err := svc.TerminateInstances(
		&ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(msg.InstanceID),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	log.Println(res)

	return res, nil
}

func main() {}
