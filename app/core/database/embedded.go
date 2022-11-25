package database

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/space-fold-technologies/aurora-service/app/core/server"
)

func init() {
	source.Register("embedded", &Embedded{})
}

type Embedded struct {
	folder     string
	migrations *source.Migrations
}

func (e *Embedded) Open(folder string) (source.Driver, error) {
	instance := &Embedded{folder: folder, migrations: source.NewMigrations()}
	if err := instance.readFolderContent(); err != nil {
		return nil, err
	}
	return instance, nil
}

func WithInstance(folder string) (source.Driver, error) {
	instance := &Embedded{folder: folder, migrations: source.NewMigrations()}
	if err := instance.readFolderContent(); err != nil {
		return nil, err
	}

	return instance, nil
}

func (e *Embedded) readFolderContent() error {
	dir, err := server.AssetDir(e.folder)
	if err != nil {
		return err
	}
	for _, name := range dir {
		migration, err := source.DefaultParse(name)
		if err != nil {
			continue
		}
		if !e.migrations.Append(migration) {
			return fmt.Errorf("failed to parse: %v", name)
		}
	}
	return nil
}

func (e *Embedded) Close() error {
	return nil
}

func (e *Embedded) First() (version uint, err error) {
	if v, ok := e.migrations.First(); !ok {
		return 0, &os.PathError{Op: "first", Path: e.folder, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *Embedded) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := e.migrations.Prev(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: e.folder, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *Embedded) Next(version uint) (nextVersion uint, err error) {
	if v, ok := e.migrations.Next(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: e.folder, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (e *Embedded) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if migration, ok := e.migrations.Up(version); ok {
		content, err := server.Asset(fmt.Sprintf("%s/%s", e.folder, migration.Raw))
		if err != nil {
			return nil, "", err
		}
		return io.NopCloser(bytes.NewReader(content)), migration.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: e.folder, Err: os.ErrNotExist}
}

func (e *Embedded) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if migration, ok := e.migrations.Up(version); ok {
		content, err := server.Asset(fmt.Sprintf("%s/%s", e.folder, migration.Raw))
		if err != nil {
			return nil, "", err
		}
		return io.NopCloser(bytes.NewReader(content)), migration.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: e.folder, Err: os.ErrNotExist}
}
