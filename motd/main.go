package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

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

type response struct {
	ID  string `json:"instance-id"`
	TTL int    `json:"ttl"`
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

func NewErrLog() *log.Logger {
	return log.New(os.Stderr, "", 0)
}

func main() {
	var (
		fnName = flag.String("fn-name", "snedd-initiator", "Lambda function name")
		runDir = flag.String("run-dir", "/run/snedd", "Temporary storage directory")
	)
	flag.Parse()

	sem := fmt.Sprintf("%s/triggered", *runDir)

	// Exit if we've run before
	if _, err := os.Stat(sem); err == nil {
		os.Exit(0)
	}

	ErrLog := NewErrLog()

	inst, err := InstanceMeta()
	if err != nil {
		ErrLog.Fatal(err)
	}

	res, err := Invoke(inst, *fnName)
	if err != nil {
		ErrLog.Fatal(err)
	}

	var success int64 = 200
	if *res.StatusCode != success {
		ErrLog.Fatal(errors.New("lambda invocation failed"))
	}

	var payload response
	if err := json.Unmarshal(res.Payload, &payload); err != nil {
		ErrLog.Fatal(errors.New("could not decode response"))
	}

	// Write the semaphore file and ignore failures
	ioutil.WriteFile(sem, []byte(strconv.Itoa(payload.TTL)), 0444)

	fmt.Printf("payload: %+v\n", payload)
	fmt.Println(banner)
}
