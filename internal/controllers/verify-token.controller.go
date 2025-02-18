package controllers

import (
	"encoding/json"
	"net/http"

	usecases "github.com/savio04/youtube-video-summarizer/domains/capatcha/useCases"
)

type VerifyTokenController struct {
	verifyTokenUseCase *usecases.VerifyTokenUseCase
}

func NewVerifyTokenController() *VerifyTokenController {
	verifyTokenUseCase := usecases.NewVerifyTokenUseCase()

	return &VerifyTokenController{
		verifyTokenUseCase,
	}
}

type VerifyTokenPayload struct {
	Token string `json:"token" validate:"required"`
}

func (verifyTokenController *VerifyTokenController) Handler(writer http.ResponseWriter, request *http.Request) {
	var body VerifyTokenPayload

	decoder := json.NewDecoder(request.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		http.Error(writer, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	err := verifyTokenController.verifyTokenUseCase.Execute(body.Token)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(writer).Encode(struct {
			Message string `json:"message"`
		}{Message: err.Error()})

		return
	}

	writer.WriteHeader(http.StatusOK)
}
