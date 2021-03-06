// Code generated for package security by go-bindata DO NOT EDIT. (@generated)
// sources:
// rules/eks.yaml
// rules/gke.yaml
// rules/kubernetes.yaml
package security

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

var _rulesEksYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x91\xcf\x6a\x1b\x31\x10\xc6\xef\xfb\x14\x1f\xce\x79\x4d\x7b\x35\xe4\x60\xda\x1c\x4a\xa0\x09\x24\x3d\x85\x50\xc6\xd2\xd8\x16\xd6\x9f\xad\x66\x14\x63\xda\xbe\x4b\x9f\xa5\x4f\x56\x24\x39\x71\x12\x42\x7b\xed\x6d\xa5\x19\xfd\xe6\xfb\xcd\x8e\xe3\x38\x8c\x88\x14\x78\x81\x8b\xcb\x1b\x2c\x8b\x26\x31\xe4\x5d\xdc\x0c\x80\x49\xb6\xdf\x8f\xcb\x2f\xb7\x57\xe3\xbb\xf7\x03\x20\x4a\x5a\x64\x81\x3d\xe5\xd8\xbb\x2c\x8b\xc9\x6e\x52\x97\xe2\x02\x3f\x06\x00\x38\x3b\xc3\xd5\x03\xe7\x07\xc7\xfb\xa1\x5d\xdc\x6e\x9d\x20\x17\xcf\x30\x5b\x36\x3b\x81\x6e\xf9\x88\x42\x5a\xb7\x13\x9d\x46\x23\x35\x9a\xcc\x87\x23\xed\xf7\xaf\x8f\xac\xe4\xbc\xf4\x8b\xe5\xf3\xd6\xd8\x82\x4f\x9c\x83\xd3\xce\x35\x29\x6a\x4e\x1e\x93\xa7\x08\x4d\x8d\x1c\x48\x9d\x21\xef\x0f\xa8\x0f\xb9\xf6\x35\xd4\x3e\xe5\x1d\x67\xc4\x64\x59\xe0\x22\xa8\x7d\x62\x4a\xc9\xd7\xa7\x81\x59\x1b\x33\xf3\xb7\xe2\x32\x5b\x58\x0e\x14\xed\x53\x32\x7c\x0a\x13\x19\xed\xc7\xcf\x49\xc1\x91\x56\x2d\xd8\x6b\x27\x43\x11\x99\xa5\x78\x6d\x73\xa6\xc9\x3b\x43\x55\x13\x2b\x36\x29\xd4\x16\x21\x2d\x99\x94\x2d\x56\x87\x06\xcc\x25\xaa\x0b\xfc\x7c\x68\x5d\xe2\xe3\x9a\x27\x32\x3b\xda\x30\x84\x4d\xc9\x4e\x0f\x3d\x85\xe5\x35\xd5\x29\x41\x36\x38\xc7\xec\xd5\x6f\x85\x93\x9e\x91\x6d\xdd\x1d\x79\x7f\x32\x96\x59\x27\xd4\xc5\xdd\x05\xd9\xdc\xe3\x7b\x3b\x03\x2e\x4e\x45\xe7\x3b\x17\x2d\xce\xcf\x31\xbb\xf6\x14\x67\x2f\x4a\x32\xb1\x39\xd5\x2f\x2e\x6f\xde\x28\x9b\x14\xd7\x6e\x53\x0d\x5d\x8a\xf3\x3a\xf5\xba\x0e\xbd\xfb\x7a\x3f\xef\x89\x1e\x63\x72\xae\x94\x35\x79\xe1\x23\xa5\xba\x2c\xde\x96\xb1\x4e\x9e\x6c\x52\x64\xa4\x8c\x90\x32\xbf\xb0\xaa\x88\x9f\xdd\xcd\xf8\x22\xca\xf9\xaf\x7a\x1f\x7a\xcf\x7f\x66\xf8\x6f\xc1\x3f\x01\x00\x00\xff\xff\xb1\xb4\xad\x60\xce\x03\x00\x00")

func rulesEksYamlBytes() ([]byte, error) {
	return bindataRead(
		_rulesEksYaml,
		"rules/eks.yaml",
	)
}

