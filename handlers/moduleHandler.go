package handlers

import (
	"os"
	"strings"

	"github.com/ItsMyEyes/install_kiyora/utils"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

var (
	AZ = "https://github.com/ItsMyEyes/az"
)

func AddModule(dir, name string) {
	switch name {
	case "az":
		azModule(dir)
		return
	}
}

func GetNameModule(dir string) string {
	read := utils.ReadFile(utils.MakeDirectoryString(dir, "go.mod"))
	mod := string(read)
	mods := strings.Split(mod, "\n")
	mod = strings.Replace(mods[0], "module ", "", -1)
	return strings.Trim(mod, "\r\n")
}

func azModule(dir string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Adding Module AZ")
	if !utils.CheckDir(utils.MakeDirectoryString(dir, "pkg", "logger")) {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Module az cant install because logger")
		return
	}

	if !utils.CheckDir(utils.MakeDirectoryString(dir, "pkg", "util")) {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Module az cant install because util")
		return
	}

	if !utils.CheckDir(utils.MakeDirectoryString(dir, "adapter", "az")) {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Module az cant install because adapter")
		return
	}

	CloningProject(w, AZ, utils.MakeDirectoryString(dir, "adapter", "az"))
}
