package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	rootFolder := "C:\\Users\\oshie\\Documents\\andi\\golang\\balance_services"
	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Replace text in the file
		fmt.Println(path)
		if err != nil {
			fmt.Printf("Error replacing text in file %s: %s\n", path, err.Error())
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %s\n", rootFolder, err.Error())
	}
}
