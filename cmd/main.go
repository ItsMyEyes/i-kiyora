package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ItsMyEyes/install_kiyora/dto"
	"github.com/ItsMyEyes/install_kiyora/handlers"
	"github.com/ItsMyEyes/install_kiyora/utils"
	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
)

var (
	// Version is the current version of the app
	answer                     string
	Kiyora                     = "Kiyora"
	Version                    = "1.0.0"
	MinimalVersionGolang       = "1.19.0"
	ProjectLink                = "https://github.com/ItsMyEyes/kiyora_v2.git"
	Replace                    = "github.com/ItsMyEyes/kiyora_v2"
	MinimalVersionGolangInt, _ = strconv.Atoi(strings.Replace(MinimalVersionGolang, ".", "", -1))
	BuildDate                  = "2023-06-19"
	Commit                     = "now"

	s = spinner.New(spinner.CharSets[11], 300*time.Millisecond)
)

func main() {
	app := cli.NewApp()
	app.Name = "ikiyora"
	app.Usage = "run ikiyora"
	app.Author = "ItsMyEyes - Andi"
	app.UsageText = "ikiyora [global options] command [command options] [arguments...]"
	app.Version = fmt.Sprintf("%s built on %s (commit: %s)", Version, BuildDate, Commit)
	app.Description = "IKiyora is a installation tool for Kiyora"
	app.Commands = []cli.Command{
		{
			Name:        "create",
			Description: "Creating a new project",
			Action:      createProject,
		},
		{
			Name:        "add",
			Description: "Add Module / Adapter",
			Action:      addModule,
		},
	}

	app.Run(os.Args)
}

func addModule(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) == 0 {
		fmt.Println("You have not entered a module name.")
		os.Exit(0)
	}

	moduleName := args[0]

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %s\n", err.Error())
		return
	}
	adapter := dir + "\\adapter"
	projectName := handlers.GetNameModule(dir)
	// check directory cli running
	if !utils.CheckDir(adapter) {
		fmt.Println("❌ Project not found. Please run this command in the project directory.")
		os.Exit(0)
	}

	// check directory adapter
	if utils.CheckDir(adapter + "\\" + moduleName) {
		fmt.Println("❌ Adapter already exists.")
		os.Exit(0)
	}

	s.Suffix = " Creating a new adapter..."
	s.Start()
	err = handlers.AddModule(adapter, moduleName)
	if err != nil {
		s.Stop()
		fmt.Printf("Error adding module: %s\n", err.Error())
		return
	}
	s.Stop()

	err = utils.ReplaceTextInFolder(adapter, Replace, projectName)
	if err != nil {
		fmt.Printf("Error replacing text in folder: %s\n", err.Error())
		return
	}

	fmt.Println("Text replacement completed.")

	utils.MakeLine()

	fmt.Println("Your project is ready to use.")
	fmt.Println("✅ You can check it")

	utils.MakeLine()
}

func createProject(_ *cli.Context) {
	var appCli dto.Cli
	handlers.Logo(Kiyora, Version)

	utils.MakeLine()

	// make input cli
	handlers.CheckVersionGolang(MinimalVersionGolangInt, MinimalVersionGolang)
	// check git
	handlers.CheckGit()

	GoPath, checkIt := utils.Getenv("GOPATH")
	if !checkIt {
		fmt.Println("You need to set GOPATH")
		os.Exit(0)
	}
	appCli.GoPath = GoPath
	fmt.Println("Your GOPATH is " + appCli.GoPath)

	utils.MakeLine()

	answer = utils.MakeInputCli("What is the name of the project?")
	if len(answer) == 0 {
		fmt.Println("You have not created a new project.")
		os.Exit(0)
	}
	appCli.NameProject = answer

	answer = utils.MakeInputCli("What name your github? ")
	if len(answer) == 0 {
		fmt.Println("You have not created a new project.")
		os.Exit(0)
	}

	appCli.GithubName = answer

	utils.MakeLine()

	fmt.Println("Your project name is " + appCli.NameProject)
	fmt.Println("Your module is " + appCli.ModuleProject())
	fmt.Println("Your path is " + appCli.PathProject())

	utils.MakeLine()
	answer = utils.MakeInputCli("Are you sure this is the correct information? (y/n)")
	if !strings.Contains(strings.ToUpper(answer), "Y") {
		fmt.Println("You have not created a new project.")
		os.Exit(0)
	}

	s.Start()
	s.Suffix = " Creating a new project..."
	if !utils.CheckDir(appCli.PathProject()) {
		fmt.Println("Cloning project... ", appCli.PathProject())
		handlers.CloningProject(ProjectLink, appCli.PathProject())
		s.Stop()
	} else {
		fmt.Println("Project already exists")
		os.Exit(0)
	}

	s.Start()
	s.Suffix = " Removing .git folder..."
	utils.RemoveFolder(fmt.Sprintf("%s\\.git", appCli.PathProject()))
	s.Stop()

	s.Start()
	s.Suffix = " Copying file..."
	utils.CopyFile(fmt.Sprintf("%s\\app.yaml.example", appCli.PathProject()), fmt.Sprintf("%s\\app.yaml", appCli.PathProject()))
	s.Stop()

	err := utils.ReplaceTextInFolder(appCli.PathProject(), Replace, appCli.ModuleProject())
	if err != nil {
		fmt.Printf("Error replacing text in folder: %s\n", err.Error())
		return
	}

	fmt.Println("Text replacement completed.")

	err = utils.ReplaceTextInFolder(appCli.PathProject(), "services_name", appCli.NameModule())
	if err != nil {
		fmt.Printf("Error replacing text in folder: %s\n", err.Error())
		return
	}

	fmt.Println("Text replacement completed.")

	handlers.RunnigMod(appCli.PathProject(), appCli.ModuleProject())

	s.Start()
	s.Suffix = " Running tidy..."
	handlers.RunningTidy(appCli.PathProject())
	s.Stop()
	utils.MakeLine()

	fmt.Println("Your project is ready to use.")
	fmt.Println("You can run it with the following commands:\n")
	fmt.Println("$ cd " + appCli.PathProject())
	fmt.Println("$ go run ./cmd/http-server/main.go")

	utils.MakeLine()
}
