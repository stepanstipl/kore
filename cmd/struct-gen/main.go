/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

const tmpl = `// Code generated by struct-gen; DO NOT EDIT.

package %s
`

func getType(filename string, line int) (*ast.File, *ast.StructType, string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, nil, "", err
	}

	fs := token.NewFileSet()
	fs.AddFile(filename, fs.Base(), int(fi.Size()))

	var res *ast.StructType
	var name string

	file, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, "", err
	}

	ast.Inspect(file, func(node ast.Node) bool {
		if node != nil && node.Pos().IsValid() {
			if fs.Position(node.Pos()).Line == line+1 {
				switch n := node.(type) {
				case *ast.TypeSpec:
					if str, isStruct := n.Type.(*ast.StructType); isStruct {
						res = str
						name = n.Name.Name
					}
				}
			}
		}
		return true
	})

	if res == nil {
		return nil, nil, "", fmt.Errorf("'%s' does not refer to a struct", name)
	}

	return file, res, name, nil
}

func main() {
	packageName := os.Getenv("GOPACKAGE")
	filename := os.Getenv("GOFILE")
	line, _ := strconv.Atoi(os.Getenv("GOLINE"))

	if packageName == "" || filename == "" {
		fmt.Println("must be called from `go generate`")
		os.Exit(1)
		return
	}

	_, _, structName, err := getType(filename, line)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
		return
	}

	//reflector := jsonschema.Reflector{}
	//reflector.ExpandedStruct = true
	//
	//typeName := packageName + "." + structName
	//value, found := types[typeName]
	//if !found {
	//	fmt.Printf("Error: unknown struct %q, you need to register it in cmd/jsonschema-gen/main.go\n", typeName)
	//	os.Exit(1)
	//	return
	//}
	//
	//schema := reflector.Reflect(value)
	//res, err := schema.MarshalJSON()
	//if err != nil {
	//	fmt.Printf("Error: %s", err.Error())
	//	os.Exit(1)
	//	return
	//}
	//
	//schemaStr := string(res)
	//tmp := map[string]interface{}{}
	//_ = json.Unmarshal([]byte(schemaStr), &tmp)
	//formattedSchemaStr, _ := json.MarshalIndent(tmp, "", "  ")
	//
	//contents := fmt.Sprintf(tmpl, packageName, structName, formattedSchemaStr)
	//filePath := strings.Replace(filename, ".go", "_schema.go", 1)
	//
	//if filePath == filename {
	//	fmt.Printf("Error: invalid source file, must end with .go: %s", filename)
	//	os.Exit(1)
	//	return
	//}
	//
	//if err := ioutil.WriteFile(filePath, []byte(contents), 0644); err != nil {
	//	fmt.Printf("Error: %s", err.Error())
	//	os.Exit(1)
	//	return
	//}
}
