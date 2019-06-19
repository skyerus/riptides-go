package handler

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler interface {
	SaveCredentials(response *http.Response, user *models.User) customError.Error
	HandleAuthorizedRequest(r *http.Request, user *models.User) customError.Error
	GetRefreshRequest(user *models.User) (*http.Request, customError.Error)
}

type RequestHandler interface {
	SendRequest(request *http.Request, user *models.User, authorized bool, allowRecursion bool) (*http.Response, customError.Error)
	sendRefreshRequest(client *http.Client, user *models.User) customError.Error
}

type requestHandler struct {
	Handler
}

func NewRequestHandler(handler Handler) RequestHandler {
	return &requestHandler{handler}
}

func (handler requestHandler) SendRequest(request *http.Request, user *models.User, authorized bool, allowRecursion bool) (*http.Response, customError.Error) {
	client := &http.Client{}

	if authorized {
		handler.Handler.HandleAuthorizedRequest(request, user)
	}

	response, err := client.Do(request)
	if err != nil {
		return response, customError.NewGenericHttpError(err)
	}
	if response.StatusCode == http.StatusUnauthorized {
		if authorized && allowRecursion {
			customErr := handler.sendRefreshRequest(client, user)
			if customErr != nil {
				return response, customErr
			}
			return handler.SendRequest(request, user, true, false)
		}
		return response, customError.NewUnauthorizedError(nil)
	}

	if response.StatusCode >= 300 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response, customError.NewGenericHttpError(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return response, customError.NewGenericHttpError(nil)
	}

	return response, nil
}

func (handler requestHandler) sendRefreshRequest(client *http.Client, user *models.User) customError.Error {
	request, customErr := handler.Handler.GetRefreshRequest(user)
	if customErr != nil {
		return customErr
	}
	response, err := client.Do(request)
	defer request.Body.Close()
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	handler.Handler.SaveCredentials(response, user)

	return nil
}
