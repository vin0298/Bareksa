package services_test

import (
	"testing"

	"github.com/bareksa/migration"
)

func TestMain(m *testing.M) {
	// Do the seeding testing
	err := migration.RunMigration()
	if err != nil {
		os.exit(1)
	}
}
