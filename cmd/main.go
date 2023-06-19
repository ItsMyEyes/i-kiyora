package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ItsMyEyes/install_kiyora/dto"
	"github.com/ItsMyEyes/install_kiyora/handlers"
	"github.com/ItsMyEyes/install_kiyora/utils"
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
)

func main() {
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

	fmt.Println("Creating a new project...")
	if !utils.CheckDir(appCli.PathProject()) {
		fmt.Println("Cloning project... ", appCli.PathProject())
		handlers.CloningProject(ProjectLink, appCli.PathProject())
	} else {
		fmt.Println("Project already exists")
		os.Exit(0)
	}

	fmt.Println("Removing .git folder...")
	utils.RemoveFolder(fmt.Sprintf("%s\\.git", appCli.PathProject()))

	fmt.Println("Copying file...")
	utils.CopyFile(fmt.Sprintf("%s\\app.yaml.example", appCli.PathProject()), fmt.Sprintf("%s\\app.yaml", appCli.PathProject()))

	err := utils.ReplaceTextInFolder(appCli.PathProject(), Replace, appCli.ModuleProject())
	if err != nil {
		fmt.Printf("Error replacing text in folder: %s\n", err.Error())
		return
	}

	fmt.Println("Text replacement completed.")

	fmt.Println("Running mod...")
	handlers.RunnigMod(appCli.PathProject(), appCli.ModuleProject())

	fmt.Println("Running tidy...")
	handlers.RunningTidy(appCli.PathProject())

	utils.MakeLine()

	fmt.Println("Your project is ready to use.")
	fmt.Println("You can run it with the following commands:\n")
	fmt.Println("$ cd " + appCli.PathProject())
	fmt.Println("$ go run ./cmd/http-server/main.go")

	utils.MakeLine()
}
