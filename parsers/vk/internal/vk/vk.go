package vk

import (
	"log"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/longpoll-bot"
)

type VKClient struct {
	Bot *longpoll.LongPoll
}

func New() *VKClient {
	token := "vk1.a.BUs5savsnH-xsEg8AZvQaxHODW9wxBMDMUa5Hkc5U9dV-WYXDQl6ruhxub949tZenRAJVwHDex8nJ6apPPZAiRB2W7WcoUYWoUAQvFkb-HwVQZJpQldW1Tsc31_qoqLnz8ff15OOXcJpCufaBelfNNPca0fDaZ9SnQE0qeMKm5EGi8Mr29t8dr11q5vOGcyLW5_KBP4i7KzWNWNfnzpCjQ" // Токен из переменной окружения

	// Создаем объект API с токеном
	vk := api.NewVK(token)

	// Получаем информацию о группе, чтобы узнать ID
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatalf("Ошибка получения информации о группе: %v", err)
	}
	groupID := group.Groups[0].ID

	// Создаем объект Long Poll
	lp, err := longpoll.NewLongPoll(vk, groupID)
	if err != nil {
		log.Fatalf("Ошибка создания Long Poll: %v", err)
	}

	return &VKClient{
		Bot: lp,
	}
}
