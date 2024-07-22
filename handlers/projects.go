package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	constants "github.com/ItsMyEyes/install_kiyora/constant"
	"github.com/ItsMyEyes/install_kiyora/dto"
	"github.com/ItsMyEyes/install_kiyora/utils"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/urfave/cli"
)

func CreateProject(_ *cli.Context) {
	var appCli dto.Cli
	Logo(constants.Kiyora, constants.Version)

	utils.MakeLine()

	// make input cli
	CheckVersionGolang(constants.MinimalVersionGolangInt, constants.MinimalVersionGolang)
	// check git
	CheckGit()

	w := wow.New(os.Stdout, spin.Get(spin.Grenade), " Checking Version")
	w.Start()
	checkUpdate, version := CheckForUpdate()
	if checkUpdate {
		w.PersistWith(spin.Spinner{Frames: []string{"😁"}}, fmt.Sprintf(" New version available: %s, Please update with command i-kiyora update", version))
	} else {
		w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, " You are using the latest version.")
	}

	GoPath, checkIt := utils.Getenv("GOPATH")
	if !checkIt {
		log.Fatal("You need to set GOPATH")
		os.Exit(0)
	}
	appCli.GoPath = GoPath

	utils.MakeLine()

	answer := utils.MakeInputCli("What is the name of the project?")
	if len(answer) == 0 {
		log.Fatal("You need to enter the project name")
		os.Exit(0)
	}
	appCli.NameProject = answer

	answer = utils.MakeInputCli("What name your github? ")
	if len(answer) == 0 {
		log.Fatal("You need to enter the github name")
		os.Exit(0)
	}

	appCli.GithubName = answer

	w = wow.New(os.Stdout, spin.Get(spin.Shark), " Creating Projects")
	// w.Start()
	utils.MakeLine()

	fmt.Println("Your project name is " + appCli.NameProject)
	fmt.Println("Your module is " + appCli.ModuleProject())
	fmt.Println("Your path is " + appCli.PathProject())

	utils.MakeLine()
	answer = utils.MakeInputCli("Are you sure this is the correct information? (y/n)")
	if !strings.Contains(strings.ToUpper(answer), "Y") {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " You have not created a new project.")
		os.Exit(0)
	}

	if !utils.CheckDir(appCli.PathProject()) {
		CloningProject(w, constants.ProjectLink, appCli.PathProject())

	} else {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Project already exists.")
		os.Exit(0)
	}

	// utils.RemoveFolder(fmt.Sprintf("%s%s.git", utils.GetPathSlash(), appCli.PathProject()))
	utils.RemoveFolder(utils.MakeDirectoryString(appCli.PathProject(), ".git"))

	// utils.CopyFile(fmt.Sprintf("%s%sconfig.yaml.example", utils.GetPathSlash(), appCli.PathProject()), fmt.Sprintf("%s%sapp.yaml", utils.GetPathSlash(), appCli.PathProject()))
	utils.CopyFile(utils.MakeDirectoryString(appCli.PathProject(), "config.yaml.example"), utils.MakeDirectoryString(appCli.PathProject(), "config.yaml"))

	utils.ReplaceTextInFolder(appCli.PathProject(), constants.Replace, appCli.ModuleProject())

	utils.ReplaceTextInFolder(appCli.PathProject(), "services_name", appCli.NameModule())

	RunnigMod(appCli.PathProject(), appCli.ModuleProject())

	RunningTidy(appCli.PathProject())
	// w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, " Project created successfully.")

	utils.MakeLine()

	fmt.Println("Your project is ready to use.")
	fmt.Println("You can run it with the following commands:")
	fmt.Println("$ cd " + appCli.PathProject())
	fmt.Println("$ go run ./cmd/http-server/main.go")

	utils.MakeLine()
}

func AddModular(ctx *cli.Context) {
	args := ctx.Args()

	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Creating Projects")
	// w.Start()
	if len(args) == 0 {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Please enter the module name.")
		os.Exit(0)
	}

	moduleName := args[0]

	dir, err := os.Getwd()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, fmt.Sprintf("Error getting current directory: %s\n", err.Error()))
		return
	}
	// adapter := dir + utils.GetPathSlash() + "adapter"
	adapter := utils.MakeDirectoryString(dir, "adapter")
	projectName := GetNameModule(dir)
	// check directory cli running
	if !utils.CheckDir(adapter) {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, "Project not found. Please run this command in the project directory.")
		os.Exit(0)
	}
	w.Text("Creating module " + moduleName + " ...").Spinner(spin.Get(spin.Shark))

	// check directory adapter
	if utils.CheckDir(utils.MakeDirectoryString(adapter, moduleName)) {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, "Module already exists.")
		os.Exit(0)
	}

	AddModule(dir, moduleName)

	utils.ReplaceTextInFolder(adapter, constants.Replace, projectName)
	// w.PersistWith()

	RunningTidy(dir)

	utils.MakeLine()

	w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, " Module created successfully.")

	utils.MakeLine()
}
