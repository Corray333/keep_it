package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/Corray333/keep_it/parsers/telegram/internal/entities"
	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	DB    *sqlx.DB
	Kafka sarama.SyncProducer
	redis *redis.Client
}

func New() *Storage {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_NAME"))
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	// kafka
	// TODO: add sasl
	brokerList := []string{"kafka:9092"}

	// Configure Sarama
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new synchronous producer
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err = redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return &Storage{
		DB:    db,
		Kafka: producer,
		redis: redis,
	}
}

func (s *Storage) NewNote(ctx context.Context, note *entities.NewNoteMessage) error {

	// Define the topic and message
	topic := "newNotes"

	encoded, err := json.Marshal(note)
	if err != nil {
		return err
	}

	// Produce a message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encoded),
	}

	fmt.Println()
	fmt.Println("Note: ", note)
	fmt.Println()
	_, _, err = s.Kafka.SendMessage(msg)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to produce message: "+err.Error())
		return err
	}

	return nil
}

func (s *Storage) SaveNote(ctx context.Context, creationDate, chatID int64, note *entities.Note) error {
	fmt.Println("Zametka: ", note.ContentDecoded)
	jsonData, err := json.Marshal(note)
	if err != nil {
		slog.Error("failed to marshal note: ", "error", err)
		return err
	}

	date := strconv.Itoa(int(creationDate))
	chat := strconv.Itoa(int(chatID))

	return s.redis.RPush(ctx, date+chat, jsonData).Err()
}

func (s *Storage) GetNotes(ctx context.Context, creationDate, chtID int64) ([]*entities.Note, error) {
	date := strconv.Itoa(int(creationDate))
	chat := strconv.Itoa(int(chtID))

	pipe := s.redis.TxPipeline()

	notesCmd := pipe.LRange(ctx, date+chat, 0, -1)

	delCmd := pipe.Del(ctx, date+chat)

	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.Error("failed to get notes: ", "error", err)
		return nil, err
	}

	notes, err := notesCmd.Result()
	if err != nil {
		slog.Error("failed to get notes: ", "error", err)
		return nil, err
	}

	var res []*entities.Note
	for _, note := range notes {
		var n *entities.Note
		if err := json.Unmarshal([]byte(note), &n); err != nil {
			return nil, err
		}
		fmt.Println("Anal: ", n.ContentDecoded)

		res = append(res, n)
	}

	if del, err := delCmd.Result(); err != nil {
		slog.Error("failed to delete notes: ", "error", err)
		return nil, err
	} else if del == 0 {
		return nil, nil
	}

	return res, nil
}
