// +build ignore

package main

import (
	"log"

	"wab"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(wab.WebUI, vfsgen.Options{
		PackageName:  "dcui",
		BuildTags:    "!dev",
		VariableName: "WebUI",
		Filename:     "webui.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