func rulesEksYaml() (*asset, error) {
	bytes, err := rulesEksYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "rules/eks.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _rulesGkeYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x59\xdb\x72\xdc\xc6\xd1\xbe\xe7\x53\xf4\x4f\x5f\xf8\x4f\x15\x97\x45\x25\x4e\x2e\x58\xe1\x05\x4d\xb3\x24\x95\x2c\x66\x8b\xa4\xe4\x0b\x97\x2b\xd5\x3b\x68\x2c\x26\x1c\xcc\xc0\x73\xd8\xe5\x3a\xc9\xbb\xe4\x59\xf2\x64\xa9\x9e\x03\x16\x58\x82\x14\x45\x93\xb1\x72\xa7\x25\x30\xdd\x5f\x37\xbe\xaf\x0f\xa3\xd9\x6c\xb6\x37\x03\x8d\x2d\x1d\xc3\xeb\x77\xe7\x30\x37\x15\x9c\x06\x6f\x66\x4e\xa0\x92\x7a\xb9\x07\x20\x4c\x95\x1e\xce\xde\xcc\x4f\xf7\x00\x9c\x47\x1f\xdc\x31\xac\xd1\xea\xf4\x46\x45\x4e\x58\xd9\x79\x69\xf4\x31\xfc\x63\x0f\x00\xe0\xab\xaf\xe0\x2f\x2b\xb2\x2b\x49\xeb\xbd\xf8\x87\xeb\x46\x3a\xb0\x41\x11\x88\x86\xc4\x8d\x03\xdf\x10\x38\xf2\x5e\xea\x25\x98\x3a\xfe\x6c\x8c\x95\xbf\x18\xed\x51\x41\x67\x2a\xc0\x82\x83\xec\x5e\x36\xfa\xef\x7f\x7d\x47\x1e\xa5\x72\xc5\xe8\x43\x87\x40\xea\x18\x13\x2a\x65\xd6\xc9\xa1\x30\xda\x5b\xa3\x66\x9d\x42\x4d\xe0\x4d\xb4\x52\x6d\x34\xb6\x52\xa0\x52\x1b\x88\x27\x61\x63\x82\x05\x4b\x9d\x92\x02\x1d\x2c\xd0\x51\x05\x46\x47\x0b\x36\x68\x2f\x5b\x02\x4b\x3f\x07\x69\xa9\x25\xed\xdd\x61\x81\x07\x6f\xdb\x0e\x85\x4f\x3f\xaf\xd6\xd2\x8b\x86\xc3\xe3\x73\x6f\xe6\xa7\x60\xea\x1a\x2c\xb5\x66\x45\x09\x4d\x4d\xe8\x83\x25\xa8\x8d\x1d\xa2\x83\x11\xba\x16\x35\x2e\x29\xf9\xbe\x83\xa8\xa2\x16\x75\x15\xfd\x73\x6e\x4b\xf6\x3b\x14\x37\x7c\xc8\x91\x08\x56\xfa\x4d\x02\x54\x51\x8d\x41\x79\x68\xdd\x12\x4e\x60\xff\xcd\x36\x71\xe5\xab\x97\xbc\x39\x20\x8d\x0b\x45\xd5\x7e\x3a\xc9\x80\x7e\x6c\xdd\xf2\x27\xf8\x7b\xfc\x0d\x20\x75\x17\xfc\xe1\x8d\xd4\x15\x9c\x9c\xc0\xfe\x5c\xa1\xde\x1f\x3d\x72\x1d\x89\xed\xf3\xd7\xef\xce\xcb\x63\x6d\xfc\xf0\x15\x61\x74\x2d\x97\xc1\x22\xb3\xe7\x30\xb9\xdd\x22\x9b\x9b\x6a\x8b\x2b\x1b\x60\xf8\xc7\x9f\xc2\x5f\x49\x17\x03\x60\x0e\x70\xea\xba\x1e\xe0\x3f\xf7\x46\x84\xbf\x6a\x24\xa9\x8a\x2a\xb8\x30\x15\xb9\x11\xdd\xaf\xde\x7c\x37\x3b\x7a\xf5\xbc\x8c\x67\x97\xae\xb8\xd4\xec\xb2\x88\xc0\xf5\x24\x1a\x51\xbc\xc7\xc7\x27\x23\x46\x40\x4b\xb0\x08\x52\xf9\xc8\x49\xd3\xb1\x7c\xce\x4c\xdb\x05\x4f\x70\xae\x97\x52\xd3\xf6\xd4\xc7\xf7\xee\x30\xda\xf9\xa1\x21\x3d\x65\x6c\xfb\xad\x0f\x7a\x80\x2d\x3a\x4f\x16\x82\x63\x67\x20\xec\xa6\xf3\x66\x69\xb1\x6b\xa4\x88\xb6\x62\x3c\xe0\x0d\xac\xc8\xca\x7a\x03\xbe\x41\x0f\xb4\x22\xbb\x89\x21\x71\xce\xa3\x82\x84\x0a\xd1\x90\x64\x33\x2b\x69\x7d\x40\x05\x2d\xb2\x28\x28\x1a\xb2\x41\x73\x42\xa3\x50\x8d\x59\x2a\xfa\xda\x41\x85\x1e\x41\x90\xf6\x64\x0f\x53\x0e\x95\x6c\xa5\x4f\xe9\xc3\x85\x54\xd2\x6f\x38\x62\xd4\x80\xde\xa3\xb8\xc9\xc4\xf0\x06\x64\xdb\x91\x75\x46\xa3\x27\xc0\x49\x28\x9f\x56\xaa\x74\xbd\x2c\x93\x5c\xab\x20\xb2\x5c\xcd\x8a\x2c\x2a\xd5\xab\xaa\x94\xad\xa1\xed\xcf\x56\xe1\xd5\x98\x0c\xbf\xa9\xee\x0a\x96\x22\x85\xa1\xda\x76\x70\x9a\xc8\x7d\x68\xb8\x18\x11\xe9\x5e\x6f\x45\x63\x43\x89\x5d\x92\x22\x74\x04\x67\x0d\x6a\x4d\x6a\xa4\xb1\x8f\x57\x97\x2f\xa2\x31\x91\x7c\x31\x43\x9d\x34\x7a\x5a\x5b\xfc\x62\x67\xcd\x4a\x72\x44\x36\xa1\x74\xb0\x92\x58\x7e\x14\x33\x2e\x96\xe8\x77\x61\x41\x56\x93\x27\x77\x90\x7a\x4a\xa4\xae\x9f\x28\xd4\xe3\x42\x9e\x31\x44\xd2\x82\xf3\x96\xb0\x8d\xd4\x21\xe7\xa9\x2a\x4f\x0f\xe1\xad\xff\x3a\x25\xdd\x92\x30\x6d\x4b\x9a\xb3\xbd\x31\x81\x53\xb3\x81\xb5\xf4\x4d\xae\x65\xc3\xe7\x97\xe7\xaf\x3f\x7c\x7f\x7a\x09\xc6\xc2\xd5\xf5\xe9\xb7\xdf\x9f\x17\xe8\x59\xa8\x09\xfe\x34\xed\xaf\xb7\xdc\x65\x89\x2e\x88\xe3\xc9\x85\xa5\xae\x41\x1b\x3d\x1b\xba\xda\xc9\xc9\xd3\xe8\xfe\xc1\xb1\x93\x07\xcc\xbe\x00\xf3\x0b\x17\x8e\x4f\xee\x17\x40\x46\xb1\xa5\xe8\xf0\xe0\xff\x9d\xc0\x7e\x4e\xf4\xae\x4d\x7e\x94\x12\xbf\xbf\x23\x98\x0b\xe3\x21\xdc\x89\x96\x29\x57\xf8\xd0\x47\x3c\xd1\x96\x3e\x66\xde\x8e\xb4\x72\x7e\x39\x3b\xfa\xfd\xf3\x6b\xa5\xe7\x67\x3c\x7f\x4f\x1f\x7a\xbc\x56\x98\xc9\xee\x0e\x89\x83\x8b\xd2\x18\xd1\x12\xfe\x3f\x67\xf5\x20\x93\xf7\x77\x60\xd1\x37\x14\xa7\x21\xcd\x86\x37\xdc\xe2\x6a\x79\xbb\xd5\x89\x03\x74\xd1\xe0\xeb\xb3\x79\x32\x86\x1a\x3c\xde\x10\x08\x8c\x35\x3b\x86\x15\xba\xa5\xc5\x8a\xa0\x43\xdf\x38\x40\x5d\xf1\xbf\x52\x89\x37\xb9\x23\x2c\xa8\x41\x55\xdf\xd3\x10\x3a\x12\xb2\xde\xf0\xeb\x08\x2e\xfe\x90\x62\x94\xab\x4e\xa1\x88\x65\x90\x12\x81\x50\x81\x25\xd7\x19\xed\x64\x69\x51\x1a\x64\x34\xc6\xde\x1c\xa9\x9a\x3b\x54\x47\xb6\x36\xb6\x85\xd0\x55\xe8\x69\x8c\xec\x69\x82\x8a\xc5\x2e\x8b\xb8\x80\x2b\xb4\x7b\x69\x61\x15\x7f\x0f\x09\x6b\xd5\x13\x79\x78\x82\x65\xb3\x2b\x98\x61\x24\xdc\x57\x06\x99\x2f\xc7\x1c\xf9\x49\xb5\xcc\xad\x5c\x71\xd3\xbf\x20\xbf\x36\xf6\x66\x77\x6f\xb9\x38\xbf\x7e\x91\x26\xd3\x65\xb7\xba\x77\xfb\x18\x0d\xdd\x39\x53\x56\x02\xba\xf5\x64\x99\x49\x58\x55\x96\x1c\xeb\xab\xb6\xa6\x8d\xee\xf8\x5d\xb2\xa9\xf9\x96\x81\x4e\x8a\x06\x5a\x42\xed\x60\xcd\xf6\xa8\x62\x86\x09\x4b\x69\xfc\x29\x8e\x3a\x22\xcb\x5e\x84\xd1\x9a\x44\xec\xda\x0b\xf2\x6b\x6e\xda\x45\x8f\x31\xf3\xa3\x96\xc5\xbc\xdc\x75\x1b\xc5\xce\x7f\x4c\xe3\xe1\xdb\xf9\xc7\x6f\xce\x64\x65\xcb\xfc\x9a\x50\x5d\x18\x4f\x69\x8a\xca\x13\x02\xcf\x53\x52\xd3\x01\x2c\x82\xe7\x66\x59\x19\x1e\x3a\xb5\x08\x69\xe5\xe9\x48\xa3\xf2\x79\x84\x5c\xb3\xea\x63\x00\xe9\xe3\xc0\x70\xbc\x62\xb1\x6c\x7a\x82\x6c\xe7\x5f\x5e\x23\xb1\xed\x13\x0a\xef\x3f\x5c\x5d\x03\x69\xc7\x13\xdc\x16\x6d\x2a\x15\x11\x6e\x58\x68\xf2\x8c\x2b\x68\xf9\x73\xa0\x61\x67\x2d\x36\x38\xfc\x08\x94\x87\x26\x1e\xfb\x14\x76\xd3\x95\xe2\x87\x21\xe4\x1e\xdd\x3a\x7e\x1b\x4b\x4e\x56\x34\x02\xd9\x59\xf3\x37\x12\xbe\x38\x4a\x23\x77\x70\x19\x1f\xd3\xac\x14\x11\x95\xd6\xc2\x3b\x81\xa4\x20\x76\xc1\xc5\x20\x00\xb5\xe1\xea\xb9\x8d\x35\xba\xc6\xde\xdd\xd3\x2b\xcc\xfc\x2e\x6b\x07\x4b\xd6\x0b\x94\x96\x4f\x4c\xaa\x19\x4f\xd6\x3b\x1f\xf6\x36\xd0\x44\x41\x99\x06\xee\xe2\xc0\x1f\x57\xe8\xc9\x7a\x52\x76\xc9\xdd\x42\x72\xf4\x87\x67\xbe\xff\x88\xa6\x4a\xd3\xc2\xe0\x0d\xef\xf8\x28\x6d\x21\x4d\xac\x18\x8a\x25\x6e\xfa\x1d\x66\x30\xcd\x8d\xaa\xcb\x00\x34\x1f\x8f\x47\xc9\xf6\xab\xd3\x50\xe0\x5c\x28\xca\xab\xe9\x59\xda\x1e\x2d\x39\x13\xac\xc8\xc3\x7f\xe8\xb2\x10\xd6\x9a\x9b\xee\xd4\xcd\x47\x6c\x2e\x4b\x56\xf4\x21\x1d\xb2\xac\x2d\xc5\xab\x09\x65\xb0\x9a\x16\x4c\x2c\x10\x0d\xae\xfa\x6d\xab\xdc\x01\xe5\xcd\xa7\x54\x34\x63\x6f\xd8\x88\x83\x56\x2e\x1b\x1f\xb9\xbe\xa0\x74\x3d\x53\x45\x43\x28\x84\xb1\x95\xd4\x4b\xb5\xe1\x70\x38\x8c\x5f\x73\x1b\xb2\xf3\xd9\x07\x7b\x31\x27\x93\xd7\xbe\xb8\x4f\x76\xc6\x28\xf7\x2c\x7c\x7f\x1c\xdd\xd9\xe9\x9c\x7d\xfe\xf8\xd7\x9f\x32\xf7\x07\x17\x1d\x27\x50\xa3\x72\x53\xbc\xdf\x09\xa5\xbf\x0e\xe1\x31\x4f\x13\xb3\xa9\x35\x96\x46\x31\x65\x1d\xc4\xe2\x91\x98\xf6\x60\x70\x67\xe9\x9d\x2f\x2a\xbe\x47\x84\xb7\x2b\x73\xb8\x8c\x82\x1b\xcb\xfc\x9b\x2f\x5c\xe6\x83\xe3\x9f\x92\x39\x7b\x6b\xd1\xe7\xeb\x4d\x4b\x71\x70\x4d\x0d\x3d\x9a\x4b\xad\xaa\xc1\x15\x41\x8d\x92\x29\xc2\x56\x6e\xfa\x55\x17\x1a\x42\xe5\x9b\x1c\xce\xbd\x9a\x1e\x49\x7a\x2c\x65\x1e\xd6\x5b\xdc\x24\x1f\x41\x3b\x2e\xbc\x41\xc5\xcf\x95\xef\x12\xca\x80\xc1\x4d\x2c\x1a\xcc\x3e\xa5\x73\x81\xdc\xaf\xd3\x73\xfe\xbe\x5f\xb8\x9e\xf3\xf7\x7c\x98\xef\x83\x50\xfe\xd7\xf4\xfc\xc4\xf8\x3e\x53\xcf\x17\xa6\x22\x7e\x3a\xb9\x3d\x1f\xfd\xf1\xa5\x44\xad\x8b\xdb\x72\x95\x73\xcf\xff\x56\x9c\x2a\xdf\x98\xb0\x6c\xc0\x93\x68\x74\x56\x64\x67\x5c\x9a\xf6\x78\x76\x8f\xcd\xf4\x00\x6e\x88\xba\xd2\x9a\xfb\x90\x61\xb8\x49\x49\x0d\x6e\xa3\x45\x1a\xfb\xee\xde\x38\x71\x5b\xf5\xd6\xc4\xde\x38\x5c\xc4\x83\x56\xe4\x92\x22\x59\x8d\xa9\x8f\xf2\x59\x61\x29\xb6\x7f\x74\xf1\x4a\x99\x7b\x7d\xb9\x41\x8d\xbe\xd1\x4f\xac\x63\xd3\xc5\xe0\xdb\x4d\xdf\xd2\x27\x4f\x45\xef\x69\xe5\xd9\xbd\xd2\xe5\x6f\x98\x2f\xd3\xc6\x65\x2b\xad\xcb\xdb\x6d\x39\x6d\xf0\x25\xef\x4f\x2f\x10\xbb\x84\x79\xfe\x42\xf0\x98\x15\x79\xa4\x99\xcf\xdb\x97\xfb\x00\x62\xb9\xe7\xad\xb9\x4f\xf7\xa0\x84\x4f\xad\xcf\x2f\x52\x16\xfe\x5b\xd1\x7e\x4e\xa0\x93\x05\xe2\x3d\xde\xc2\xdc\x54\xe3\xff\xef\x39\xfa\xd3\x33\x57\x08\xf9\x4b\x7f\x21\xd5\xe2\x2d\x74\xa6\x72\xdc\xab\xcb\x16\x3f\x5d\x25\xae\xb2\x78\x62\x27\x5d\xa1\x0a\x5c\x1a\x0c\x28\xb3\x06\x61\x82\x62\x11\x57\x64\x21\x78\xa9\xd8\xfe\x76\x4b\x77\xd9\x2e\x5c\xe7\xb5\x9e\x6e\x3b\x12\x9e\xaa\x9e\xf9\xd2\xc1\xab\x57\x47\x63\x1c\xd3\x22\x3e\x8b\x9e\x14\x61\x16\x69\x6e\xda\xec\x37\xd6\xf4\xb4\x40\xa4\xbf\x95\x2e\x5e\x03\x8e\xb1\xa5\x39\x7d\x70\x85\xf0\x0c\x42\x2d\x1f\xee\xb7\x6a\xd9\x2d\xde\xb2\xfb\x39\x59\x86\x04\x7f\xe6\x84\x3e\xc4\xd4\xfe\xbb\x4b\xc7\x9f\xf0\x4b\xe8\xca\xcf\x14\xc2\x7f\x02\x00\x00\xff\xff\xfd\x17\xee\x47\x3a\x20\x00\x00")

