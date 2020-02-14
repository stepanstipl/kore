// +build ignore

package main

import (
	"log"

	"github.com/appvia/kore/pkg/clusterman"
	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(clusterman.Manifests, vfsgen.Options{
		PackageName:  "clusterman",
		VariableName: "Manifests",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
