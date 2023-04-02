package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "nre",
		Usage:                  "regex expression with natural language",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Suggest:                true,

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "detail",
				Aliases: []string{"d"},
				Usage:   "with explanation of the regex expression",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "api",
				Aliases:  []string{"key"},
				EnvVars:  []string{"OPENAI_API_KEY"},
				Usage:    "specify the OpenAI API key, it will be set to env var $OPENAI_API_KEY by default",
			},
		},

		Action: func(ctx *cli.Context) error {
			query := "pattern: " + ctx.Args().First()

			messages := []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a helpful programmer assistant, the user describe a pattern and you answer with a regex matching that pattern"},
				{Role: "system", Content: "Output the regex expression only, no any explanation, unless the user requires you to output the explanation"},
				{Role: "user", Content: query},
			}
			if ctx.Bool("detail") {
				messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: "with explanation"})
			}

			client := openai.NewClient(ctx.String("api"))
			stream, err := client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
				Model:       "gpt-3.5-turbo",
				Messages:    messages,
				Temperature: 0.1,
				Stream:      true,
			})
			if err != nil {
				return err
			}
			defer stream.Close()

			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}

				if err != nil {
					return err
				}

				fmt.Print(response.Choices[0].Delta.Content)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
