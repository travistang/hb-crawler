package api

import (
	"hb-crawler/rating-gain/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	CredentialsEndpointRoot   = "/credentials"
	CredentialsList           = "/"
	CredentialsCreateEndpoint = "/"
)

type CredentialsApiHandler struct {
	repo *database.DatabaseRepository
}

func (handler *CredentialsApiHandler) credentialsListHandler(c *gin.Context) {
	credentials, err := handler.repo.Login.GetAllAvailableAccounts()
	if err != nil {
		logrus.Warnf("Failed to retrieve data: %+v\n", err)
		reportError(c, http.StatusInternalServerError, "Failed to retrieve data")
		return
	}
	sendJSONPayload(c, http.StatusOK, credentials)
}

func (handler *CredentialsApiHandler) createCredentialHandler(c *gin.Context) {
	var account database.Account
	if err := c.Bind(&account); err != nil {
		reportError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := handler.repo.Login.CreateAccount(&account); err != nil {
		logrus.Warnf("Failed to create account: %+v\n", err)
		reportError(c, http.StatusInternalServerError, "Failed to create account")
		return
	}
	sendJSONPayload(c, http.StatusOK, gin.H{"ok": true})
}

func (handler *CredentialsApiHandler) Register(api *gin.Engine) {
	router := api.Group(CredentialsEndpointRoot)
	router.GET(CredentialsList, handler.credentialsListHandler)
	router.POST(CredentialsCreateEndpoint, handler.createCredentialHandler)
}
