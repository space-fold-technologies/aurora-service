// Code generated for package server by go-bindata DO NOT EDIT. (@generated)
// sources:
// resources/application.yml
// resources/boot.txt
// resources/permissions.json
// resources/settings.yml
// resources/migrations/1_core_tables_schema.down.sql
// resources/migrations/1_core_tables_schema.up.sql
package server

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

var _resourcesApplicationYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\xca\x31\x6f\x85\x20\x10\x00\xe0\xdd\xc4\xff\x70\x1b\x53\x5f\x00\x51\x9f\xcc\x5d\x3b\xbd\xee\xe6\x0a\x77\x91\x44\x39\x03\x68\xd2\x7f\xdf\x74\xe8\xd4\xf1\x4b\xbe\x8c\x07\x79\xc0\xab\x48\xc1\xb7\x4a\xe5\x4e\x81\xfa\xee\xa6\x52\x93\x64\x0f\xe6\xa1\x1f\xba\xef\x0a\x9d\x52\x53\x93\xf2\xfd\xff\x52\xbe\x53\x91\x7c\x50\x6e\x1e\x54\xa4\x9b\x76\x39\x7f\xa5\xfa\x2e\xc8\x71\xa4\xb6\x6e\x58\x37\x0f\xd1\x44\x64\x33\x2d\xac\x9f\x93\x9b\xc3\x30\x32\x07\x33\xd0\x12\x71\x32\x96\xc7\x39\x30\x31\x8e\x8e\xfb\xee\xeb\x4a\x7b\x5c\x23\x36\xf2\xa0\x5e\x57\x86\x17\x9d\x0d\xc0\x2c\x60\x9c\xb7\xce\xbb\x01\x3e\xde\x3f\xc1\x6a\x6b\xd5\xdf\xa6\x53\xc2\xb6\x56\x0a\x1e\x94\x19\xb5\x73\xb3\xb6\xcf\x41\xfd\x04\x00\x00\xff\xff\xd0\xed\x29\x1d\xe0\x00\x00\x00")

func resourcesApplicationYmlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesApplicationYml,
		"resources/application.yml",
	)
}

func resourcesApplicationYml() (*asset, error) {
	bytes, err := resourcesApplicationYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/application.yml", size: 224, mode: os.FileMode(511), modTime: time.Unix(1669345184, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesBootTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x90\xe1\x09\xc4\x20\x0c\x85\xff\x0b\xee\xf0\x6d\xd0\x85\x02\x2e\xd2\x2d\x0e\x6e\xc0\x9b\xe4\x50\x5b\x94\x44\x1a\x2f\x7f\x4e\x82\x0d\xcf\x2f\xbe\xfa\xa0\xb4\x55\xbf\xd0\xf6\x72\x2b\xd3\x91\x56\x72\xe2\xa0\x2d\x39\x41\xa0\xef\x5d\x99\x8f\x94\x42\x4e\x27\x7c\xde\x2f\x55\x72\x69\xf5\xa2\xab\xb3\x50\xe3\x1e\x86\xb7\xd8\xfa\x03\xbd\x2d\x45\x79\x19\xf3\x00\x31\x84\x61\x54\x5f\xfe\x7c\x4d\x88\xc0\x18\xd9\xda\x8e\xf5\x37\xce\xd6\xfc\x5a\x13\xc9\xe1\xa6\xe8\x11\x43\x58\x19\x89\x9b\xa2\x4f\x2c\x8c\x18\x53\x77\x27\x8b\x88\x83\x94\x11\x72\x22\xb6\xfe\x31\xf8\x0d\x00\x00\xff\xff\x61\x15\x01\x0e\x38\x04\x00\x00")

func resourcesBootTxtBytes() ([]byte, error) {
	return bindataRead(
		_resourcesBootTxt,
		"resources/boot.txt",
	)
}

