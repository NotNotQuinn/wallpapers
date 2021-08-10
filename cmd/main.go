package main

import (
	"fmt"
	"os"

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
				Name:        "config",
				Aliases:     []string{"c"},
				TakesFile:   true,
				Value:       "./wallpaperconf.json",
				DefaultText: "./wallpaperconf.json",
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
					fmt.Printf("Unable to load config in '%s'\nTo set the path use the -c flag.\n", c.String("config"))
					os.Exit(1)
				}
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "random",
				Usage: "Changes the wallpaper to a random one",
				Action: func(c *cli.Context) error {
					return wallpapers.ChangeToRandom()
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
	// "zooms" into the images until they fill the screen
	must(wp.SetMode(wp.Crop))
}
