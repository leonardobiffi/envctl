package setup

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/leonardobiffi/envctl/config"
	"github.com/leonardobiffi/envctl/platform/aws"
)

type awsAnswers struct {
	Profile string
	Region  string
}

type envAnswers struct {
	Environment string
	Secret      string
}

func askAwsQuestions() *awsAnswers {
	questions := []*survey.Question{
		{
			Name: "profile",
			Prompt: &survey.Select{
				Message: "AWS profile:",
				Options: aws.GetProfiles(),
			},
		},
		{
			Name: "region",
			Prompt: &survey.Select{
				Message: "Region:",
				Options: aws.GetRegions(),
			},
		},
	}

	answers := &awsAnswers{}

	err := survey.Ask(questions, answers)
	answers.Region = aws.GetRegionCode(answers.Region)

	if err != nil {
		fmt.Println(err)
	}

	return answers
}

func initEnvQuestions(environments []*envAnswers) []*envAnswers {
	questions := []*survey.Question{
		{
			Name: "environment",
			Prompt: &survey.Input{
				Message: "Add an environment (eg. dev):",
			},
		},
		{
			Name: "secret",
			Prompt: &survey.Input{
				Message: "Secret Name:",
			},
		},
	}

	answers := &envAnswers{}
	err := survey.Ask(questions, answers)

	if err != nil {
		fmt.Println(err)
	}

	environments = append(environments, answers)
	addMoreEnv := false
	confirm := &survey.Confirm{
		Message: "Would you like to add another environment?",
	}

	survey.AskOne(confirm, &addMoreEnv)

	if !addMoreEnv {
		return environments
	}

	return initEnvQuestions(environments)
}

func askEnvQuestions() map[string]string {
	environments := []*envAnswers{}
	environmentsMap := map[string]string{}

	environments = initEnvQuestions(environments)

	for _, environment := range environments {
		environmentsMap[environment.Environment] = environment.Secret
	}

	return environmentsMap
}

func askDefaultEnvQuestion(envAnswers map[string]string) string {
	envNames := make([]string, 0, len(envAnswers))
	for k := range envAnswers {
		envNames = append(envNames, k)
	}

	defaultEnvQuestion := &survey.Select{
		Message: "Set a default env:",
		Options: envNames,
	}

	var answer string
	survey.AskOne(defaultEnvQuestion, &answer)

	return answer
}

// Run starts the setup instructions.
func Run() {
	awsAnswers := askAwsQuestions()
	envAnswers := askEnvQuestions()

	defaultEnvAnswer := askDefaultEnvQuestion(envAnswers)

	configuration := &config.Config{
		Profile:            awsAnswers.Profile,
		Region:             awsAnswers.Region,
		Environments:       envAnswers,
		DefaultEnvironment: defaultEnvAnswer,
	}

	config.CreateConfig(configuration)
}
