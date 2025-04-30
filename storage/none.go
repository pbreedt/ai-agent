package storage

import "github.com/firebase/genkit/go/ai"

type NoStorage struct{}

func NewNoStorage() Storage {
	return &NoStorage{}
}

func (n *NoStorage) GetHistory(id ...string) []*ai.Message {
	return nil
}

func (n *NoStorage) StoreMessage(message *ai.Message, ids ...string) error {
	return nil
}
