# Envctl
![GitHub release](https://img.shields.io/github/release/leonardobiffi/envctl.svg?style=flat)
![GitHub](https://img.shields.io/github/license/leonardobiffi/envctl.svg?style=flat)

A simple CLI tool to run a process with secrets from AWS Secrets Manager or AWS Parameter Store

forked from [pratishshr/envault](https://github.com/pratishshr/envault)

## About

Envctl focuses on integrating AWS Secrets Manager in your application with ease without having to write a single line of code in your source files. Simply run your commands with the envctl CLI and the secrets will be injected in that process.

## Table Of Contents
1. [Install envctl](#1-install-envctl)
2. [Verify Installation](#2-verify-installation)
3. [AWS Credentials](#3-aws-credentials)
4. [Setup](#4-setup)
5. [List Secrets](#5-list-secrets)
5. [Update Secrets](#6-update-secrets)
7. [Run With Secrets](#7-run-with-secrets)
8. [Usage with CI/CD](#8-usage-with-cicd)
9. [Using custom .env files](#89-using-custom-env-files)

## Usage

### 1. Install envctl:

```shell
curl -fsSL https://raw.githubusercontent.com/leonardobiffi/envctl/master/scripts/install.sh | sh
```

Note: 
If your architecture is not supported, clone this repo and run `go build` to generate a binary.
Then, simply place the binary in your local `bin`.


### 2. Verify Installation:

```shell
envctl
```

### 3. AWS Credentials

Before using envctl, you have to provide your AWS credentials. This allows envctl to fetch secrets from the AWS Secrets Manager. Also, make sure you have the correct access for your credentials.

Simply create `~/.aws/credentials` file for storing AWS credentials. <br/>
Example:

```
[example-profile]
aws_access_key_id = xxxxxx
aws_secret_access_key = xxxxxx
```
To know more about AWS configurations, view [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)

### 4. Setup

Go to your project directory and run `setup` command to initiate the setup process.

```shell
envctl setup
```

- Choose your AWS profile that was setup earlier. <br>
- Choose the AWS Region where your secrets are kept.
- You can also add a deployment environment associated with the secret name. You may add any number of environment you want.
- Set a default env

```
 Example:

 AWS profile: default
 Region: US West (Oregon)
 Add an environment (eg. dev): dev
 Secret Name: api/dev
 Add an environment (eg. dev): uat
 Secret Name: api/uat
```

`envctl.json` file will be created in your project directory like below.
```json
{
  "profile": "default",
  "region": "us-west-2",
  "environments": {
    "dev": "api/dev",
    "uat": "api/uat"
  },
  "defaultEnv": "dev"
}
```

**If you do not want a project-specific config file, you can skip the above step.**

### 5. List secrets

```shell
envctl list -e dev
```
```shell
envctl list -e uat
```
Here `dev` and `uat` are the environments you specified in `envctl.json`.


If you have not setup a `envctl.json` file, you can still pass `--secret` or `-s` flag with the secrets path.
This will use the `default` profile from your `~/.aws/credentials` file.
```shell
envctl list --secret=api/dev
```
```shell
envctl list --secret=api/uat
```

You also can list environments from AWS Parameter Store

```shell
envctl list --parameter=api/dev
```
```shell
envctl list --parameter=api/uat
```

### 6. Update secrets

This will update secrets with content in .env file

```shell
envctl update --secret=/dev/service/app --envfile .env
```

Or update secret on Parameter Store
```shell
envctl update --parameter=/dev/service/app --envfile .env
```

### 7. Run with secrets

```shell
envctl run 'yarn build' -e dev
```
This will inject the secrets from `dev` to the `yarn build` process.

Similarly, if you have not setup a `envctl.json` file, you can still pass `--secret` or `-s` flag with the secrets path.
This will use the `default` profile from your `~/.aws/credentials` file.

```shell
envctl run 'yarn build' --secret=api/dev
```

### 8. Usage with CI/CD:

Instead of setting up a `~/.aws/credentials` file. You can also use the following environment variables to set up your AWS credentials.

| Variable | Description |
|-----------|----------|
| AWS_ACCESS_KEY_ID | Your AWS access key|
| AWS_SECRET_ACCESS_KEY | Your AWS secret key|
| AWS_REGION | AWS region where you added your secret|
| ENVIRONMENT | Environment which you set in envctl.json |
| SECRET_NAME | AWS Secret Name |
| PARAMETER_NAME | AWS Parameter Store Path |

### 9. Using custom .env files
If you want to inject environment keys from a file instead of using AWS Secrets Manager. You can use the`-ef` flag.

```shell
envctl run 'envctl run 'go run main.go' -ef env/staging.env
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
