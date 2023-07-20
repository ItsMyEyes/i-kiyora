package handlers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	github_pkg "github.com/ItsMyEyes/install_kiyora/pkg/github"
	"github.com/ItsMyEyes/install_kiyora/utils"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/urfave/cli"
)

func CheckPath(source string, custom bool) (string, error) {
	if custom && source == "" {
		return "", errors.New("PATH not found, are you custom PATH? you must run with --source")
	}
	checkDir := source
	if source == "" {
		checkDir = "C:\\i-kiyora"
	}
	if _, err := os.Stat(checkDir); os.IsNotExist(err) {
		return "", errors.New("PATH not found, are you custom PATH? you must run with --source")
	}

	if !custom {
		split := os.Getenv("PATH")
		if !strings.Contains(split, fmt.Sprintf("%s%s", checkDir, utils.GetPathSlash())) {
			return "", errors.New("PATH not found, are you custom PATH? you must run with --source")
		}
		return fmt.Sprintf("%s%s", checkDir, utils.GetPathSlash()), nil
	}

	return fmt.Sprintf("%s%s", source, utils.GetPathSlash()), nil
}

func UpdateBinary(d *cli.Context) {
	w := wow.New(os.Stdout, spin.Get(spin.Grenade), " Checking Version")
	w.Start()
	getValueFlag := d.String("path")

	requiredUpdate, version := CheckForUpdate()
	if !requiredUpdate {
		w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, fmt.Sprintf(" You are using the latest version. %s", version))
		os.Exit(0)
	}

	w.Text(fmt.Sprintf("😁 New version available: %s, Updating...", version))

	path, err := CheckPath(getValueFlag, d.IsSet("path"))
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Path %s, %s", path, err.Error()))
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, fmt.Sprintf(" Path %s", path))
	Latest, err := github_pkg.GetLatest()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Error %s", err.Error()))
		os.Exit(0)
	}

	defer func() {
		if r := recover(); r != nil {
			w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Error %s", r))
			os.Exit(0)
		}

		w = wow.New(os.Stdout, spin.Get(spin.Grenade), " Cleaning")
		w.Start()
		err = Latest.DeleteArchinve()
		if err != nil {
			w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Error %s", err.Error()))
			os.Exit(0)
		}
		w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, fmt.Sprintf(" Cleaned and update to %s", *Latest.RepoRelease.TagName))
	}()

	w = wow.New(os.Stdout, spin.Get(spin.Grenade), " Downloading")
	w.Start()
	err = Latest.DownloadAndCopy()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Error %s", err.Error()))
		os.Exit(0)
	}
	w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, fmt.Sprintf(" Downloaded"))

	w = wow.New(os.Stdout, spin.Get(spin.Grenade), " Extracting")
	w.Start()
	err = Latest.Extract()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Error %s", err.Error()))
		os.Exit(0)
	}
	w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, fmt.Sprintf(" Extracted"))

	w = wow.New(os.Stdout, spin.Get(spin.Grenade), " Move")
	w.Start()
	err = Latest.Move(path)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌ "}}, fmt.Sprintf(" Error %s", err.Error()))
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✅ "}}, fmt.Sprintf(" Moved"))
}
