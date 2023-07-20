package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	constants "github.com/ItsMyEyes/install_kiyora/constant"
	"github.com/ItsMyEyes/install_kiyora/dto"
)

func ReleaseLatest() (res *dto.ResponseRelease) {
	// http client with header
	req, err := http.NewRequest("GET", "https://api.github.com/repos/ItsMyEyes/i-kiyora/releases/latest", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("User-Agent", "golang-http-client")
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return
	}

	return
}

func CheckForUpdate() (bool, string) {
	latest := ReleaseLatest()
	return latest.TagName != constants.Version, latest.TagName
}
