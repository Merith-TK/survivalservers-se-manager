package main

import (
	"fmt"
	"os"

	"github.com/secsy/goftp"
)

func isDirectory(client *goftp.Client, remotePath string) (bool, error) {
	entries, err := client.ReadDir(remotePath)
	if err != nil {
		return false, err
	}

	// If there are any entries, assume it's a directory
	return len(entries) > 0, nil
}

func pullData(client *goftp.Client) {
	fmt.Println("Pulling data from FTP server")
	// for each entry in conf.Files, key is remote, value is local
	for _, file := range conf.Files {
		fmt.Println("Checking", file.Remote)
		// check if file is a folder or a file by checking remote
		remoteFile, fileErr := client.Stat(file.Remote)
		if fileErr != nil {
			panic(fileErr)
		}
		if remoteFile == nil {
			continue
		}
		if remoteFile.IsDir() {
			fetchFolder(client, file.Remote, file.Local)
		} else {
			fetchFile(client, file.Remote, file.Local)
		}
	}
}
func fetchFolder(client *goftp.Client, remote string, local string) {
	fmt.Println("Fetching folder", remote, "to", local)
	// Create the folder
	err := os.MkdirAll(local, 0755)
	if err != nil {
		panic(err)
	}
	entries, err := client.ReadDir(remote)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		// catch if folder is . or .. and skip
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}
		if entry.IsDir() {
			fetchFolder(client, remote+"/"+entry.Name(), local+"/"+entry.Name())
		} else {
			if same, _ := areFilesSame(local+"/"+entry.Name(), remote+"/"+entry.Name(), client); same {
				fmt.Println("Skipping", remote+"/"+entry.Name())
				continue
			}
			fetchFile(client, remote+"/"+entry.Name(), local+"/"+entry.Name())
		}
	}
}
func fetchFile(client *goftp.Client, remote string, local string) {
	fmt.Println("Fetching", remote, "to", local)
	// Create the file
	outFile, err := os.Create(local)
	if err != nil {
		panic(err)
	}
	fetchErr := client.Retrieve(remote, outFile)
	if fetchErr != nil {
		panic(fetchErr)
	}
	outFile.Close() // Close the file
}

func areFilesSame(localPath, remotePath string, client *goftp.Client) (bool, error) {
	localFileInfo, err := os.Stat(localPath)
	if err != nil {
		return false, err
	}

	remoteFileInfo, err := client.Stat(remotePath)
	if err != nil {
		return false, err
	}

	// Compare file size and modification times
	return localFileInfo.Size() == remoteFileInfo.Size(), nil
}
