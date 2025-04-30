package storage

import "github.com/firebase/genkit/go/ai"

type Storage interface {
	GetHistory(id ...string) []*ai.Message
	StoreMessage(message *ai.Message, ids ...string) error
}
