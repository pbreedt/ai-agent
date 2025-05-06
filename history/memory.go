package history

import (
	"github.com/firebase/genkit/go/ai"
)

type InMemoryHistory struct {
	chatHistory []*ai.Message
	maxHistory  int
}

func NewInMemoryHistory(max int) *InMemoryHistory {
	return &InMemoryHistory{
		maxHistory: max,
	}
}

func (m *InMemoryHistory) GetLast(last int) []*ai.Message {
	if len(m.chatHistory) > last {
		return m.chatHistory[len(m.chatHistory)-last:]
	}
	return m.chatHistory
}

func (m *InMemoryHistory) GetAll() []*ai.Message {
	return m.chatHistory
}

func (m *InMemoryHistory) Store(message *ai.Message) error {
	// log.Println("Storing message:", message.Text)
	if m.maxHistory > 0 && len(m.chatHistory) >= m.maxHistory {
		m.chatHistory = m.chatHistory[1:]
	}
	m.chatHistory = append(m.chatHistory, message)
	return nil
}
