package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/savio04/youtube-video-summarizer/internal/env"
)

type VerifyTokenUseCase struct{}

func NewVerifyTokenUseCase() *VerifyTokenUseCase {
	return &VerifyTokenUseCase{}
}

func (useCase *VerifyTokenUseCase) Execute(token string) error {
	const baseUrl = env.GetEnvOrDie("RECAPTCHA_BASE_URL")
	const secret = env.GetEnvOrDie("RECAPTCHA_SECRET_TOKEN")

	requestBody, _ := json.Marshal(map[string]interface{}{
		"secret":   "6LfifNcqAAAAAFb8HNuvB9CpGOKPVkKgNAuvz3Do",
		"response": token,
	})

	data := url.Values{}
	data.Set("username", "teste")
	data.Set("password", "123456")

	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("result", result)

	if !result["success"].(bool) {
		return fmt.Errorf("invalid token")
	}

	return nil
}
