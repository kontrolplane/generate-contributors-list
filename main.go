package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

type Configuration struct {
	Token      string `env:"INPUT_TOKEN"`
	Owner      string `env:"INPUT_OWNER"`
	Repository string `env:"INPUT_REPOSITORY"`
	Size       int    `env:"INPUT_SIZE" envDefault:"50"`
	File       string `env:"INPUT_FILE" envDefault:"README.md"`
	Limit      int    `env:"INPUT_LIMIT" envDefault:"70"`
	LogLevel   string `env:"INPUT_LOG_LEVEL" envDefault:"info"`
}

type Contributor struct {
	Username string
	Avatar   string
	Profile  string
}

const guard string = "[//]: kontrolplane/contributors"

func main() {
	cfg := Configuration{}
	ctx := context.Background()

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("unable to parse the environment variables", err)
		os.Exit(1)
	}

	var ll slog.Level

	if cfg.LogLevel == "debug" {
		ll = slog.LevelDebug
	} else {
		ll = slog.LevelInfo
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     ll,
	})
	logger := slog.New(logHandler)

	var client *github.Client

	if cfg.Token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	contributors, err := fetchContributors(cfg, ctx, client)
	if err != nil {
		logger.Error("error fetching contributors", slog.Any("error", err))
		os.Exit(1)
	}

	file, err := os.ReadFile(cfg.File)
	if err != nil {
		logger.Error("error reading file", slog.Any("error", err))
		os.Exit(1)
	}

	updatedContent, err := updateContributors(cfg, contributors, string(file))
	if err != nil {
		logger.Error("updating contributors failed", slog.Any("error", err))
		os.Exit(1)
	}

	err = os.WriteFile(cfg.File, []byte(updatedContent), 0644)
	if err != nil {
		logger.Error("error writing file", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Debug(updatedContent)
	logger.Info("contributors section updated successfully")
}

func fetchContributors(cfg Configuration, ctx context.Context, client *github.Client) ([]Contributor, error) {

	clientOptions := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: cfg.Limit},
	}

	var contributors []Contributor
	for {
		contributorList, resp, err := client.Repositories.ListContributors(ctx, cfg.Owner, cfg.Repository, clientOptions)
		if err != nil {
			return nil, err
		}

		for _, contributor := range contributorList {
			contributors = append(contributors, Contributor{
				Username: contributor.GetLogin(),
				Avatar:   contributor.GetAvatarURL(),
				Profile:  fmt.Sprintf("https://github.com/%s", contributor.GetLogin()),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		clientOptions.Page = resp.NextPage
	}

	return contributors, nil
}

func updateContributors(cfg Configuration, contributors []Contributor, content string) (string, error) {

	r := regexp.MustCompile(fmt.Sprintf(`(?s)%s\n*.*?\n*%s`, regexp.QuoteMeta(guard), regexp.QuoteMeta(guard)))

	contributorsHTML := generateContributors(cfg, contributors)

	if r.MatchString(content) {
		content = r.ReplaceAllString(content, fmt.Sprintf("%s\n\n%s\n\n%s", guard, contributorsHTML, guard))
	} else {
		// todo: possible option to add it to the content regardless of guards
		// content += fmt.Sprintf("\n\n%s\n%s\n%s", guard, contributorsHTML, guard)
		return "", errors.New("no guards found used for protecting the markdown file")
	}

	return content, nil
}

func generateContributors(cfg Configuration, contributors []Contributor) string {
	var htmlBuilder strings.Builder
	for i, contributor := range contributors {
		htmlBuilder.WriteString(fmt.Sprintf(
			`<a href="%s"><img src="%s" title="%s" width="%d" height="%d"></a>%s`,
			contributor.Profile,
			contributor.Avatar,
			contributor.Username,
			cfg.Size,
			cfg.Size,

			func() string {
				if i < len(contributors)-1 {
					return "\n"
				}
				return ""
			}(),
		))
	}
	return htmlBuilder.String()
}
