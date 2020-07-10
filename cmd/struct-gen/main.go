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
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getJSONSchema(filename string, line int) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	fs := token.NewFileSet()
	fs.AddFile(filename, fs.Base(), int(fi.Size()))

	file, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}

	var jsonSchema string

	ast.Inspect(file, func(node ast.Node) bool {
		if node != nil && node.Pos().IsValid() {
			if fs.Position(node.Pos()).Line == line+1 {
				switch n := node.(type) {
				case *ast.BasicLit:
					if n.Kind == token.STRING {
						jsonSchema, err = strconv.Unquote(n.Value)
						if err != nil {
							err = fmt.Errorf("failed to unquote JSON schema string: %w", err)
						}
					}
				}
			}
		}
		return true
	})

	if err != nil {
		return "", err
	}

	if jsonSchema == "" {
		return "", errors.New("No var/const definition found containing a JSON schema string")
	}

	return jsonSchema, nil
}

func run() error {
	packageName := os.Getenv("GOPACKAGE")
	filename := os.Getenv("GOFILE")
	line, _ := strconv.Atoi(os.Getenv("GOLINE"))

	if packageName == "" || filename == "" {
		return errors.New("struct-gen: must be called from `go generate`")
	}

	if len(os.Args) < 2 || strings.TrimSpace(os.Args[1]) == "" {
		return fmt.Errorf("struct-gen: the first argument must be the struct name")
	}

	structName := strings.TrimSpace(os.Args[1])

	jsonSchema, err := getJSONSchema(filename, line)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile(os.TempDir(), "struct-gen-schema-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := f.WriteString(jsonSchema); err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	filePath := strings.Replace(filename, ".go", "_struct.go", 1)

	if filePath == filename {
		return fmt.Errorf("invalid source file %q, must end with '_schema.go", filename)
	}

	cmd := exec.Command("go",
		"run",
		"github.com/idubinskiy/schematyper",
		"-o",
		filePath,
		"--package",
		packageName,
		"--generator",
		"struct-gen",
		"--root-type",
		structName,
		"--ptr-for-omit",
		f.Name(),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("generating Go struct from JSON schema failed: %s", string(out))
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
