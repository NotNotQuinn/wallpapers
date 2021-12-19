package main

import (
	"errors"
	"fmt"
	"math/rand"
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
			Name: "random",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "playlist",
					Aliases: []string{"p"},
					Usage:   "Select a random wallpaper only from this playlist.",
				},
				&cli.BoolFlag{
					Name:    "untagged",
					Aliases: []string{"u"},
					Usage:   "Only show wallpapers not in any playlist. Ignored if -p used.",
				},
			},
			Usage: "Changes the wallpaper to a random one",
			Action: func(c *cli.Context) error {
				key := c.String("playlist")
				if key == "" {
					if c.Bool("untagged") {
						return wallpapers.SetRandomUntagged()
					} else {
						return wallpapers.SetRandom()
					}
				}
				playlists, err := wallpapers.LoadPlaylists()
				if err != nil {
					return fmt.Errorf("failed to load playlists: %w", err)
				}
				if playlists[key] == nil {
					return fmt.Errorf("no such playlist '%s' exists", key)
				}
				err = wallpapers.SetFromFile(playlists[key][rand.Intn(len(playlists[key]))])
				return err
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
			Name:      "add",
			Usage:     "Adds the current wallpaper to the named playlists.",
			ArgsUsage: "<playlists1> <playlist2>....",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return errors.New("please provide at least one playlist")
				}
				// create the playlist if new
				// and add the current wallpaper to it
				current, err := wallpapers.CurrentFilePath()
				if err != nil {
					return fmt.Errorf("failed to get current wallpaper: %w", err)
				}
				playlists, err := wallpapers.LoadPlaylists()
				if err != nil {
					return fmt.Errorf("failed to load playlists: %w", err)
				}
				keys := c.Args().Slice()
				for _, key := range keys {
					playlists[key] = append(playlists[key], current)
				}
				err = wallpapers.SavePlaylists()
				if err != nil {
					return fmt.Errorf("failed to save playlists: %w", err)
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "Lists all playlists.",
			Flags: []cli.Flag{&cli.BoolFlag{
				Name:     "long",
				Aliases:  []string{"l"},
				Usage:    "Output all playlists in full.",
				Required: false,
			}},
			Action: func(c *cli.Context) error {
				playlists, err := wallpapers.LoadPlaylists()
				if err != nil {
					return err
				}
				fmt.Println("Listing all playlists:")
				if len(playlists) == 0 {
					fmt.Println("<none>")
				}
				for i, list := range playlists {
					fmt.Printf("%s:  %d wallpapers\n", i, len(list))
					if c.Bool("long") {
						for _, v := range list {
							fmt.Printf("  %s\n", v)
						}
					}
				}
				return nil
			},
		},
		{
			Name:  "show",
			Usage: "Shows current wallpaper",
			Action: func(c *cli.Context) error {
				return wallpapers.PrintCurrentWallpaperInfo()
			},
		},
		{
			Name:      "remove",
			Aliases:   []string{"rem"},
			Usage:     "Removes current wallpaper from the named playlist.",
			ArgsUsage: "<playlist>",
			Action: func(c *cli.Context) error {
				keys := c.Args().Slice()
				playlists, err := wallpapers.LoadPlaylists()
				if err != nil {
					return fmt.Errorf("failed to load playlists: %w", err)
				}
				file, err := wallpapers.CurrentFilePath()
				if err != nil {
					return fmt.Errorf("failed to get current wallpaper: %w", err)
				}
				for _, key := range keys {
					var tmp = make([]string, 0)
					for _, v := range playlists[key] {
						if v != file {
							tmp = append(tmp, v)
						}
					}
					playlists[key] = tmp
				}

				return wallpapers.SavePlaylists()
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
