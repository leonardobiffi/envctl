package main

import (
	"fmt"

	"github.com/leonardobiffi/envctl/cli"
	"github.com/leonardobiffi/envctl/version"
)

func main() {
	info := &cli.Info{
		Name:        "envctl",
		Version:     version.Version,
		Description: "Is a simple CLI tool which runs a process with secrets from AWS Secrets Manager",
		AuthorName:  "Leonardo Biffi",
		AuthorEmail: "leonardobiffi@outlook.com",
	}

	err := cli.Initialize(info)
	if err != nil {
		fmt.Println(err)
	}
}
