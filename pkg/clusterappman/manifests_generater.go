// +build ignore

package main

import (
	"log"

	"github.com/appvia/kore/pkg/clusterappman"
	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(clusterappman.Manifests, vfsgen.Options{
		PackageName:  "clusterappman",
		VariableName: "Manifests",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
