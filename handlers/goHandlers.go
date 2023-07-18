package handlers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ItsMyEyes/install_kiyora/utils"
	"github.com/common-nighthawk/go-figure"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func RunnigMod(dir string, mod string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Create Modfd")
	w.Start()
	cmd := exec.Command("go", "mod", "init", mod)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Cant init mod "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, " Mod init")
}

func RunningTidy(dir string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Running Tidy")
	w.Start()
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Cant running tidy "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, " Tidy")
}

func CloningProject(w *wow.Wow, project string, nameProject string) {
	w.Text("Cloning project").Spinner(spin.Get(spin.Shark))
	cmd := exec.Command("git", "clone", project, nameProject)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Cant clone project "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, " Project cloned")
}

func Logo(Name string, Version string) {
	myFigure := figure.NewFigure(Name, "", true)
	myFigure.Print()
	utils.NCenter(70, "Version "+Version).WriteTo(os.Stdout)
}

func CheckGit() {
	cmd := exec.Command("git", "--version")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("You are using git version " + out.String())
}

func CheckVersionGolang(MinimalVersionGolangInt int, MinimalVersionGolang string) {
	cmd := exec.Command("go", "version")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		fmt.Print(err)
	}
	// get version interger
	version := strings.Split(out.String(), " ")[2]
	version = strings.Replace(version, "go", "", 1)
	versionInt, _ := strconv.Atoi(strings.Replace(version, ".", "", -1))
	if versionInt < MinimalVersionGolangInt {
		fmt.Print("You need to install go version " + MinimalVersionGolang + " or higher")
	}

	fmt.Println("You are using go version " + version)
}
