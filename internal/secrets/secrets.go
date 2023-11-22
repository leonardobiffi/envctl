package secrets

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/leonardobiffi/envctl/config"
	"github.com/leonardobiffi/envctl/platform/aws"
	"github.com/leonardobiffi/envctl/util/system/exit"
)

// TODO: add service function

// GetSecrets sets appropriate config and fetches secrets from aws
func GetSecrets(secretName string, env string, region string, profile string, envFile string) map[string]string {
	if envFile != "" {
		secrets, err := godotenv.Read(envFile)

		if err != nil {
			exit.Error("Could not read env file " + envFile)
		}

		return secrets
	}

	conf := config.GetConfig()

	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}

	if env == "" {
		env = conf.DefaultEnvironment
	}

	if secretName == "" {
		secretName = os.Getenv("SECRET_NAME")
	}

	if secretName == "" && env == "" {
		exit.Error("Secret Name is required to list environments. Set -secret flag.")
	}

	if secretName == "" && env != "" {
		if _, ok := conf.Environments[env]; !ok {
			exit.Error("Environment '" + env + "' does not exist.")
		}

		secretName = conf.Environments[env]
	}

	if region == "" {
		region = conf.Region
	}

	if profile == "" {
		profile = conf.Profile
	}

	return aws.GetSecrets(profile, region, secretName)
}

// UpdateSecrets sets appropriate config and updates secrets on aws.
func UpdateSecrets(secretName string, env string, region string, profile string, envFile string) {
	if envFile == "" {
		exit.Error("Env file is required to update secrets. Set --envfile flag.")
	}

	if envFile != "" {
		secrets, err := godotenv.Read(envFile)

		if err != nil {
			exit.Error("Could not read env file " + envFile)
		}

		conf := config.GetConfig()

		if env == "" {
			env = os.Getenv("ENVIRONMENT")
		}

		if env == "" {
			env = conf.DefaultEnvironment
		}

		if secretName == "" {
			secretName = os.Getenv("SECRET_NAME")
		}

		if secretName == "" && env == "" {
			exit.Error("Secret Name is required to list environments. Set -secret flag.")
		}

		if secretName == "" && env != "" {
			if _, ok := conf.Environments[env]; !ok {
				exit.Error("Environment '" + env + "' does not exist.")
			}

			secretName = conf.Environments[env]
		}

		if region == "" {
			region = conf.Region
		}

		if profile == "" {
			profile = conf.Profile
		}

		aws.UpdateSecrets(profile, region, secretName, secrets)
	}
}
