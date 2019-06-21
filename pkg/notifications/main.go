package notifications

import (
	"bytes"
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"log"
	"net/http"
	"os"
)

func PushNotification(token string, message string) customError.Error {
	client := &http.Client{}
	notification := models.Notification{Notification: message}

	bodyBytes, err := json.Marshal(notification)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	b := bytes.NewBuffer(bodyBytes)

	request, err := http.NewRequest("POST", os.Getenv("CHAT_SOCKET_URL") + "/socket/push", b)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	if response.StatusCode >= 300 {
		log.Println("Bad response from notifications")
		return customError.NewHttpError(http.StatusInternalServerError, "Bad response from notifications", nil)
	}

	return nil
}
