package cli

import (
	"fmt"
	"strings"

	"github.com/leonardobiffi/envctl/internal/cli/setup"
	"github.com/leonardobiffi/envctl/internal/secrets"
	"github.com/leonardobiffi/envctl/util/shell"
	"github.com/leonardobiffi/envctl/util/system/exit"
)

// Setup prompts user for required settings and creates a envctl.json file
func Setup() {
	setup.Run()
}

// Run given command with the secrets from given Secret Manager.
func Run(secretName string, command string, env string, region string, profile string, envFile string) {
	if command == "" {
		exit.Error("Command to run is not specified. Add command as 'envctl run [command]'")
	}

	shell.Execute(command, secrets.GetSecrets(secretName, env, region, profile, envFile))
}

// List all environment from Secrets Manager
func List(secretName string, env string, region string, profile string, envFile string, upper bool) {
	for key, value := range secrets.GetSecrets(secretName, env, region, profile, envFile) {
		if upper {
			key = strings.ToUpper(key)
		}

		fmt.Println(key + "=" + value)
	}
}

// Update all environment from env file to Secrets Manager
func Update(secretName string, env string, region string, profile string, envFile string) {
	secrets.UpdateSecrets(secretName, env, region, profile, envFile)
}
