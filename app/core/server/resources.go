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

	info := bindataFileInfo{name: "resources/application.yml", size: 224, mode: os.FileMode(511), modTime: time.Unix(1667305046, 0)}
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

	info := bindataFileInfo{name: "resources/boot.txt", size: 1080, mode: os.FileMode(511), modTime: time.Unix(1666820639, 0)}
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

	info := bindataFileInfo{name: "resources/permissions.json", size: 543, mode: os.FileMode(511), modTime: time.Unix(1667603553, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesSettingsYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x90\x4d\x6b\x1b\x31\x10\x86\xef\x02\xfd\x87\x21\x3e\x4b\xd9\x0f\xc7\x9b\xea\xe6\x3a\x86\x96\x90\x6c\xb1\x31\xa1\x47\xb1\x1a\xc7\xa2\x5a\x69\xd1\xcc\x1a\xf6\xdf\x97\xdd\xa6\xc9\xc1\xe8\x22\x1e\x9e\x77\x66\x78\x87\x9c\xae\xde\x61\x36\xf0\xd4\xee\x9e\xf7\x07\x75\x7c\xdb\x1e\x5e\xa4\x18\x72\x3a\xfb\x80\xca\xf9\x8c\x1d\xa7\x3c\x19\xb8\x47\xee\xee\xed\x98\x53\xb6\xb0\x82\x5f\x96\x2f\x90\x22\xcc\x16\xd0\x44\x8c\x3d\x70\x82\xf7\x94\x1c\x10\x8f\xe7\xb3\x14\x0e\xcf\x76\x0c\xac\xba\x30\x12\xcf\x2b\x3e\x00\xac\xe0\xfb\x7e\xb7\x3d\x1d\xf7\xf0\xbb\x3d\xc1\xf1\xb4\x7b\xbe\x91\x55\xf0\xc4\x18\x95\x75\x2e\x23\x91\x81\x42\x2f\xcf\x54\x75\xd3\xdc\xda\xd6\x5d\x31\xb3\x27\xfc\x0a\x94\x85\x6e\xf4\xb7\x8d\x5e\x7f\xd9\x8c\xb6\x37\x60\x5d\xef\xa3\x14\x2e\xf5\xd6\x47\x03\xcd\xa3\x2e\xd7\xa5\xae\xca\x47\x5d\x36\xb5\x8e\x7e\xd0\x3e\xc1\x0a\xfe\x9f\xd7\xbe\xbd\xc2\x53\xfb\xb2\xfd\xf9\x2a\xc5\x12\x55\x84\x44\x3e\x45\xe5\xc6\x6c\xd9\xa7\x68\xa0\x5a\xc3\x0a\x7e\xb4\xa7\xc3\x51\x8a\x4b\x22\xfe\x3c\x57\x8a\x21\x65\x36\x50\x17\x45\x21\x45\x87\x99\x55\x46\x4a\xe1\x3a\xd7\x71\x17\x90\x09\x63\x97\xa7\x81\xef\xa4\xf8\xf8\xcd\x03\xa5\x00\xe8\xb1\x5f\x6a\xdf\x3c\x3c\xd4\x9b\x19\x78\xc6\x7f\xfb\xc8\xc0\x7a\x06\x83\xcd\x36\x04\x0c\x9e\x7a\x03\xd5\x4c\xc8\x06\x56\x01\xe3\x3b\x5f\x0c\x94\x4b\xea\x0f\x4e\x9f\xa4\xae\xfe\x06\x00\x00\xff\xff\x5f\xe0\xf4\x97\xef\x01\x00\x00")

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

	info := bindataFileInfo{name: "resources/settings.yml", size: 495, mode: os.FileMode(511), modTime: time.Unix(1668376007, 0)}
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

	info := bindataFileInfo{name: "resources/migrations/1_core_tables_schema.down.sql", size: 323, mode: os.FileMode(511), modTime: time.Unix(1667609381, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resourcesMigrations1_core_tables_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\x4d\x6f\xe2\x30\x10\x3d\x83\xc4\x7f\x98\x5b\x41\xaa\x56\xdd\x8f\xf6\x52\xf5\x90\xc2\x94\x46\x0b\x86\x0d\x66\xbb\x3d\x45\x2e\x71\x2b\xab\xf9\x52\x6c\x58\xf1\xef\x57\x21\x81\xc4\xc1\xa1\x29\x5b\xb4\x7b\x04\xbf\xcc\x8c\xdf\x7b\x33\x9e\xbe\x83\x16\x45\xa0\xd6\xed\x08\x41\x71\x16\xb8\xea\x09\xba\x9d\x36\x00\x08\x0f\x6c\x42\x71\x88\x0e\x4c\x1d\x7b\x6c\x39\x8f\xf0\x1d\x1f\xc1\x9a\xd3\x89\x4d\xfa\x0e\x8e\x91\xd0\xf3\x1c\xc9\x43\x25\x9e\x05\x4f\xe0\xa7\xe5\xf4\xef\x2d\xa7\xfb\xf5\xaa\x07\x64\x42\x81\xcc\x47\xa3\x0c\x14\xb2\x80\xef\x8e\x2f\x2f\x8a\x63\x98\x13\xfb\xc7\x1c\x33\x94\xc7\xe5\x22\x11\xb1\x12\x51\x08\x14\x7f\xd1\x1d\xaa\xd3\xee\x5d\x77\xda\x9d\xb6\x56\xf0\xc2\x5f\x4a\xc5\x93\x93\xd5\x5c\x5b\x8d\xe1\x4a\x5f\x2e\x6b\xae\xa4\xd6\x71\x81\xfa\x7c\x51\xcd\xc1\x3c\x2f\xe1\x52\xee\x10\x57\xdf\xea\xe2\x44\xaf\x3c\x2c\xa5\xbb\xaa\xc1\xa5\x55\xc9\x98\x2d\xf8\xa1\x90\x87\xe9\xe4\x2c\x90\x19\xa1\xdb\xbf\x4a\xc4\xea\xe5\x6f\x3c\x53\x7b\x5a\x12\xa1\x5b\x84\x3a\xdf\x7e\xd5\xcb\x50\x77\x13\x07\xed\x21\xa9\xa0\x7a\xe0\xe0\x1d\x3a\x48\xfa\x38\x2b\x29\xdd\x4d\x4f\x26\x04\x06\x38\x42\x8a\xd0\xb7\x66\x7d\x6b\x80\xfb\x81\xb6\x29\xca\x51\x72\x83\x9b\x43\x18\x29\x09\x23\x8f\xa7\x5f\xfc\x03\x77\xbd\x49\xfd\x31\xf6\xbb\x7c\x8f\xfd\x3e\x4c\x1d\x23\xb5\x2c\x8e\x7d\xb1\x60\xe9\xf5\x4f\xc5\x70\x33\x86\xde\xd0\x41\x2e\x98\xcf\x6b\x24\xf8\xcb\xf6\xf0\x99\x54\xae\xc7\x63\x3f\x5a\x07\x3c\x54\x40\xed\x31\xce\xa8\x35\x9e\xfe\xf7\x7d\x51\x14\x7d\x2a\xed\xca\xf6\xa8\x25\x50\x04\xec\x85\xbb\xcb\x44\x94\x06\xec\xde\x84\x4d\x78\x1c\x25\x2a\x97\xb6\x90\x95\x27\x2b\xb1\xe0\x6e\x5d\x31\x05\x50\x31\xb5\x94\x25\x1b\x99\x7a\x88\x7b\x2e\x2b\xe9\x57\xb5\x49\x14\xc4\x3e\x57\x15\xd0\xbe\x36\xfa\x9d\x35\x89\x2a\xdd\x72\x50\xa9\xea\x58\x8f\x42\xc5\x44\x98\xb9\xe4\x14\x52\x89\xf8\xc0\xf8\xc8\x07\x8c\xfb\xcc\x02\xe1\xaf\x6b\x74\x6c\x24\xf6\x66\x16\xd7\x9e\x7e\x20\x8f\xfb\x01\xf3\xd4\x5a\xa4\xed\xd3\xd0\xbc\x69\x78\xb8\x12\x49\x14\x6e\xba\x66\xc5\x12\xc1\x9e\xfc\x63\x1e\x97\x57\xbe\x36\x2e\x52\xd9\xe9\x8a\xf9\x4b\x6e\x9e\x63\x51\x6c\xde\xc0\xf2\x51\xc5\x92\x17\xae\x8c\x32\x1b\x6f\x13\xf3\x24\x10\x52\x1e\x37\xbe\x6b\xb7\x41\x63\xaa\xa5\x3c\xca\xbc\x5a\x12\xc3\x60\x08\xc5\xe2\xd5\xd5\x40\x06\x6f\x37\x69\x00\x1e\x30\xe1\x1f\xca\x14\x33\x29\x7f\x47\x89\xd7\x9c\xdd\xec\xca\xc5\x22\xb6\xf9\xfd\x01\x5b\x58\x1e\xe7\xd0\x0a\x96\x43\x34\xb3\x6f\x15\x38\xf5\x23\x63\x62\xa1\x30\x5a\x23\x2e\x4a\xbe\x7c\x1f\x23\xda\x87\xa7\xe1\x45\x4f\x51\x0e\xa4\x77\x53\xa3\x99\xe2\xd8\xc3\xf4\x66\xcb\xd8\x63\x8a\xbb\xd5\x5d\xc2\xba\xa3\xe8\xc0\x7c\x3a\x48\xb1\x13\xa2\xbf\xd8\x9d\x76\xeb\x16\x87\x36\xe9\xb4\x5b\xad\x1c\xa2\x0f\x46\xd8\x94\x0e\x00\x33\xa4\xd5\x35\xe5\xa6\xdb\xb7\x66\x08\x0f\xf7\x48\x80\xe0\xc3\xa7\xfc\x81\xbc\x81\xb3\x01\x4e\x47\x93\x47\x1c\x9c\x01\xdd\x1e\x6a\x4f\x1f\x8e\x66\xb8\xb7\xf4\x20\x19\xf4\x76\xe9\x1e\xee\xd1\x41\x10\xde\x4d\xfa\xad\x3e\xc6\xaf\x3b\xed\x16\x92\xc1\xf5\x9f\x00\x00\x00\xff\xff\x48\xef\x98\xc8\xb3\x0e\x00\x00")

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

	info := bindataFileInfo{name: "resources/migrations/1_core_tables_schema.up.sql", size: 3763, mode: os.FileMode(511), modTime: time.Unix(1668378496, 0)}
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
