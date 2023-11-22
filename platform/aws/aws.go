package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/leonardobiffi/envctl/util/system/exit"
)

// TODO: add service function

// GetSecrets returns a map of secrets from AWS Secrets Manager.
func GetSecrets(profile string, region string, secretName string) map[string]string {
	var opts []func(*config.LoadOptions) error
	ctx := context.TODO()
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		exit.Error(err.Error())
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		exit.Error(err.Error())
	}

	secretString := *result.SecretString
	data := map[string]string{}

	json.Unmarshal([]byte(secretString), &data)
	return data
}

// UpdateSecrets adds a secret to AWS Secrets Manager from .env file
func UpdateSecrets(profile string, region string, secretName string, data map[string]string) {
	var opts []func(*config.LoadOptions) error
	ctx := context.TODO()
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		exit.Error(err.Error())
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	secretString, _ := json.Marshal(data)

	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(secretName),
		SecretString: aws.String(string(secretString)),
	}

	_, err = svc.UpdateSecret(ctx, input)
	if err != nil {
		exit.Error(err.Error())
	}
}
