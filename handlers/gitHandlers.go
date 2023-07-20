package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	github_pkg "github.com/ItsMyEyes/install_kiyora/pkg/github"
	"github.com/ItsMyEyes/install_kiyora/utils"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/urfave/cli"
)

func GitSSL(d *cli.Context) {
	oss := github_pkg.OSS

	getUrlCA := d.String("url-ca")
	configPath := d.String("config-path")
	checkPermission(oss)

	if oss == "windows" {
		GitSSLIATwindows(getUrlCA, configPath)
	}

	if oss == "linux" {
		GitSSLIATLinux(getUrlCA, configPath)
	}

	if oss == "darwin" {
		GitSSLIATDarwin(getUrlCA, configPath)
	}

	fmt.Println("Git SSL not support for " + oss)
}

func checkPermission(oss string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Checking Permission")
	w.Start()
	if oss == "windows" {
		cmd := exec.Command("net", "session")

		// Attempt to run the command
		err := cmd.Run()

		// Check the error returned by the command
		if err != nil {
			w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL need Administrator permission, please run as Administrator")
			os.Exit(0)
		}
	}

	if oss == "linux" || oss == "darwin" {
		cmd := exec.Command("sudo", "echo", "test")
		err := cmd.Run()
		if err != nil {
			w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL need Administrator permission, please run with sudo")
			os.Exit(0)
		}
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✔"}}, " Git SSL have permission")
}

func GitSSLIATwindows(urlCa string, pathSSL string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Adding Git SSL at windows")
	w.Start()
	if urlCa == "" {
		urlCa = "http://cokro4.ru/IAT-ROOT.crt"
	}

	configSsl, err := getConfigSSL()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant get config "+err.Error())
		os.Exit(0)
	}

	if configSsl == "" {
		configSsl = pathSSL
	}

	ca := downloadConfig(w, urlCa)

	configSsl = strings.Replace(configSsl, "/", utils.GetPathSlash(), -1)
	updateGitSSL(w, configSsl, ca)

	w.PersistWith(spin.Spinner{Frames: []string{"✔"}}, " Git SSL installed")
	os.Exit(0)
}

func GitSSLIATLinux(urlCa string, pathSSL string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Adding Git SSL at "+github_pkg.OSS)
	w.Start()
	if urlCa == "" {
		urlCa = "http://cokro4.ru/IAT-ROOT.crt"
	}
	readFIle := "IAT_CA.crt"
	err := downloadFile(urlCa, readFIle)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant download "+err.Error())
		os.Exit(0)
	}

	if pathSSL == "" {
		pathSSL = "/usr/local/share/ca-certificates/"
	}

	err = copyFile(readFIle, pathSSL)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant copy "+err.Error())
		os.Exit(0)
	}

	cmd := exec.Command("sudo", "update-ca-certificates")
	err = cmd.Run()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant update "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✔"}}, " Git SSL installed")
	os.Exit(0)
}

func GitSSLIATDarwin(urlCa string, pathSSL string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Adding Git SSL at "+github_pkg.OSS)
	w.Start()
	if urlCa == "" {
		urlCa = "http://cokro4.ru/IAT-ROOT.crt"
	}
	readFIle := "IAT_CA.crt"
	err := downloadFile(urlCa, readFIle)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant download "+err.Error())
		os.Exit(0)
	}

	if pathSSL == "" {
		pathSSL = "/Library/Keychains/System.keychain"
	}

	cmd := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", pathSSL, readFIle)
	err = cmd.Run()
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant update "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✔"}}, " Git SSL installed")
	os.Exit(0)
}

func downloadConfig(w *wow.Wow, urlCa string) string {
	readFIle := "IAT_CA.crt"
	err := downloadFile(urlCa, readFIle)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant download "+err.Error())
		os.Exit(0)
	}

	read, err := readFile(readFIle)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant read file "+err.Error())
		os.Exit(0)
	}
	ca := string(read)
	return ca
}

func getLatestSSLCrt(w *wow.Wow, configSsl string) string {
	read, err := readFile(configSsl)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant read "+err.Error())
		os.Exit(0)
	}
	catLatest := string(read)

	if strings.Contains(catLatest, "IAT CA ROOT") {
		w.PersistWith(spin.Spinner{Frames: []string{"✔"}}, " Git SSL already installed")
		os.Exit(0)
	}
	return catLatest
}

func updateGitSSL(w *wow.Wow, configSsl string, ca string) {
	catLatest := getLatestSSLCrt(w, configSsl)

	catLatest = catLatest + "\n" + "# IAT CA ROOT\n" + ca
	err := updateFile(configSsl, []byte(catLatest))
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Git SSL cant update "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✔"}}, " Git SSL installed")
	os.Exit(0)
}

func getConfigSSL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "http.sslCAInfo")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.Trim(out.String(), "\n"), nil
}

func downloadFile(url string, dest string) error {
	if url == "" {
		return fmt.Errorf("url is empty")
	}
	if dest == "" {
		return fmt.Errorf("dest is empty")
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// // Save the asset file locally
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return io.ReadAll(f)
}

func updateFile(name string, data []byte) error {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func removeFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(source string, dest string) error {
	input, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
