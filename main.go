package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/thinceller/scrapgox/client"
)

var commands = []*cli.Command{
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list scrapbox pages",
		Action:  cmdList,
		Flags:   flags,
	},
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "token",
		Aliases: []string{"t"},
		Usage:   "Set token for request (Cookie: connect.sid)",
	},
}

func cmdList(c *cli.Context) error {
	parsedUrl, err := url.ParseRequestURI(client.DefaultHost)
	if err != nil {
		return err
	}

	token := c.String("token")
	scrapgoxClient, err := client.NewClient(parsedUrl, token, client.DefaultUserAgent)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// TODO: cli の context から project を取得する
	project := c.Args().Get(0)
	// TODO: cli の context から query を取得する
	query := c.Args().Get(1)
	pages, err := scrapgoxClient.GetPages(project, query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, v := range pages {
		fmt.Println(v.Title)
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:     "scrapgox",
		Usage:    "scrapbox cli tool",
		Commands: commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
