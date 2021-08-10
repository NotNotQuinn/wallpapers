package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/notnotquinn/wallpapers/conf"
	"github.com/notnotquinn/wallpapers/wallpapers"
	wp "github.com/reujab/wallpaper"
	"github.com/urfave/cli/v2"
)

func main() {
	loadConfig := func(c *cli.Context) error {
		err := conf.SetPath(c.Path("config"))
		if err != nil && strings.HasSuffix(err.Error(), ": The system cannot find the file specified.") {
			fmt.Println("This subcommand requires flag \"config\" to be set.")
			os.Exit(1)
		}
		return err
	}

	app := cli.NewApp()
	app.Name = "Wallpaper Updater"
	app.Version = "0.0.1-dev"
	app.Usage = "Randomize your wallpapers"
	app.Description = `Randomly set your wallpaper background.`
	app.Copyright = "(c) MIT"
	app.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			DefaultText: conf.DefaultPath,
			Usage:       "config file",
		},
	}
	app.Authors = []*cli.Author{{
		Name:  "Quinn T",
		Email: "quinn.github@gmail.com",
	}}
	app.Commands = []*cli.Command{
		{
			Name:      "random",
			UsageText: "wallpaper --config /path/to/config.json random",
			Usage:     "Changes the wallpaper to a random one",
			Before:    loadConfig,
			Action: func(c *cli.Context) error {
				return wallpapers.ChangeToRandom()
			},
		},
	}
	app.UseShortOptionHandling = true

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
