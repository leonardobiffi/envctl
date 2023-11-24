package parameters

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/leonardobiffi/envctl/config"
	"github.com/leonardobiffi/envctl/platform/aws"
	"github.com/leonardobiffi/envctl/util/system/exit"
)

func GetParameters(path string, env string, region string, profile string, envFile string) map[string]string {
	if envFile != "" {
		parameters, err := godotenv.Read(envFile)

		if err != nil {
			exit.Error("Could not read env file " + envFile)
		}

		return parameters
	}

	conf := config.GetConfig()

	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}

	if env == "" {
		env = conf.DefaultEnvironment
	}

	if path == "" {
		path = os.Getenv("PARAMETER_NAME")
	}

	if path == "" && env == "" {
		exit.Error("Parameter Name is required to list environments. Set -parameter flag.")
	}

	if path == "" && env != "" {
		if _, ok := conf.Environments[env]; !ok {
			exit.Error("Environment '" + env + "' does not exist.")
		}

		path = conf.Environments[env]
	}

	if region == "" {
		region = conf.Region
	}

	if profile == "" {
		profile = conf.Profile
	}

	return aws.GetParameters(profile, region, path)
}

// UpdateParameters sets appropriate config and updates parameters on aws.
func UpdateParameters(parameterPath string, env string, region string, profile string, envFile string) {
	if envFile == "" {
		exit.Error("Env file is required to update parameters. Set --envfile flag.")
	}

	if envFile != "" {
		parameters, err := godotenv.Read(envFile)

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

		if parameterPath == "" {
			parameterPath = os.Getenv("PARAMETER_NAME")
		}

		if parameterPath == "" && env == "" {
			exit.Error("Parameter Path is required to list environments. Set -parameter flag.")
		}

		if parameterPath == "" && env != "" {
			if _, ok := conf.Environments[env]; !ok {
				exit.Error("Environment '" + env + "' does not exist.")
			}

			parameterPath = conf.Environments[env]
		}

		if region == "" {
			region = conf.Region
		}

		if profile == "" {
			profile = conf.Profile
		}

		aws.UpdateParameters(profile, region, parameterPath, parameters)
	}
}
