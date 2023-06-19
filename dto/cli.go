package dto

import (
	"fmt"
	"strings"
)

type Cli struct {
	GoPath        string
	NameProject   string
	RootDirectory string
	GithubName    string
}

func (c *Cli) NameModule() string {
	result := strings.ToLower(c.NameProject)
	result = strings.Replace(result, " ", "_", -1)
	return result
}

func (c *Cli) ModuleProject() string {
	return fmt.Sprintf("github.com/%s/%s", c.GithubName, c.NameModule())
}

func (c *Cli) PathProject() string {
	return fmt.Sprintf("%s\\%s", c.GoPath, c.NameModule())
}
