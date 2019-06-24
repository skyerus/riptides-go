package GoogleService

import "github.com/skyerus/riptides-go/pkg/google"

type googleService struct {

}

func NewGoogleService() google.Service {
	return &googleService{}
}

