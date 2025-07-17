package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
	repositoriesDomain "github.com/savio04/youtube-video-summarizer/domains/video/repositories"
	"github.com/savio04/youtube-video-summarizer/internal/database/repositories"
	"github.com/savio04/youtube-video-summarizer/internal/env"
	"github.com/savio04/youtube-video-summarizer/internal/utils"
)

type ResponseTrancription struct {
	Text string `json:"text"`
}

var client redis.Client

func Init() error {
	ctx := context.Background()

	password := env.GetEnvOrDie("POSTGRES_PASSWORD")

	connection := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: password,
		DB:       0,
	})

	_, err := connection.Ping(ctx).Result()
	if err != nil {
		return err
	}

	client = *connection

	go ConsumeQueue(utils.QueueVideoProcessing)

	return nil
}

func InsertIntoQueue(queueName, item string) error {
	ctx := context.Background()

	err := client.LPush(ctx, queueName, item).Err()
	if err != nil {
		return fmt.Errorf("erro ao inserir item na fila: %w", err)
	}
	return nil
}

func ConsumeQueue(queueName string) {
	ctx := context.Background()

	for {
		result, err := client.BLPop(ctx, 0*time.Second, queueName).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("Erro ao consumir fila: %v", err)
			break
		}

		if len(result) > 1 {
			item := result[1]
			fmt.Printf("Processando item: %s\n", item)

			videoRepository := repositories.NewDbVideoRepository()

			fmt.Println("Baixando áudio...")

			newStatus := "DOWNLOADED_AUDIO"

			filePath, err := downloadAudio(item)
			if err != nil {
				newStatus = "FAILED"

				err = videoRepository.UpdateByExternalId(item, &repositoriesDomain.UpdateParams{
					Status:  &newStatus,
					Summary: nil,
				})

				log.Println("Erro ao baixar áudio: ", err)

				continue
			}

			err = videoRepository.UpdateByExternalId(item, &repositoriesDomain.UpdateParams{
				Status:  &newStatus,
				Summary: nil,
			})
			if err != nil {
				log.Println("Erro ao atualizar video: ", err)
				continue
			}

			log.Println("Áudio baixado! Enviando para transcrição...")

			transcription, err := transcribeAudio(*filePath)
			if err != nil {
				newStatus = "FAILED"

				err = videoRepository.UpdateByExternalId(item, &repositoriesDomain.UpdateParams{
					Status:  &newStatus,
					Summary: nil,
				})

				log.Println("Erro na transcrição do áudio: ", err)

				continue
			}

			newStatus = "COMPLETED"

			summary, err := summarizeText(*transcription)
			if err != nil {
				newStatus = "FAILED"

				err = videoRepository.UpdateByExternalId(item, &repositoriesDomain.UpdateParams{
					Status:  &newStatus,
					Summary: nil,
				})

				log.Println("Erro ao fazer resumo do video: ", err)
				continue
			}

			log.Println("Resumo feito com sucesso")

			err = videoRepository.UpdateByExternalId(item, &repositoriesDomain.UpdateParams{Summary: &summary, Status: &newStatus})
			if err != nil {
				log.Println("Erro ao atualizar video: ", err)
				continue
			}

			log.Println("Item " + item + " finalizado!")
		}
	}
}

func downloadAudio(videoId string) (*string, error) {
	filePath := "tmp/" + videoId + ".mp3"
	cmd := exec.Command("yt-dlp", "-x", "--cookies", "cookies.txt", "--audio-format", "mp3", videoId, "-o", filePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &filePath, nil
}

func transcribeAudio(filePath string) (*string, error) {
	fmt.Println("Arquivo sendo enviado:", filePath)

	url := env.GetEnvOrDie("GROQ_BASE_URL")
	token := env.GetEnvOrDie("GROQ_API_KEY")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	defer utils.RemoveFile(filePath)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("model", "whisper-large-v3")
	_ = writer.WriteField("language", "pt")

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url+"/v1/audio/transcriptions", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &ResponseTrancription{}
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}

	response := result.Text

	if response == "" {
		log.Println("transcription is empty")
		return nil, fmt.Errorf("transcription is empty")
	}

	return &response, nil
}

func summarizeText(text string) (string, error) {
	url := env.GetEnvOrDie("GROQ_BASE_URL")
	token := env.GetEnvOrDie("GROQ_API_KEY")

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{"role": "system", "content": "Resuma sempre em português a transcrição do vídeo destacando os principais assuntos abordados. Estruture o resumo em parágrafos bem organizados, separando cada ideia com uma quebra de linha (\n\n)."},
			{"role": "user", "content": text},
		},
		"temperature": 0.5,
	})

	req, err := http.NewRequest("POST", url+"/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	choices := result["choices"].([]interface{})

	return choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string), nil
}
