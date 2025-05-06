package history

import "github.com/firebase/genkit/go/ai"

type History interface {
	GetLast(last int) []*ai.Message
	GetAll() []*ai.Message
	Store(message *ai.Message) error
}
