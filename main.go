package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/cbroglie/mustache"
)

// Root directory path
// TODO: Hard code
const DirPath = "."

func main(){
	// Each file in the root directory
	// (from: https://flaviocopes.com/go-list-files/)
  filepath.Walk(DirPath, func(path string, info os.FileInfo, err error) error {
  	// If path is not directory path
  	if !info.IsDir(){
  		// Just print file path
  		fmt.Println(path)
  	}
  	return nil
	})
	data, _ := mustache.Render("hello {{c}}", map[string]string{"c": "world"})
	fmt.Println(data)
}
