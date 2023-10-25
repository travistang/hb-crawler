package database

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestDecrypt(t *testing.T) {
	if _, err := createCipher(); err != nil {
		log.Errorf("Failed to create cipher: %+v\n", err)
		t.Fail()
	}
}
