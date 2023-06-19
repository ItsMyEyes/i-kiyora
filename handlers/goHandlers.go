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
)

func RunnigMod(dir string, mod string) {
	cmd := exec.Command("go", "mod", "init", mod)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Sucess running Mod" + out.String())
}

func RunningTidy(dir string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n ✅ Sucess running tidy" + out.String())
}

func CloningProject(project string, nameProject string) {
	cmd := exec.Command("git", "clone", project, nameProject)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n ✅ Sucess cloning" + out.String())
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
