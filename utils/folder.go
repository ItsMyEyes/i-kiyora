package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func RemoveFolder(dir string) {
	cmd := os.RemoveAll(dir)
	if cmd != nil {
		log.Fatal(cmd)
	}

	fmt.Println("Sucess remove folder " + dir)
}

func CopyFile(src, dst string) error {
	fmt.Println(src, dst)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
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
