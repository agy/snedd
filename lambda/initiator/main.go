package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/fullsailor/pkcs7"
)

type Message struct {
	IDDoc string `json:"pkcs7"`
	TTL   uint   `json:"ttl"`
}

func awsCerts() ([]*x509.Certificate, error) {
	const AWSCert = `-----BEGIN CERTIFICATE-----
MIIC7TCCAq0CCQCWukjZ5V4aZzAJBgcqhkjOOAQDMFwxCzAJBgNVBAYTAlVTMRkw
FwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYD
VQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQzAeFw0xMjAxMDUxMjU2MTJaFw0z
ODAxMDUxMjU2MTJaMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9u
IFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNl
cnZpY2VzIExMQzCCAbcwggEsBgcqhkjOOAQBMIIBHwKBgQCjkvcS2bb1VQ4yt/5e
ih5OO6kK/n1Lzllr7D8ZwtQP8fOEpp5E2ng+D6Ud1Z1gYipr58Kj3nssSNpI6bX3
VyIQzK7wLclnd/YozqNNmgIyZecN7EglK9ITHJLP+x8FtUpt3QbyYXJdmVMegN6P
hviYt5JH/nYl4hh3Pa1HJdskgQIVALVJ3ER11+Ko4tP6nwvHwh6+ERYRAoGBAI1j
k+tkqMVHuAFcvAGKocTgsjJem6/5qomzJuKDmbJNu9Qxw3rAotXau8Qe+MBcJl/U
hhy1KHVpCGl9fueQ2s6IL0CaO/buycU1CiYQk40KNHCcHfNiZbdlx1E9rpUp7bnF
lRa2v1ntMX3caRVDdbtPEWmdxSCYsYFDk4mZrOLBA4GEAAKBgEbmeve5f8LIE/Gf
MNmP9CM5eovQOGx5ho8WqD+aTebs+k2tn92BBPqeZqpWRa5P/+jrdKml1qx4llHW
MXrs3IgIb6+hUIB+S8dz8/mmO0bpr76RoZVCXYab2CZedFut7qc3WUH9+EUAH5mw
vSeDCOUMYQR7R9LINYwouHIziqQYMAkGByqGSM44BAMDLwAwLAIUWXBlk40xTwSw
7HX32MxXYruse9ACFBNGmdX2ZBrVNGrN9N2f6ROk0k9K
-----END CERTIFICATE-----
`

	var certs []*x509.Certificate

	decoded, rest := pem.Decode([]byte(AWSCert))
	if len(rest) != 0 {
		return certs, errors.New("invalid AWS cert")
	}

	pub, err := x509.ParseCertificate(decoded.Bytes)
	if err != nil {
		return certs, err
	}

	if pub == nil {
		return certs, errors.New("invalid x509 certificate")
	}

	certs = append(certs, pub)

	return certs, nil
}

func decodeIdentityDocument(IDDoc string) (*ec2metadata.EC2InstanceIdentityDocument, error) {
	// wrap the ID document so that we may PEM decode it
	wrappedDoc := fmt.Sprintf("-----BEGIN PKCS7-----\n%s\n-----END PKCS7-----", IDDoc)

	decoded, rest := pem.Decode([]byte(wrappedDoc))
	if len(rest) != 0 {
		return nil, errors.New("invalid PKCS7 cert")
	}

	data, err := pkcs7.Parse(decoded.Bytes)
	if err != nil {
		return nil, err
	}

	certs, err := awsCerts()
	if err != nil {
		return nil, err
	}

	data.Certificates = certs

	if data.Verify() != nil {
		return nil, errors.New("invalid ID document")
	}

	doc := new(ec2metadata.EC2InstanceIdentityDocument)
	if err := json.Unmarshal([]byte(data.Content), doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func Handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	stateMachineARN := os.Getenv("STATEMACHINEARN")
	if stateMachineARN == "" {
		return nil, errors.New("invalid state machine ARN")
	}

	ttl := os.Getenv("TTL")
	if ttl == "" {
		ttl = "30"
	}

	// While golang handles strconv.Atoi("") correctly, the eawsy shim
	// seems to choke on it.
	expiry, err := strconv.Atoi(ttl)
	if err != nil {
		return nil, err
	}

	msg := new(Message)
	if err := json.Unmarshal(evt, &msg); err != nil {
		return nil, err
	}

	doc, err := decodeIdentityDocument(msg.IDDoc)
	if err != nil {
		return nil, err
	}

	sess := session.Must(session.NewSession())
	svc := sfn.New(sess)

	execInput := fmt.Sprintf(`{"ttl": %d, "instance-id": "%s"}`, expiry, doc.InstanceID)
	execName := fmt.Sprintf("snedd-%s", doc.InstanceID)

	input := &sfn.StartExecutionInput{
		Input:           aws.String(execInput),
		Name:            aws.String(execName),
		StateMachineArn: aws.String(stateMachineARN),
	}

	res, err := svc.StartExecution(input)
	if err != nil {
		return nil, err
	}

	log.Println(res)

	return res, nil
}

func main() {}
