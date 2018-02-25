package main

import (
	"fmt"
	"log"
	"os"

	"github.com/er28-0652/go-zoomus/zoomus"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type Repo struct {
	Owner string
	Name  string
}

type Build struct {
	Tag      string
	Event    string
	Number   int
	Commit   string
	Ref      string
	Branch   string
	Author   string
	Pull     string
	Message  string
	DeployTo string
	Status   string
	Link     string
	Started  int64
	Created  int64
}

type Config struct {
	Webhook string
	Token   string
}

type Job struct {
	Started int64
}

type Plugin struct {
	Repo   Repo
	Build  Build
	Config Config
	Job    Job
}

func (p *Plugin) Exec() error {
	msg := zoomus.Message{
		Title:   "drone notification",
		Summary: message(p.Repo, p.Build),
		Body:    fallback(p.Repo, p.Build),
	}

	zoom, err := zoomus.NewClient(p.Config.Webhook, p.Config.Token)
	if err != nil {
		return errors.Wrap(err, "fail to init Client")
	}

	log.Printf("Webhook: %s\n", zoom.WebhookURL.String())
	log.Printf("Token: %s\n", zoom.Header["X-Zoom-Token"])
	log.Printf("Title: %s\n", msg.Title)
	log.Printf("Summary: %s\n", msg.Summary)
	log.Printf("Body: %s\n", msg.Body)

	err = zoom.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err, "fail to send message")
	}
	log.Println("done")
	return nil
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* [%s|%s/%s#%s] (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Commit:   c.String("commit.sha"),
			Ref:      c.String("commit.ref"),
			Branch:   c.String("commit.branch"),
			Author:   c.String("commit.author"),
			Pull:     c.String("commit.pull"),
			Message:  c.String("commit.message"),
			DeployTo: c.String("build.deployTo"),
			Link:     c.String("build.link"),
			Started:  c.Int64("build.started"),
			Created:  c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Webhook: c.String("webhook"),
			Token:   c.String("token"),
		},
	}

	err := plugin.Exec()
	if err != nil {
		return cli.NewExitError(err, -1)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "zoom plugin"
	app.Action = run
	app.Flags = []cli.Flag{
		// zoom config variables
		cli.StringFlag{
			Name:   "webhook",
			Usage:  "zoom webhook url",
			EnvVar: "ZOOM_WEBHOOK,PLUGIN_WEBHOOK",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "zoom webhook token",
			EnvVar: "ZOOM_TOKEN,PLUGIN_TOKEN",
		},

		// drone envs
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.pull",
			Usage:  "git pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:   "build.deployTo",
			Usage:  "environment deployed to",
			EnvVar: "DRONE_DEPLOY_TO",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},
	}

	app.Run(os.Args)
}
