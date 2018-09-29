package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/cbroglie/mustache"
)


func replaceInDir(dirPath string) error {
	// Each file in the root directory
	// (from: https://flaviocopes.com/go-list-files/)
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// If path is not directory path
		if !info.IsDir(){
			// Just print file path
			fmt.Printf("====== %s ======\n", path)
			data, _  := mustache.RenderFile(path, map[string]string{"c": "world"}) // TODO: Hard code
			ioutil.WriteFile(path, []byte(data), info.Mode())
		}
		return nil
	})
	return nil
}

func main(){
	// Get root directory path
	// TODO: Error handling
	dirPath := os.Args[1]
	// Replace files in the directory
	replaceInDir(dirPath)
}
