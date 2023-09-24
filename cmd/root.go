package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	awsS "github.com/aws/aws-sdk-go/aws"
	"github.com/mattn/go-colorable"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/view"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	profile, region string
	version         = "dev"
	commit          = "dev"
	date            = "today"
	rootCmd         = &cobra.Command{
		Use:   `cloudlens`,
		Short: `cli for aws services`,
		Long:  `cli for aws services[s3, ec2, security-groups, iam]`,
		Run:   run,
	}
	out = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(versionCmd(), updateCmd())
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Read aws profile")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "Read aws region")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) {
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile("./cloudlens.log", mod, 0644)
	if err != nil {
		log.Printf("Could not open cloudlens.log. Writing logs to stdout instead.")
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if err == nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	}
	//TODO profiles and regions should under aws
	//var sess *session.Session
	var regions []string
	app := view.NewApp()

	profile := awsS.String(os.Getenv(AWS_PROFILE))
	profiles := []string{*profile}
	region := awsS.String(os.Getenv(AWS_DEFAULT_REGION))
	regions = append(regions, *region)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("aws session init failed -- %v", err))
	}
	ctx := context.WithValue(context.Background(), internal.KeySession, cfg)
	if err := app.Init(ctx, profiles, regions, version); err != nil {
		panic(fmt.Sprintf("app init failed -- %v", err))
	}
	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}
}
