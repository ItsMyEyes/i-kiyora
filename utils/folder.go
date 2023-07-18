package utils

import (
	"io"
	"log"
	"os"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func RemoveFolder(dir string) {
	cmd := os.RemoveAll(dir)
	if cmd != nil {
		log.Fatal(cmd)
	}
}

func CopyFile(src, dst string) {
	w := wow.New(os.Stdout, spin.Get(spin.Shark), " Copying File")
	w.Start()
	in, err := os.Open(src)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Cant copy file "+err.Error())
		os.Exit(0)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Cant copy file "+err.Error())
		os.Exit(0)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"❌"}}, " Cant copy file "+err.Error())
		os.Exit(0)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, " File copied")
}

func ReadFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		log.Fatal(err)
	}

	return bs
}

func CheckDir(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func Getenv(env string) (string, bool) {
	if value, ok := os.LookupEnv(env); ok {
		return value, true
	}
	return "", false
}
