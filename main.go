package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func main() {

	address := os.Getenv("address")

	parts := strings.Split(address, ".")
	accountID := parts[0]
	region := parts[3]

	sess := session.Must(session.NewSession())
	ecrClient := ecr.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	out := &output{
		Error:         "",
		State:         "success",
		ServerAddress: address,
	}
	defer finish(out)

	authTokenResp, err := ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(accountID)},
	})
	if err != nil {
		out.Error = err.Error()
		out.State = "error"
		return
	}

	if len(authTokenResp.AuthorizationData) < 1 {
		out.Error = "no ecr registry credentials found"
		out.State = "error"
		return
	}

	aD := authTokenResp.AuthorizationData[0]

	base := aws.StringValue(aD.AuthorizationToken)
	creds, _ := base64.StdEncoding.DecodeString(base)
	credParts := strings.Split(string(creds), ":")
	username := credParts[0]
	password := credParts[1]

	out.Username = username
	out.Password = password
}

func finish(out *output) {
	o, _ := json.Marshal(out)

	fmt.Printf("%s\n", o)
}

type output struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ServerAddress string `json:"serverAddress"`
	Error         string `json:"error"`
	State         string `json:"state"`
}
