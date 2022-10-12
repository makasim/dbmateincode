package dbmateincode

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/amacneil/dbmate/pkg/dbmate"
)

type Config struct {
	AutoDumpSchema      bool
	DatabaseURL         *url.URL
	MigrationsDir       fs.FS
	TemporaryDir        string
	MigrationsTableName string
	SchemaFile          string
	Verbose             bool
	WaitBefore          bool
	WaitInterval        time.Duration
	WaitTimeout         time.Duration
	Log                 io.Writer
}

func NewConfig(databaseURL *url.URL, migrationsDir fs.FS) Config {
	return Config{
		DatabaseURL:         databaseURL,
		MigrationsDir:       migrationsDir,
		AutoDumpSchema:      true,
		MigrationsTableName: dbmate.DefaultMigrationsTableName,
		SchemaFile:          dbmate.DefaultSchemaFile,
		WaitBefore:          false,
		WaitInterval:        dbmate.DefaultWaitInterval,
		WaitTimeout:         dbmate.DefaultWaitTimeout,
		Log:                 os.Stdout,
	}

}

func Migrate(cfg Config) error {
	var tmpDir string
	if cfg.TemporaryDir != `` {
		tmpDir = cfg.TemporaryDir
	} else {
		var err error
		tmpDir, err = os.MkdirTemp(``, ``)
		if err != nil {
			return fmt.Errorf("mkdir temp: %w", err)
		}
		defer os.RemoveAll(tmpDir)
	}

	if err := fs.WalkDir(cfg.MigrationsDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := cfg.MigrationsDir.Open(path)
		if err != nil {
			return fmt.Errorf("open migration file: %s: %s", path, err)
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("read migration file: %s: %s", path, err)
		}

		tmpPath := filepath.Join(tmpDir, d.Name())

		err = os.WriteFile(tmpPath, b, 0644)
		if err != nil {
			return fmt.Errorf("write tmp migration file: %s: %s", tmpPath, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("migrations dir: walk: %s", err)
	}

	dbMate := &dbmate.DB{
		AutoDumpSchema:      cfg.AutoDumpSchema,
		DatabaseURL:         cfg.DatabaseURL,
		MigrationsDir:       tmpDir,
		MigrationsTableName: cfg.MigrationsTableName,
		SchemaFile:          cfg.SchemaFile,
		Verbose:             cfg.Verbose,
		WaitBefore:          cfg.WaitBefore,
		WaitInterval:        cfg.WaitInterval,
		WaitTimeout:         cfg.WaitTimeout,
		Log:                 cfg.Log,
	}

	if err := dbMate.Migrate(); err != nil {
		return fmt.Errorf("dbmate: migrate: %w", err)
	}

	return nil
}
