package aws

import (
	"context"
	"errors"
	"fmt"
	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
)

type Profiles struct {
	Data  []string
	Error string
}
type credentialProvider struct {
	awsV2.Credentials
}

func (c credentialProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey, SessionToken: os.Getenv("AWS_SESSION_TOKEN")}, nil
}

func (c credentialProvider) IsExpired() bool {
	return c.Expired()
}

func GetCfg(profile, region string) (awsV2.Config, error) {
	cfg, err := awsV2Config.LoadDefaultConfig(
		context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("failed to load config")
		return awsV2.Config{}, err
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("failed to read credentials")
		return awsV2.Config{}, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return awsV2.Config{}, errors.New("AWS Credentials expired")
	}
	return cfg, err
}
