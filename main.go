package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:    "Wallpaper Updater",
		Version: "0.0.1-dev",
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	// err := wp.SetMode(wp.Fit)
	// if err != nil {
	// 	panic(err)
	// }

	// err = wp.SetFromURL("https://cdn.betterttv.net/emote/54fbf00a01abde735115de5c/3x")
	// if err != nil {
	// 	panic(err)
	// }
	// path, err := wp.Get()
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = fmt.Println(path)
	// if err != nil {
	// 	panic(err)
	// }
}
