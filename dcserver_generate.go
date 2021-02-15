// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"

	"github.com/mrecachinas/dcserver/internal/ui"
)

func main() {
	err := vfsgen.Generate(ui.WebUI, vfsgen.Options{
		PackageName:  "dcui",
		BuildTags:    "!dev",
		VariableName: "WebUI",
		Filename:     "webui.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
