package constants

import (
	"strconv"
	"strings"
)

var (
	Answer                     string
	Kiyora                     = "Kiyora"
	Version                    = "v0.3.4"
	MinimalVersionGolang       = "1.22.0"
	ProjectLink                = "https://github.com/ItsMyEyes/kiyora_v3.git"
	Replace                    = "github.com/ItsMyEyes/kiyora_v3"
	MinimalVersionGolangInt, _ = strconv.Atoi(strings.Replace(MinimalVersionGolang, ".", "", -1))
	BuildDate                  = "2024-07-22"
	Commit                     = "now"
)
