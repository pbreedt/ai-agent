package history

import "github.com/firebase/genkit/go/ai"

type ChatHistory interface {
	GetLast(last int) []*ai.Message
	GetAll() []*ai.Message
	Store(message *ai.Message) error
}
