package main

import (
	"fmt"
	"os"
	"time"

	constants "github.com/ItsMyEyes/install_kiyora/constant"
	"github.com/ItsMyEyes/install_kiyora/handlers"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ikiyora"
	app.Usage = "run ikiyora"
	app.Author = "ItsMyEyes - Andi"
	app.UsageText = "ikiyora [command] [arguments...]"
	app.Version = fmt.Sprintf("%s built on %s (commit: %s)", constants.Version, constants.BuildDate, constants.Commit)
	app.Description = "IKiyora is a installation tool for Kiyora"
	app.Commands = []cli.Command{
		{
			Name:        "create",
			Description: "Creating a new project",
			UsageText:   "Creating a new project",
			Usage:       "Creating a new project",
			Action:      handlers.CreateProject,
		},
		{
			Name:        "add",
			Description: "Add Module / Adapters",
			UsageText:   "Add Module / Adapters",
			Usage:       "Add Module / Adapters",
			Action:      handlers.AddModular,
		},
		{
			Name:        "loading",
			Description: "Add Module / Adapters",
			UsageText:   "Add Module / Adapters",
			Usage:       "Add Module / Adapters",
			Action:      testLoading,
		},
		{
			Name:        "update",
			Description: "Update ikiyora",
			UsageText:   "Update ikiyora",
			Usage:       "Update ikiyora",
			Action:      handlers.UpdateBinary,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "path",
					Usage:  "path of ikiyora",
					EnvVar: "IKIYORA_PATH",
				},
			},
		},
	}

	app.Run(os.Args)
}

func testLoading(_ *cli.Context) {
	w := wow.New(os.Stdout, spin.Get(spin.Runner), " Such Spins")
	w.Start()
	time.Sleep(2 * time.Second)
	w.Text("Very emojis").Spinner(spin.Get(spin.Shark))
	time.Sleep(2 * time.Second)
	w.PersistWith(spin.Spinner{Frames: []string{"👍"}}, " Wow!")

	z := wow.New(os.Stdout, spin.Get(spin.Runner), " Such Spins")
	z.Start()
	time.Sleep(2 * time.Second)
	z.Text("Very emojis").Spinner(spin.Get(spin.Shark))
	time.Sleep(2 * time.Second)
	z.PersistWith(spin.Spinner{Frames: []string{"👍"}}, " Wow!")
}
