package main

import (
	"fmt"
	"os"
	"time"

	"github.com/notnotquinn/wallpapers/conf"
	"github.com/notnotquinn/wallpapers/wallpapers"
	wp "github.com/reujab/wallpaper"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:                   "Wallpaper Updater",
		Version:                "0.0.1-dev",
		UseShortOptionHandling: true,
		Usage:                  "Randomize your wallpapers",
		Copyright:              "(c) MIT",
		Flags: []cli.Flag{
			&cli.StringFlag{
				// config file
				Name:        "config",
				Aliases:     []string{"c"},
				DefaultText: "./config.json",
				TakesFile:   true,
				Usage:       "config file",
			},
		},
		Authors: []*cli.Author{{
			Name:  "Quinn T",
			Email: "quinn.github@gmail.com",
		}},
		Before: func(c *cli.Context) error {
			if c.IsSet("config") {
				err := conf.Load(c.String("config"))
				if err != nil {
					return fmt.Errorf("config not loaded: %w", err)
				}
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "debug",
				Usage: "Debug command for private use",
				Action: func(c *cli.Context) error {
					for range make([]struct{}, 4) {
						fmt.Println(conf.Conf)
						cat := []wallpapers.WallpaperCatagory{wallpapers.AsiaRussia}
						link, _, err := wallpapers.NewURL(cat, cat)
						if err != nil {
							return err
						}
						err = wp.SetFromURL(link)
						if err != nil {
							return err
						}
						time.Sleep(time.Second * 5)
					}
					return nil
				},
			},
		},
	}

	// Some things that need to get done:
	// - SET WALLPAPERS EVERY SO OFTEN AND GET THEM FROM THE INTERNET
	// - Add ourselves to start on boot. (cross platform may be difficult)
	// - Some configuration options
	// - NSFW detection and block toggle.
	//   they provide API key to 3rd party service if they want it - no backend

	// Some things that would be nice:
	// - CLI can be used to interact with the currently running thing (like docker)
	//   Needs 2 binaries (?) one is a local server and the CLI can connect to it and invoke things.
	//   Seems like way too much work, could be easier to just set one in watch mode
	//   on the config, that way the CLI can edit the config file and the other will
	//   pick up on that and change accordingly.
	err := app.Run(os.Args)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	err := wp.SetMode(wp.Fit)
	if err != nil {
		panic(err)
	}
}
