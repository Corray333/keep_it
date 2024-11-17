package file

import (
	"io"
	"log/slog"
	"os"
	"time"

	"math/rand"
)

type FileManager struct{}

func New() *FileManager {
	return &FileManager{}
}

func generateRandomString(length int) string {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}
	return string(b)
}

func (f *FileManager) SaveImage(img io.Reader, name string) (string, error) {
	filePath := os.Getenv("FILE_PATH") + "/images/" + name + generateRandomString(10) + ".png"
	newFile, err := os.Create(filePath)
	if err != nil {
		slog.Error("Failed to create file: " + err.Error())
		return "", err
	}

	defer newFile.Close()

	_, err = io.Copy(newFile, img)
	if err != nil {
		slog.Error("Failed to copy image: " + err.Error())
		return "", err
	}

	return filePath, nil
}
