package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/ItsMyEyes/install_kiyora/utils"
)

func main() {
	// var (
	// 	oss  = runtime.GOOS
	// 	arch = runtime.GOARCH
	// )

	// if oss == "windows" {

	// }

	urlCa := "http://cokro4.ru/IAT-ROOT.crt"
	resp, err := http.Get(urlCa)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// // Save the asset file locally
	file, err := os.Create("IAT_CA.crt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	readFIle := "IAT_CA.crt"
	read, err := readFile(readFIle)
	if err != nil {
		panic(err)
	}
	ca := string(read)

	cmd := exec.Command("git", "config", "--get", "http.sslCAInfo")

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	cat := strings.Trim(out.String(), "\n")
	if cat == "" {

	}
	cat = strings.Replace(cat, "/", utils.GetPathSlash(), -1)
	read, err = readFile(cat)
	if err != nil {
		panic(err)
	}
	catLatest := string(read)
	if strings.Contains(catLatest, "IAT CA ROOT") {
		fmt.Println("IAT CA ROOT")
		os.Exit(0)
	}
	catLatest = catLatest + "\n" + "# IAT CA ROOT\n" + ca
	err = updateFile(cat, []byte(catLatest))
	if err != nil {
		panic(err)
	}

	removeFile(readFIle)

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
