package history

import "github.com/firebase/genkit/go/ai"

type NoHistory struct{}

func NewNoHistory() History {
	return &NoHistory{}
}

func (n *NoHistory) GetLast(last int) []*ai.Message {
	return nil
}

func (n *NoHistory) GetAll() []*ai.Message {
	return nil
}

func (n *NoHistory) Store(message *ai.Message) error {
	return nil
}