func resourcesBootTxt() (*asset, error) {
	bytes, err := resourcesBootTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/boot.txt", size: 1080, mode: os.FileMode(511), modTime: time.Unix(1669345184, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesPermissionsJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\x4d\xaa\xc3\x30\x0c\x06\xf7\x81\xdc\xc1\x64\xfd\x78\x17\x2a\x5d\x98\x58\xa5\x06\x5b\x32\xfe\x09\xf4\xf6\x5d\xa8\xb8\xb5\xa4\xed\x4c\x32\x9f\xf0\x6d\xdf\x9c\x73\xee\xf0\xa5\xb4\xff\xb3\x82\xef\x70\xfc\xfd\xb2\x51\x82\x62\x01\x4a\xa2\xd7\xca\x22\x3e\xa8\x66\xdf\x23\xe1\x2a\xda\x13\x52\x9a\xe8\x4c\xa3\x75\xa8\x6a\x6b\x72\xb1\x37\xb9\xd5\x9f\xb2\x42\xa6\xeb\xfb\x13\xe0\x15\x2b\x61\x06\xec\x6a\x68\x71\x56\x74\xf9\x40\x84\x91\x02\xa8\x22\x43\x71\x37\x43\xab\xcf\x46\x84\x3b\xf8\xac\xc2\x0c\xad\x06\x1b\xd1\x18\xcd\x78\x57\x86\x56\x83\x8d\x38\x9b\xe1\x27\xbc\x6f\xf7\x77\x00\x00\x00\xff\xff\x24\x7f\xb8\x98\x1f\x02\x00\x00")

func resourcesPermissionsJsonBytes() ([]byte, error) {
	return bindataRead(
		_resourcesPermissionsJson,
		"resources/permissions.json",
	)
}

func resourcesPermissionsJson() (*asset, error) {
	bytes, err := resourcesPermissionsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/permissions.json", size: 543, mode: os.FileMode(511), modTime: time.Unix(1669345184, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesSettingsYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x91\x4f\x6f\xdb\x30\x0c\xc5\xef\x02\xf4\x1d\x78\xeb\x49\x6a\xec\xa4\x49\xa7\xdb\xb0\xed\x34\x0c\x03\xba\xc3\xce\x9c\x4d\x27\x5c\xf4\xc7\xa0\xe8\xac\xf9\xf6\x83\xb5\xb4\x05\x56\xf8\x60\xe0\xc7\xf7\xf0\x28\xbe\x59\xca\x85\x47\x92\x00\x9f\xbf\x7f\xfa\xfa\xe5\xc9\xfd\xf8\xf9\xf1\xe9\x9b\x35\xb3\x94\x89\x23\xb9\x91\x85\x06\x2d\x72\x0d\x70\x4f\x3a\xdc\xe3\x22\x45\xd0\x9a\x91\x26\x5c\xa2\xba\x21\x2e\x55\x57\xfb\x0d\xbc\x9b\xb8\xc8\x55\x29\x3b\x1c\x47\xa1\x5a\x03\xdc\x6d\x7c\xfb\x42\xbf\x3d\x1c\xee\xde\xeb\x71\xbc\x90\x28\x57\x7a\xb3\x74\x1b\x7f\xf0\x1f\xf6\x7e\xf7\xa6\x56\xc2\x14\x00\xc7\xc4\xd9\x9a\xb1\x24\xe4\x1c\xe0\xf0\xe8\xbb\x5d\xe7\xfb\xee\xd1\x77\x87\xad\xcf\x3c\x7b\x2e\xd6\x34\x95\xab\x54\x2b\x97\xec\xc6\x45\x50\xb9\xe4\x00\xfd\xce\x9a\x53\xa9\x1a\xe0\xb6\x92\x35\x73\x11\x0d\xd0\xf5\x9b\x8d\x35\x27\xd5\xb9\x06\x98\x30\x56\xb2\x66\x20\x51\x27\x54\x4b\xbc\x90\xb8\x8c\x89\x02\x44\xd2\x4a\x79\x90\xeb\xac\xff\x0b\x28\x21\xc7\x00\xed\xe7\x87\x92\xac\xc9\xa4\x7f\x8a\x9c\x6f\xd6\xd7\x73\xbd\xe0\x59\x68\xe2\xe7\x00\x2f\x07\x9e\xe3\x72\xe4\x5c\x83\x35\x00\x0e\x54\x90\x26\x3e\x5b\x73\x8b\x5b\xf7\x5f\x27\x89\x52\xeb\x66\xff\xf0\xb0\xdd\xaf\x80\x95\xfe\x3d\xaf\x06\xd8\xad\x60\x46\xc1\x18\x29\x72\x4d\x01\xfa\x95\x54\x8c\xea\x22\xe5\xa3\x9e\x02\x74\xcd\x75\xa6\xeb\x2b\xd9\xf6\xd6\xe0\x91\xf2\xda\x08\x53\xd6\x96\x23\xa4\xc2\xd4\xaa\x68\xb1\xf8\xcc\x69\x49\xee\x37\x6b\x2b\xbf\xdb\x34\xfc\x0b\x87\x73\x99\x26\xc7\x59\x49\x2e\x18\x03\x3c\xac\x58\x39\x51\x59\xb4\xc9\x56\xe1\xdf\x00\x00\x00\xff\xff\x4a\x2a\x97\xb0\x74\x02\x00\x00")

func resourcesSettingsYmlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesSettingsYml,
		"resources/settings.yml",
	)
}

func resourcesSettingsYml() (*asset, error) {
	bytes, err := resourcesSettingsYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/settings.yml", size: 628, mode: os.FileMode(511), modTime: time.Unix(1670057445, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesMigrations1_core_tables_schemaDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8e\x41\x0a\x02\x31\x0c\x45\xf7\x82\x77\xe8\x3d\x5c\x29\x0e\x83\x20\x28\xc5\x7d\xc9\x4c\xb3\x28\xb4\x69\x49\xd2\x01\x6f\x2f\xa3\xc2\xa0\xed\xf6\xfd\x9f\x9f\x77\xb6\xb7\xbb\x79\xd8\xcb\x38\x0e\xd6\xd4\xe2\x41\xd1\x45\x10\x75\x1e\x4b\xcc\xcf\x84\xa4\x87\xfd\xee\xd3\x3a\x9e\xae\x83\xa9\x82\xec\x0a\x72\x0a\x22\x21\x93\x74\x52\x45\x48\x5d\x3e\xfd\xc2\x6d\xa5\x89\x90\x96\xc0\x99\xd6\xf7\x6e\x01\x0e\x30\x45\x6c\x4a\x73\x26\x85\x40\x9d\xe5\x4d\xbe\x89\x28\xfb\x76\x69\x35\x76\x73\xac\xa2\xc8\x7f\xe6\x5f\xda\x3f\x79\xc3\x57\x00\x00\x00\xff\xff\x53\x1f\x9b\x15\x43\x01\x00\x00")

func resourcesMigrations1_core_tables_schemaDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesMigrations1_core_tables_schemaDownSql,
		"resources/migrations/1_core_tables_schema.down.sql",
	)
}

