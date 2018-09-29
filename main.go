package main

import (
	"fmt"
	"github.com/cbroglie/mustache"
)

func main(){
  data, _ := mustache.Render("hello {{c}}", map[string]string{"c": "world"})
  fmt.Println(data)
}
