package queue

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/redis/go-redis/v9"
	"github.com/savio04/youtube-video-summarizer/internal/env"
	"github.com/savio04/youtube-video-summarizer/internal/utils"
)

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
			//
			// videoRepository := repositoriesdb.NewDbVideoRepository()
			//
			// video, err := videoRepository.FindOne(&repositories.FindOneVideoParams{
			// 	ExternalId: &item,
			// })
			// if err != nil {
			// 	fmt.Printf("Erro ao ler video %v", err)
			// 	continue
			// }

			fmt.Println("Baixando áudio...")

			err = downloadAudio(item)
			if err != nil {
				log.Println("Erro ao baixar áudio: ", err)
				continue
			}

			fmt.Println("Áudio baixado! Enviando para transcrição...")
		}
	}
}

func downloadAudio(videoId string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoId)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels()

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}

	defer stream.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}

	return nil
}

// func transcribeAudio(audioData []byte) (string, error) {
// 	apiKey := os.Getenv("GROQ_API_KEY")
// 	if apiKey == "" {
// 		return "", fmt.Errorf("a variável de ambiente GROQ_API_KEY não está definida")
// 	}
//
// 	client := openai.NewClient(apiKey)
// 	ctx := context.Background()
//
// 	// Criar um arquivo temporário para armazenar o áudio
// 	tempFile, err := os.CreateTemp("", "audio-*.mp3")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer os.Remove(tempFile.Name()) // Excluir após uso
//
// 	_, err = tempFile.Write(audioData)
// 	if err != nil {
// 		return "", err
// 	}
// 	tempFile.Close()
//
// 	// Criar a requisição para transcrição
// 	req := openai.AudioRequest{
// 		Model:       "whisper-large-v3",
// 		FilePath:    tempFile.Name(),
// 		Prompt:      "Specify context or spelling",
// 		Format:      "json",
// 		Language:    "en",
// 		Temperature: 0.0,
// 	}
//
// 	resp, err := client.CreateTranscription(ctx, req)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return resp.Text, nil
// }
