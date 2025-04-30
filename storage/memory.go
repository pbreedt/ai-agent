package storage

import (
	"github.com/firebase/genkit/go/ai"
)

type MemoryStorage struct {
	chatHistory []*ai.Message
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{}
}

func (m *MemoryStorage) GetHistory(id ...string) []*ai.Message {
	return m.chatHistory
}

func (m *MemoryStorage) StoreMessage(message *ai.Message, ids ...string) error {
	// log.Println("Storing message:", message.Text)
	m.chatHistory = append(m.chatHistory, message)
	return nil
}
