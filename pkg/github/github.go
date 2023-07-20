package github_pkg

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	OSS = runtime.GOOS
	Arc = runtime.GOARCH
	// Set your GitHub personal access token here
	token = "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"

	// Specify the repository details
	owner = "ItsMyEyes"
	repo  = "i-kiyora"
	ctx   = context.Background()
	ts    = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc     = oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
)

type Result struct {
	RepoRelease *github.RepositoryRelease
}

func GetLatest() (*Result, error) {
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return &Result{
		RepoRelease: release,
	}, nil
}

func (release *Result) GetAssetsName() string {
	assetName := fmt.Sprintf("i-kiyora_%s_%s_%s.tar.gz", strings.Replace(*release.RepoRelease.TagName, "v", "", -1), OSS, Arc)
	return assetName
}

func (release *Result) GetAssets() *github.ReleaseAsset {
	var asset *github.ReleaseAsset
	for _, a := range release.RepoRelease.Assets {
		if strings.EqualFold(a.GetName(), release.GetAssetsName()) {
			asset = &a
			break
		}
	}
	return asset
}

func (release *Result) DownloadAndCopy() error {
	asset := release.GetAssets()
	resp, err := http.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// // Save the asset file locally
	file, err := os.Create(release.GetAssetsName())
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

func (release *Result) Extract() error {
	// // Extract the tar.gz archive
	err := extractTarGz(release.GetAssetsName(), "./")
	if err != nil {
		return fmt.Errorf("error extracting archive: %v\n", err)
	}
	// // Make the binary executable
	err = os.Chmod(release.GetAssetsName(), 0755)
	if err != nil {
		return fmt.Errorf("error making binary executable: %v\n", err)
	}

	return nil
}

func (release *Result) Move(source string) error {
	err := os.Rename("./i-kiyora.exe", source+"\\i-kiyora.exe")
	if err != nil {
		return fmt.Errorf("error moving file: %v\n", err)
	}

	return nil
}

func (release *Result) DeleteArchinve() error {
	err := os.Remove(release.GetAssetsName())
	if err != nil {
		return fmt.Errorf("error deleting archive: %v\n", err)
	}

	return nil
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