func resourcesMigrations1_core_tables_schemaDownSql() (*asset, error) {
	bytes, err := resourcesMigrations1_core_tables_schemaDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/migrations/1_core_tables_schema.down.sql", size: 323, mode: os.FileMode(511), modTime: time.Unix(1669345184, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesMigrations1_core_tables_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\x4d\x6f\xe2\x30\x10\x3d\x13\x29\xff\x61\x6e\x05\xa9\x5a\x75\x3f\xda\x4b\xd5\x43\x0a\x53\x1a\x2d\x18\x36\x98\xed\xf6\x14\xb9\xc4\xad\xac\xe6\x4b\xb1\x61\xc5\xbf\x5f\x41\x02\x89\x83\x43\x29\x5b\xb4\x7b\x04\xbf\xcc\x8c\xdf\x7b\x33\x9e\xae\x87\x0e\x45\xa0\xce\xed\x00\x41\x71\x16\xf9\xea\x09\xda\xb6\x05\x00\x22\x00\x97\x50\xec\xa3\x07\x63\xcf\x1d\x3a\xde\x23\x7c\xc7\x47\x70\xa6\x74\xe4\x92\xae\x87\x43\x24\xf4\xbc\x40\xf2\x58\x89\x67\xc1\x33\xf8\xe9\x78\xdd\x7b\xc7\x6b\x7f\xbd\xea\x00\x19\x51\x20\xd3\xc1\x20\x07\xc5\x2c\xe2\xdb\xe3\xcb\x8b\xf2\x18\xa6\xc4\xfd\x31\xc5\x1c\x15\x70\x39\xcb\x44\xaa\x44\x12\x03\xc5\x5f\x74\x8b\xb2\xad\xce\xb5\x6d\xd9\x96\x56\xf0\x2c\x9c\x4b\xc5\xb3\x93\xd5\xdc\x58\x8d\xe1\x4a\x5f\x2e\x1b\xae\xa4\x96\x69\x89\xfa\x7c\x51\xcf\xc1\x82\x20\xe3\x52\x6e\x11\x57\xdf\x1a\xe2\xac\xb2\xc9\x94\xcd\xf8\x3e\xe8\x7e\x9a\x38\x8b\x64\x4e\xd4\xe6\xaf\x0a\x61\x7a\x59\x6b\x2f\x34\x9e\x56\xc8\x6d\x97\xa1\xce\x37\x5f\x75\x72\xd4\xdd\xc8\x43\xb7\x4f\x6a\xa8\x0e\x78\x78\x87\x1e\x92\x2e\x4e\x2a\x0a\xb6\x57\x27\x23\x02\x3d\x1c\x20\x45\xe8\x3a\x93\xae\xd3\xc3\xdd\x40\x9b\x14\xd5\x28\x85\x71\xcd\x21\x8c\x94\xc4\x49\xc0\x57\x5f\xfc\x03\xd7\xbc\x49\xfd\x31\xb6\xba\x7c\x8f\xad\x3e\x4c\x1d\x23\xb5\x2c\x4d\x43\x31\x63\xab\xeb\x9f\x8a\xe1\xc3\x18\x7a\x43\x07\x39\x63\x21\x6f\x90\xe0\x2f\xdb\x23\x64\x52\xf9\x01\x4f\xc3\x64\x19\xf1\x58\x01\x75\x87\x38\xa1\xce\x70\xfc\xdf\xf7\x45\x59\xf4\xa9\xb4\xab\xda\xa3\x91\x40\x11\xb1\x17\xee\xcf\x33\x51\x19\x9c\x3b\x93\x33\xe3\x69\x92\xa9\x42\xda\x52\x56\x9e\x2d\xc4\x8c\xfb\x4d\xc5\x94\x40\xc5\xd4\x5c\x56\x6c\x64\xea\x21\x1e\xf8\xac\xa2\x5f\xdd\x26\x49\x94\x86\x5c\xd5\x40\xbb\xda\xe8\x77\xd6\x24\xaa\x75\xcb\x5e\xa5\xea\x63\x3d\x89\x15\x13\x71\xee\x92\x53\x48\x25\xd2\x3d\xe3\xa3\x18\x30\xfe\x33\x8b\x44\xb8\x6c\xd0\xf1\x20\xb1\xd7\xb3\x78\x7f\x2f\xcd\xd3\x80\x9d\x9c\xe6\xdd\x80\x45\x65\x5a\xa4\xcd\xcb\x71\x78\x4f\xf1\x78\x21\xb2\x24\x5e\x37\xd5\x82\x65\x82\x3d\x85\xc7\xbc\x3d\xaf\x7c\x69\xdc\x9f\xf2\xd3\x05\x0b\xe7\xdc\x3c\xe6\x92\xd4\xbc\x78\x15\x93\x8c\x65\x2f\x5c\x19\x5d\x60\xbc\x4d\xca\xb3\x48\x48\x79\xdc\x74\x6f\x5c\x02\x8d\xa9\xe6\xf2\x28\x6f\x6b\x49\x0c\x73\x23\x16\xb3\x57\x5f\x03\x19\xac\x7f\x48\x7f\xf0\x88\x89\x70\x5f\xa6\x94\x49\xf9\x3b\xc9\x82\xc3\xd9\xcd\xaf\x5c\xee\x69\xeb\xdf\x1f\xb0\xa4\x15\x71\xf6\x6d\x68\x05\x44\x33\xfb\x46\x81\x53\xbf\x41\x26\x16\x4a\xa3\x1d\xc4\x45\xc5\x97\xef\x63\x44\xfb\xf0\x34\xbc\xe8\x29\xaa\x81\xf4\x6e\x3a\x68\xa6\x78\x6e\x7f\x75\xb3\x7c\x24\xfa\xf5\x55\xc3\xb9\xa3\xe8\xc1\x74\xdc\x5b\x61\x47\x44\x7f\xd0\x6d\xab\x75\x8b\x7d\x97\xd8\x56\xab\x55\x40\xf4\xc1\x08\xeb\xd2\x01\x60\x82\xb4\xbe\xc5\xdc\xb4\xbb\xce\x04\xe1\xe1\x1e\x09\x10\x7c\xf8\x54\xbc\x9f\x37\x70\xd6\xc3\xf1\x60\xf4\x88\xbd\x33\xa0\x9b\x43\xed\x65\xc4\xc1\x04\x77\x76\x22\x24\xbd\xce\x36\xdd\xc3\x3d\x7a\x08\x22\xb8\x59\x7d\xab\x8f\xf1\x6b\xdb\x6a\x21\xe9\x5d\xdb\xd6\x9f\x00\x00\x00\xff\xff\xa7\x2b\x0e\xb8\xac\x0e\x00\x00")

