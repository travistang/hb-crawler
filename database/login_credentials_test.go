package database

import (
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	Username = "opulent_umpires0w@icloud.com"
)

func TestCachedCredential(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	db, err := InitializeDatabase("../db.sqlite")
	if err != nil {
		t.Errorf("Failed to initialize database")
		return
	}
	defer db.Close()

	repo := CreateLoginCredentialRepository(db)

	credential := LoginCredentialRecord{}
	if err := repo.GetCredential(&credential, Username, int64(0)); err != nil {
		t.Errorf("Failed to retrieve credential: %+v\n", err)
		return
	}

	if credential.Username != Username {
		t.Errorf("Failed to retrieve credential: \n expected %s\n got %s\n", Username, credential.Username)
	}
}
