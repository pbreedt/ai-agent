package history

import (
	"fmt"
	"testing"

	"github.com/firebase/genkit/go/ai"
)

func TestHistory1(t *testing.T) {
	m := NewInMemoryHistory(10)
	err := m.Store(ai.NewUserTextMessage("hello"))
	if err != nil {
		t.Error(err)
	}

	msgs := m.GetLast(10)
	if len(msgs) != 1 {
		t.Errorf("expected 1 message, got %d", len(msgs))
	}

	if msgs[0].Content[0].Text != "hello" {
		t.Errorf("expected 'hello', got %s", msgs[0].Content[0].Text)
	}
}

// Test getting the last 10 messages, even if there are more
func TestHistory15(t *testing.T) {
	m := NewInMemoryHistory(10)
	for i := range 15 {
		err := m.Store(ai.NewUserTextMessage(fmt.Sprintf("hello %d", i)))
		if err != nil {
			t.Error(err)
		}
	}

	msgs := m.GetLast(10)
	if len(msgs) != 10 {
		t.Errorf("expected 10 message, got %d", len(msgs))
	}

	if msgs[0].Content[0].Text != "hello 5" {
		t.Errorf("expected 'hello 5', got %s", msgs[0].Content[0].Text)
	}

	if msgs[9].Content[0].Text != "hello 14" {
		t.Errorf("expected 'hello 14', got %s", msgs[9].Content[0].Text)
	}
}

// If max history is set to 10, then no more than 10 messages should be stored
func TestMaxHistory(t *testing.T) {
	m := NewInMemoryHistory(10)
	for i := range 15 {
		err := m.Store(ai.NewUserTextMessage(fmt.Sprintf("hello %d", i)))
		if err != nil {
			t.Error(err)
		}
	}

	// return only 10 even if 15 were stored and all messages were requested
	msgs := m.GetAll()
	if len(msgs) != 10 {
		t.Errorf("expected 10 message, got %d", len(msgs))
	}

	if msgs[0].Content[0].Text != "hello 5" {
		t.Errorf("expected 'hello 5', got %s", msgs[0].Content[0].Text)
	}

	if msgs[9].Content[0].Text != "hello 14" {
		t.Errorf("expected 'hello 14', got %s", msgs[9].Content[0].Text)
	}
}
