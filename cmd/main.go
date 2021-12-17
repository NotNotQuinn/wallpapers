package main

import (
	"fmt"
	"os"

	"github.com/notnotquinn/wallpapers/wallpapers"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "wallpaper"
	app.Version = "0.1.1"
	app.Usage = "Randomize your wallpapers"
	app.Commands = []*cli.Command{
		{
			Name:  "random",
			Usage: "Changes the wallpaper to a random one",
			Action: func(c *cli.Context) error {
				return wallpapers.SetRandom()
			},
		},
		{
			Name:  "downloadone",
			Usage: "Downloads a random wallpaper.",
			Action: func(c *cli.Context) error {
				url, catagory, err := wallpapers.NewRandomURL(nil, nil)
				if err != nil {
					return err
				}
				path, _ := wallpapers.CalculatePath(url)
				fmt.Printf("Downloading: %s\n  %s\n  to %s\n", url, catagory, path)
				_, err = wallpapers.AddUrl(url)
				if err != nil {
					return err
				}
				fmt.Println("done")
				return nil
			},
		},
		{
			Name:  "add",
			Usage: "Adds the current wallpaper to the named playlist.",
			Action: func(c *cli.Context) error {
				// create the playlist if new
				// and add the current wallpaper to it
				current, err := wallpapers.CurrentFilePath()
				if err != nil {
					return err
				}
				playlists, err := wallpapers.LoadPlaylists()
				if err != nil {
					return err
				}
				(*playlists)[c.Args().First()] = append((*playlists)[c.Args().First()], current)
				return wallpapers.SavePlaylists()
			},
		},
		{
			Name:  "list",
			Usage: "Lists all playlists.",
			Action: func(c *cli.Context) error {
				playlists, err := wallpapers.LoadPlaylists()
				if err != nil {
					return err
				}
				fmt.Println("Listing all playlists:")
				if len(*playlists) == 0 {
					fmt.Println("<none>")
				}
				for i, list := range *playlists {
					fmt.Printf("%s:\n", i)
					for _, v := range list {
						fmt.Printf("  %s\n", v)
					}
				}
				return nil
			},
		},
	}
	app.UseShortOptionHandling = true

	// Some things that need to get done:
	// - Add ourselves to start on boot. (cross platform may be difficult)
	// - Some configuration options (i guess??)
	// - NSFW detection and block toggle. (not needed if users can curate their own lists of images)
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
