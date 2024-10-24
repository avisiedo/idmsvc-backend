package datastore

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/podengo-project/idmsvc-backend/internal/config"
)

const dbMirationScriptPath = "./scripts/db/migrations"

// From hmscontent
func CreateMigrationFile(migrationName string) error {
	// datetime format in YYYYMMDDhhmmss - uses the reference time Mon Jan 2 15:04:05 MST 2006
	datetime := time.Now().Format("20060102150405")

	filenameUp := fmt.Sprintf(dbMirationScriptPath+"/%s_%s.up.sql", datetime, migrationName)
	filenameDown := fmt.Sprintf(dbMirationScriptPath+"/%s_%s.down.sql", datetime, migrationName)

	migrationTemplate := fmt.Sprintf(`
-- File created by: %s new %s
BEGIN;
-- your migration here
COMMIT;
`, os.Args[0], migrationName)

	f, err := os.Create(filenameUp)
	if err != nil {
		return err
	}
	_, err = f.WriteString(migrationTemplate)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	f, _ = os.Create(filenameDown)
	if err != nil {
		return err
	}
	_, err = f.WriteString(migrationTemplate)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	return nil
}

func MigrateDb(config *config.Config, direction string, steps int) error {
	if config == nil {
		return fmt.Errorf("'config' cannot be nil")
	}
	_, m, err := NewDbMigration(config)
	if err != nil {
		return err
	}

	// show current database version
	version, dirty, verr := m.Version()
	if verr == nil {
		slog.Info(
			"Current database migration status",
			slog.Uint64("version", uint64(version)),
			slog.Bool("dirty", dirty),
		)
	} else if verr == migrate.ErrNilVersion {
		slog.Info("No database version")
	}

	switch direction {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			steps *= -1
			err = m.Steps(steps)
		} else {
			err = m.Down()
		}
	default:
		return fmt.Errorf("'direction' should be 'up' or 'down' but was found '%s'", direction)
	}

	if err != nil && err == migrate.ErrNoChange {
		slog.Info("No new migrations")
		return nil
	} else if err != nil {
		slog.Error("Error running migration", slog.String("error", err.Error()))
		// Force back to previous migration version. If errors running version 1,
		// drop everything (which would just be the schema_migrations table).
		// This is safe if migrations are wrapped in transaction.
		previousMigrationVersion, err := getPreviousMigrationVersion(m)
		if err != nil {
			return err
		}
		if previousMigrationVersion == 0 {
			if err = m.Drop(); err != nil {
				return err
			}
		} else {
			if err = m.Force(previousMigrationVersion); err != nil {
				return err
			}
		}
	}

	version, dirty, verr = m.Version()
	if verr == nil {
		slog.Info(
			"New database migration status",
			slog.Uint64("version", uint64(version)),
			slog.Bool("dirty", dirty),
		)
	}

	return err

}

func getPreviousMigrationVersion(m *migrate.Migrate) (int, error) {
	var f *os.File
	f, err := os.Open(dbMirationScriptPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	migrationFileNames, _ := f.Readdirnames(0)
	version, _, _ := m.Version()
	var previousMigrationIndex int
	var datetimes []int

	for _, name := range migrationFileNames {
		nameArr := strings.Split(name, "_")
		datetime, _ := strconv.Atoi(nameArr[0])
		datetimes = append(datetimes, datetime)
	}
	previousMigrationIndex = sort.IntSlice(datetimes).Search(int(version)) - 1
	if previousMigrationIndex == -1 {
		return 0, err
	} else {
		return datetimes[previousMigrationIndex], err
	}
}

func MigrateUp(config *config.Config, steps int) error {
	return MigrateDb(config, "up", steps)
}

func MigrateDown(config *config.Config, steps int) error {
	return MigrateDb(config, "down", steps)
}
