package utils

import (
	"fmt"
	"os"
)

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Erro ao deletar o arquivo:", err)
		return err
	}

	fmt.Println("Arquivo temporario deletado com sucesso!")

	return nil
}
