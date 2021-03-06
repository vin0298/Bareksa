package migration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"

	repo "github.com/bareksa/repository"
)

// This file is to help do initial setups and for testing purposes.

var (
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Dir(b)
)

func GetRootFolderPath() string {
	dir, _ := filepath.Split(Basepath)
	return dir
}

func GetStringFromFile(filename string) (*string, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.New("File not found")
	}

	queryString := fmt.Sprintf("%s", content)

	return &queryString, nil
}

func RunMigration() error {
	newsRepo, err := repo.SetupDatabase()
	if err != nil {
		return err
	}

	basePath := GetRootFolderPath()
	migrationQueryString, err := GetStringFromFile(basePath + "migration/schema.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(*migrationQueryString)

	if err != nil {
		return err
	}

	return nil
}

func RunSeeder() error {
	newsRepo, err := repo.SetupDatabase()
	if err != nil {
		return err
	}

	basePath := GetRootFolderPath()
	seederQueryString, err := GetStringFromFile(basepath + "migration/seeder.sql")

	if err != nil {
		return err
	}

	newsRepo.db.Exec(*seederQueryString)

	return nil
}
