package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

const (
	banner = `
          _  __       _           _                   _
 ___  ___| |/ _|   __| | ___  ___| |_ _ __ _   _  ___| |_
/ __|/ _ \ | |_   / _' |/ _ \/ __| __| '__| | | |/ __| __|
\__ \  __/ |  _| | (_| |  __/\__ \ |_| |  | |_| | (__| |_
|___/\___|_|_|    \__,_|\___||___/\__|_|   \__,_|\___|\__|


 ___  ___  __ _ _   _  ___ _ __   ___ ___
/ __|/ _ \/ _' | | | |/ _ \ '_ \ / __/ _ \
\__ \  __/ (_| | |_| |  __/ | | | (_|  __/
|___/\___|\__, |\__,_|\___|_| |_|\___\___|
             |_|
 _       _ _   _       _           _
(_)_ __ (_) |_(_) __ _| |_ ___  __| |
| | '_ \| | __| |/ _' | __/ _ \/ _' |
| | | | | | |_| | (_| | ||  __/ (_| |
|_|_| |_|_|\__|_|\__,_|\__\___|\__,_|`
)

type instance struct {
	Cert   string `json:"pkcs7"`
	ID     string `json:"instance-id"`
	Region string `json:"-"`
}

func InstanceMeta() (*instance, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	meta := ec2metadata.New(sess)

	pkcs7, err := meta.GetDynamicData("instance-identity/pkcs7")
	if err != nil {
		return nil, err
	}

	doc, err := meta.GetInstanceIdentityDocument()
	if err != nil {
		return nil, err
	}

	return &instance{
		Cert:   pkcs7,
		ID:     doc.InstanceID,
		Region: doc.Region,
	}, nil
}

func Invoke(i *instance, fnName string) (*lambda.InvokeOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(i.Region),
	})
	if err != nil {
		return nil, err
	}

	svc := lambda.New(sess)

	payload, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(fnName),
		InvocationType: aws.String("RequestResponse"),
		Payload:        payload,
	}

	res, err := svc.Invoke(input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	var (
		fnName = flag.String("fn-name", "snedd-initiator", "Lambda function name")
		runDir = flag.String("run-dir", "/run/snedd", "Temporary storage directory")
	)
	flag.Parse()

	inst, err := InstanceMeta()
	if err != nil {
		panic(err)
	}

	res, err := Invoke(inst, *fnName)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	// get ttl

	fmt.Println(*runDir)
	fmt.Println(banner)
}
