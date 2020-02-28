// +build tools
// Add hard dependency see - https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md#tools-as-dependencies

package tools

import (
	_ "github.com/shurcooL/vfsgen/cmd/vfsgendev"
)
