package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ItsMyEyes/install_kiyora/utils"
)

var (
	AZ = "https://github.com/ItsMyEyes/az"
)

func AddModule(dir, name string) error {
	switch name {
	case "az":
		return azModule(dir)
	}
	return errors.New("Module not found")
}

func GetNameModule(dir string) string {
	read := utils.ReadFile(dir + "\\go.mod")
	mod := string(read)
	mods := strings.Split(mod, "\n")
	mod = strings.Replace(mods[0], "module ", "", -1)
	return strings.Trim(mod, "\r\n")
}

func azModule(dir string) error {
	fmt.Println("📦 Installing AZ Module")
	CloningProject(AZ, dir+"\\az")

	return nil
}
