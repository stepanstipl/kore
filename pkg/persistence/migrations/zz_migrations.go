// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// pkg/persistence/migrations/files/.1595965872_migrate_users.up.sql.swp
// pkg/persistence/migrations/files/1595965872_migrate_users.up.sql
package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1595965872_migrate_usersUpSqlSwp = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\xda\xcd\x4e\xf2\x40\x14\x06\xe0\xc3\xb7\xff\x82\xd1\x1b\x18\x70\x81\x26\xda\xa1\x18\x14\x34\xae\x4c\x8d\x24\x02\x09\xa0\xae\x4c\x53\xda\xa1\x8e\xfc\xb4\xce\x0f\x89\x1b\xbd\x04\xaf\xc2\x3b\xd4\xbd\xa9\x60\x74\x63\x64\x25\x51\xdf\x67\x33\xf3\xe6\x9c\x4c\xce\x76\x32\xd3\x2f\x9f\x37\x9a\xac\xe6\x54\x88\x88\x56\x88\xa2\xa7\xa2\x7f\xe9\xee\x53\xf8\x4c\x74\x2d\xb4\xa1\x05\x68\x13\xa8\x81\x55\xb7\x5f\xf5\xdd\x65\x07\xf2\x38\xe1\x5a\x85\x3c\x96\xe6\xca\xf6\x9d\x30\x19\xf3\x20\x4d\xa7\x32\xe0\xc3\x44\x09\x9e\x0e\x63\x9e\x0a\xa5\xa5\x36\x62\x12\x0a\x3e\x96\xb1\x0a\x8c\x4c\x26\x9a\x0f\xe4\x48\x68\xee\x56\xeb\xd5\xfa\x6e\xb5\xb6\x57\xf1\x67\x35\xe1\x5b\x2d\x94\x76\x6c\xea\xe8\x9b\xd1\x22\xf3\x02\xfc\x39\xd6\x0c\xb6\x6b\xff\x69\xa7\xe2\x96\xb3\xb8\x5e\x2c\xb0\xb5\xd5\xb3\x65\x4f\x05\x00\x00\x00\x00\x00\xdf\xc8\xa4\x39\xba\x27\xa2\x7f\xf3\x9c\xfb\x64\x05\x00\x00\x00\x00\x00\x00\x80\x9f\x2b\x88\x88\x92\x3c\xd1\x43\x7e\xf6\xfe\xff\x76\xdf\xcf\xf2\x63\x7e\xc9\xc3\x01\x00\x00\x00\x00\x00\x00\xfc\x22\x8d\x56\xd7\xeb\xf4\x58\xa3\xd5\x6b\x7f\xd8\x32\x19\x89\x89\x91\x46\x0a\xcd\x36\xac\x16\xca\x97\xd1\x16\x4b\x55\x32\x95\x91\x50\xbe\x18\x07\x72\xf4\x9e\x37\x59\xd7\x3b\xf5\x8e\x7a\x2c\x6b\x9a\xd7\x4a\x5a\x27\x25\x76\xdc\x69\x37\xd9\xeb\x17\x79\x66\xd9\xc5\x89\xd7\xf1\x98\x75\x64\xc4\x0a\x87\xcc\x3d\xa0\x97\x00\x00\x00\xff\xff\xc0\x90\x9b\xb6\x00\x30\x00\x00")

func _1595965872_migrate_usersUpSqlSwpBytes() ([]byte, error) {
	return bindataRead(
		__1595965872_migrate_usersUpSqlSwp,
		".1595965872_migrate_users.up.sql.swp",
	)
}

func _1595965872_migrate_usersUpSqlSwp() (*asset, error) {
	bytes, err := _1595965872_migrate_usersUpSqlSwpBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: ".1595965872_migrate_users.up.sql.swp", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1595965872_migrate_usersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\xf4\x0b\x76\x0d\x0a\x51\xf0\xf4\x0b\xf1\x57\xc8\x4c\x49\xcd\x2b\xc9\x2c\xc9\x4c\x2d\x56\xd0\x28\x2d\x4e\x2d\x8a\xcf\x4c\xd1\x51\x28\x28\xca\x2f\xcb\x4c\x49\x2d\x8a\x4f\xcd\x4d\xcc\xcc\x41\xf0\x35\x15\x82\x5d\x7d\x5c\x9d\x43\x14\x40\x8a\xa0\x72\xea\xc5\xc5\xf9\xea\x0a\x6e\x41\xfe\xbe\x0a\x20\xfd\xc5\x0a\xa5\x0a\xe1\x1e\xae\x41\xae\x0a\xa5\x7a\x99\x29\x0a\x8a\xb6\x0a\x86\xd6\x5c\x80\x00\x00\x00\xff\xff\xed\x1f\x97\x34\x71\x00\x00\x00")

func _1595965872_migrate_usersUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1595965872_migrate_usersUpSql,
		"1595965872_migrate_users.up.sql",
	)
}

func _1595965872_migrate_usersUpSql() (*asset, error) {
	bytes, err := _1595965872_migrate_usersUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1595965872_migrate_users.up.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	".1595965872_migrate_users.up.sql.swp": _1595965872_migrate_usersUpSqlSwp,
	"1595965872_migrate_users.up.sql":      _1595965872_migrate_usersUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	".1595965872_migrate_users.up.sql.swp": {_1595965872_migrate_usersUpSqlSwp, map[string]*bintree{}},
	"1595965872_migrate_users.up.sql":      {_1595965872_migrate_usersUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
