package handlers

import (
	constants "github.com/ItsMyEyes/install_kiyora/constant"
	github_pkg "github.com/ItsMyEyes/install_kiyora/pkg/github"
)

func CheckForUpdate() (bool, string) {
	latest := github_pkg.ReleaseLatest()
	return latest.GetTagName() != constants.Version, latest.GetTagName()
}
