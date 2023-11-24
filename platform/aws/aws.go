package aws

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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

// Get Parameters returns a map of parameters from AWS Systems Manager Parameter Store.
func GetParameters(profile string, region string, path string) map[string]string {
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

	// Create Parameter Store client
	svc := ssm.NewFromConfig(cfg)

	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	result, err := svc.GetParametersByPath(ctx, input)
	if err != nil {
		exit.Error(err.Error())
	}

	data := map[string]string{}
	for _, parameter := range result.Parameters {
		data[sanitizePath(*parameter.Name, path)] = *parameter.Value
	}

	for result.NextToken != nil {
		input.NextToken = result.NextToken
		result, err = svc.GetParametersByPath(ctx, input)
		if err != nil {
			exit.Error(err.Error())
		}

		for _, parameter := range result.Parameters {
			data[sanitizePath(*parameter.Name, path)] = *parameter.Value
		}
	}

	return data
}

// UpdateParameters adds a parameter to AWS Systems Manager Parameter Store from .env file
func UpdateParameters(profile string, region string, path string, data map[string]string) {
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

	// Create Parameter Store client
	svc := ssm.NewFromConfig(cfg)

	for key, value := range data {
		input := &ssm.PutParameterInput{
			Name:      aws.String(formatPath(path) + key),
			Value:     aws.String(value),
			Overwrite: aws.Bool(true),
			Type:      types.ParameterTypeString,
		}

		_, err = svc.PutParameter(ctx, input)
		if err != nil {
			exit.Error(err.Error())
		}
	}
}

// remove prefix of path name
func sanitizePath(name, path string) string {
	new := strings.ReplaceAll(name, path, "")
	return strings.TrimPrefix(new, "/")
}

// add suffix to path
func formatPath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}
