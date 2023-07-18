package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	oss := runtime.GOOS
	arch := runtime.GOARCH
	// Set your GitHub personal access token here
	token := "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"

	// Specify the repository details
	owner := "ItsMyEyes"
	repo := "i-kiyora"

	checkDir := "C:\\i-kiyora"
	if _, err := os.Stat(checkDir); os.IsNotExist(err) {
		log.Fatal("PATH not found, are you custom PATH? you must run with --source")
	}

	split := os.Getenv("PATH")
	if !strings.Contains(split, fmt.Sprintf("%s\\", checkDir)) {
		log.Fatal("PATH not found, are you custom PATH? you must run with --source")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		fmt.Printf("Error retrieving release: %v\n", err)
		return
	}

	// Specify the name of the binary file in the release assets
	assetName := fmt.Sprintf("i-kiyora_%s_%s_%s.tar.gz", strings.Replace(*release.TagName, "v", "", -1), oss, arch)

	// Find the asset by name
	var asset *github.ReleaseAsset
	for _, a := range release.Assets {
		fmt.Println(a.GetName())
		if strings.EqualFold(a.GetName(), assetName) {
			asset = &a
			break
		}
	}

	if asset == nil {
		fmt.Printf("Asset '%s' not found in the release.\n", assetName)
		return
	}

	// // Download the asset file
	resp, err := http.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		fmt.Printf("Error downloading asset: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// // Save the asset file locally
	file, err := os.Create(assetName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return
	}

	// // Extract the tar.gz archive
	err = extractTarGz(assetName, "./")
	if err != nil {
		fmt.Printf("Error extracting archive: %v\n", err)
		return
	}
	// // Make the binary executable
	err = os.Chmod(assetName, 0755)
	if err != nil {
		fmt.Printf("Error making binary executable: %v\n", err)
		return
	}

	// move file
	err = os.Rename("./i-kiyora.exe", checkDir+"\\i-kiyora.exe")
	if err != nil {
		fmt.Printf("Error moving file: %v\n", err)
		return
	}
}

func extractTarGz(archivePath string, destinationDir string) error {
	// Open the tar.gz archive file
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// Create a gzip reader to read the compressed data
	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader to read the tar archive
	tarReader := tar.NewReader(gzipReader)

	// Extract each file in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// Reached end of tar archive
			break
		}
		if err != nil {
			return err
		}

		// Determine the destination path for the extracted file
		destPath := filepath.Join(destinationDir, header.Name)

		// Check if the file is a directory
		if header.FileInfo().IsDir() {
			// Create the directory in the destination path
			err := os.MkdirAll(destPath, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			continue
		}

		// Create the parent directory for the file
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			return err
		}

		// Create the destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy the file data from the tar archive to the destination file
		_, err = io.Copy(destFile, tarReader)
		if err != nil {
			return err
		}
	}

	// // Delete the archive file
	// err = os.Remove(archivePath)
	// if err != nil {
	// 	return err
	// }

	return nil
}
