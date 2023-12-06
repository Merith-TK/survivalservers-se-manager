package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileHash struct {
	FileName string `json:"file_name"`
	Hash     string `json:"hash"`
}

type FolderHash struct {
	FolderPath string     `json:"folder_path"`
	Hashes     []FileHash `json:"hashes"`
}

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}

func generateHashes(folderPath string) (FolderHash, error) {
	fileInfos, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return FolderHash{}, err
	}

	var fileHashes []FileHash

	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()
		if !fileInfo.IsDir() && !strings.HasSuffix(fileName, ".md") {
			filePath := filepath.Join(folderPath, fileName)
			hash, err := calculateMD5(filePath)
			if err != nil {
				return FolderHash{}, err
			}
			fileHashes = append(fileHashes, FileHash{FileName: fileName, Hash: hash})
		}
	}

	sort.Slice(fileHashes, func(i, j int) bool {
		return fileHashes[i].FileName < fileHashes[j].FileName
	})

	folderHash := FolderHash{
		FolderPath: folderPath,
		Hashes:     fileHashes,
	}

	return folderHash, nil
}

func saveJSON(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, jsonData, 0644)
}

func processFolder(folderPath string) error {
	// ignore .git so it doesnt break itself
	if strings.HasSuffix(folderPath, string(os.PathSeparator)+".git") {
		return nil
	}

	folderHash, err := generateHashes(folderPath)
	if err != nil {
		return err
	}

	fmt.Printf("Folder: %s\n", folderPath)

	jsonFilePath := filepath.Join(folderPath, ".hashes.json")
	if err := saveJSON(jsonFilePath, folderHash); err != nil {
		return err
	}

	fmt.Printf("JSON file saved: %s\n", jsonFilePath)

	return nil
}

func walkFolders(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return processFolder(path)
		}

		return nil
	})
}

func main() {
	rootPath := "." // Change this to the root folder you want to start the traversal from
	if err := walkFolders(rootPath); err != nil {
		fmt.Println("Error:", err)
	}
}
