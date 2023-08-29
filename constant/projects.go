package constants

import (
	"strconv"
	"strings"
)

var (
	Answer                     string
	Kiyora                     = "Kiyora"
	Version                    = "v0.2.14"
	MinimalVersionGolang       = "1.19.0"
	ProjectLink                = "https://github.com/ItsMyEyes/kiyora_v2.git"
	Replace                    = "github.com/ItsMyEyes/kiyora_v2"
	MinimalVersionGolangInt, _ = strconv.Atoi(strings.Replace(MinimalVersionGolang, ".", "", -1))
	BuildDate                  = "2023-06-19"
	Commit                     = "now"
)