func resourcesMigrations1_core_tables_schemaUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_resourcesMigrations1_core_tables_schemaUpSql,
		"resources/migrations/1_core_tables_schema.up.sql",
	)
}

func resourcesMigrations1_core_tables_schemaUpSql() (*asset, error) {
	bytes, err := resourcesMigrations1_core_tables_schemaUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "resources/migrations/1_core_tables_schema.up.sql", size: 3756, mode: os.FileMode(511), modTime: time.Unix(1670057445, 0)}
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
	"resources/application.yml":                          resourcesApplicationYml,
	"resources/boot.txt":                                 resourcesBootTxt,
	"resources/permissions.json":                         resourcesPermissionsJson,
	"resources/settings.yml":                             resourcesSettingsYml,
	"resources/migrations/1_core_tables_schema.down.sql": resourcesMigrations1_core_tables_schemaDownSql,
	"resources/migrations/1_core_tables_schema.up.sql":   resourcesMigrations1_core_tables_schemaUpSql,
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
	"resources": &bintree{nil, map[string]*bintree{
		"application.yml": &bintree{resourcesApplicationYml, map[string]*bintree{}},
		"boot.txt":        &bintree{resourcesBootTxt, map[string]*bintree{}},
		"migrations": &bintree{nil, map[string]*bintree{
			"1_core_tables_schema.down.sql": &bintree{resourcesMigrations1_core_tables_schemaDownSql, map[string]*bintree{}},
			"1_core_tables_schema.up.sql":   &bintree{resourcesMigrations1_core_tables_schemaUpSql, map[string]*bintree{}},
		}},
		"permissions.json": &bintree{resourcesPermissionsJson, map[string]*bintree{}},
		"settings.yml":     &bintree{resourcesSettingsYml, map[string]*bintree{}},
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
