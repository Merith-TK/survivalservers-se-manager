package main

import (
	"fmt"
	"os"

	"github.com/secsy/goftp"
)

func pullData(client *goftp.Client) {
	fmt.Println("Pulling data from FTP server")
	// for each entry in conf.Files, key is remote, value is local
	for _, file := range conf.Files {
		fetchFolder(client, file.Remote, file.Local)
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