func rulesGkeYamlBytes() ([]byte, error) {
	return bindataRead(
		_rulesGkeYaml,
		"rules/gke.yaml",
	)
}

func rulesGkeYaml() (*asset, error) {
	bytes, err := rulesGkeYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "rules/gke.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _rulesKubernetesYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x54\xc1\x6e\x1b\x37\x10\xbd\xeb\x2b\x5e\x9d\x43\x2e\xb6\xea\x5e\x0d\xf8\xe0\xd6\x01\x22\xb8\x40\x04\xd7\x41\x0f\x41\x10\x50\xe4\x68\x77\x60\x2e\xb9\xe0\xcc\x6a\x2b\xb4\xf9\xf7\x62\xb8\x5a\x59\x6e\xec\x1e\x82\x1e\x0a\x9d\xb4\x24\xdf\xbc\x79\xf3\xe6\x5d\x5c\x5c\x2c\x2e\x90\x5c\x47\x57\xb8\x19\xb4\xc5\xba\xe4\x3f\xf6\x58\xad\x71\xef\x52\x43\xb2\x00\x7c\x0e\x76\xf8\xf1\xe1\xfd\x6a\x7d\x71\xf9\xd3\x02\x10\x75\x3a\xc8\x15\x46\x57\x12\xa7\x66\x01\x04\x12\x5f\xb8\x57\xce\xe9\x0a\x7f\x2d\x00\xe0\xcd\x1b\x7c\xd8\x51\xd9\x31\x8d\x8b\xfa\xe1\xa1\x65\x41\x19\x22\xc1\xb7\xe4\x1f\x05\xda\x12\x12\xe9\x98\xcb\x23\xb6\x1c\x95\x8a\x20\xa7\xfa\xd9\x0d\xda\x52\x52\xf6\xce\x20\xd1\x1b\xa9\xe5\x62\xc6\xbd\x25\x75\x1c\x65\xfa\xff\x7b\x4b\x09\x77\xb9\x10\x7c\x21\xa7\x24\x70\xb8\x1b\x36\x54\x12\xd9\x1f\x1f\x07\x51\x2a\xe7\x60\xc5\x20\x76\x9a\x5e\x44\xaf\x58\x65\x48\xd6\x10\x38\x09\x07\x82\xb6\x4e\x67\x00\x68\x3e\x7d\x47\x70\xde\x93\x88\x7d\x36\xc2\x87\x5b\x4b\xac\x14\x2c\x15\x6c\x43\xa2\xe8\x8b\xf3\xca\x9e\xec\x5e\x21\xd1\xc2\x5e\xeb\x83\xd5\x1a\x2e\x84\x62\x10\xa5\x2a\x0d\x4a\x6e\x13\x29\x60\xb3\x47\xa0\xad\x1b\xa2\x9a\x1a\xae\x62\x9d\x92\xc0\x63\xca\x63\x82\x90\x22\x6f\x0d\x67\x7a\xff\x24\xcf\xaa\xeb\x9d\xd7\x59\xf4\x97\xc5\x44\xa0\x3e\xe6\x3d\x05\xb0\x40\xc8\x0f\xa5\x52\xdc\x10\x72\x4f\x69\xee\x8a\x93\x56\x1d\xcf\x2b\x96\x64\xf0\x16\x89\xac\x6f\x57\xf6\xa6\x28\x4b\xd5\xa1\x57\xa3\x5e\x7b\x1c\xd2\x4c\x56\x30\xb2\xb6\x79\xd0\x63\xe3\x26\xad\xb6\x34\x49\x6d\xa4\x97\x78\x9f\x47\xda\xd9\x78\xc6\x96\x0a\xa1\xcf\x22\xbc\x89\x74\x5e\xcb\xd7\x3b\x90\x36\x0f\x31\x18\xb5\x19\x87\xc2\x44\x30\x8b\xe9\xf8\x04\x27\x07\x94\x7d\x1e\x0a\x5c\xe8\x38\xb1\x68\x71\x9a\x2b\x97\x18\x8f\x33\x7b\x1a\x18\xb6\x25\x77\x55\x3a\x73\xe6\xec\xdd\xde\xf9\x47\x67\xa5\x4d\x18\xd6\xfd\xa4\xe5\x3c\x95\x4e\x1a\x5c\xe3\xec\x26\xc6\xb9\xac\xf4\xe4\x79\xcb\x14\x26\x6b\x53\x80\x4b\x01\x3e\x77\x7d\x64\x97\xf4\x6c\x7a\xdf\x47\x97\x3e\x75\xd2\x7c\xc6\x9f\xf5\x3f\xc0\xa9\x1f\x74\xf9\xc8\x29\xe0\xfa\x1a\x67\xeb\xe8\xd2\xd9\xb3\x23\x03\x5e\xfa\x9c\xb6\xdc\x0c\xa5\x8e\x6f\x69\xd3\xac\x4b\x7a\x13\x63\x1e\x29\xac\xd6\xf2\xe9\xcb\x67\x23\x74\xb9\xac\xbf\x1f\x2f\x67\x0c\x23\x7a\x65\x4c\x9f\x1b\x60\x5a\xf1\x9e\x4a\xc7\xaa\x93\x03\xe6\x35\x9c\x14\xaf\x16\x30\xa7\x1c\xf9\x7c\x5d\x9c\xda\xf0\x5f\x9b\xf8\x65\xba\xf3\x7f\xea\xc3\x9f\x52\xfa\xba\x38\x06\xde\xed\x61\xa0\x0f\xe4\x3a\xdc\xe7\x48\xc7\xb4\xbb\xfb\xf8\xf3\xbb\x2f\xf7\x1f\x7e\x7d\xf7\xdb\x7f\x9a\x78\xbc\x3d\x9a\x48\xad\x66\xc9\x91\xe6\xc8\x9b\x0d\xd9\xf0\x8e\x2c\x56\x32\xba\xc1\xb7\x70\xa3\x7b\x2d\xf9\x6c\xb7\xbf\x85\x1b\x5b\xf6\xad\x69\xc1\xa9\xa5\xc2\x3a\x25\x8a\x8b\xd1\xe2\xaf\xd4\xc0\x12\xa2\x4e\xb0\xa3\xb2\x47\xcc\x59\xa6\x85\x34\xc7\xb2\x58\x79\x5b\x53\x77\xe4\x53\xf7\xa8\x06\x4f\x8c\x53\x99\x8e\xba\x0d\x15\x79\x35\x6f\x66\x4e\x7d\xe1\x1d\x47\xb2\x05\xc9\xdb\xda\xe3\xc9\x73\xab\x55\x19\xb4\xdc\xb4\xe7\x18\xe9\x6d\x80\x0c\x4d\x43\xa2\x93\xd5\x5a\x97\x9a\x1a\x18\x19\x6f\x4d\xd2\xb7\xdf\xb5\xa4\xb7\x2f\xc9\x5d\xc8\xe7\xae\xa3\x14\x28\x58\x8c\x5a\x2e\x7d\xcf\x8a\x56\xb8\xab\xeb\xd7\x2d\x7e\xe0\x62\xf6\x3a\xb8\xeb\xf8\xec\x87\x6b\x9c\x59\x57\xff\xb4\xf8\xb7\x7c\x53\xd6\x67\x84\x27\x31\xea\xd9\x6c\xe8\xbf\x03\x00\x00\xff\xff\xa8\x67\xdf\xcd\xc5\x07\x00\x00")

func rulesKubernetesYamlBytes() ([]byte, error) {
	return bindataRead(
		_rulesKubernetesYaml,
		"rules/kubernetes.yaml",
	)
}

func rulesKubernetesYaml() (*asset, error) {
	bytes, err := rulesKubernetesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "rules/kubernetes.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"rules/eks.yaml":        rulesEksYaml,
	"rules/gke.yaml":        rulesGkeYaml,
	"rules/kubernetes.yaml": rulesKubernetesYaml,
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
	"rules": {nil, map[string]*bintree{
		"eks.yaml":        {rulesEksYaml, map[string]*bintree{}},
		"gke.yaml":        {rulesGkeYaml, map[string]*bintree{}},
		"kubernetes.yaml": {rulesKubernetesYaml, map[string]*bintree{}},
	}},
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
