package ai

import "github.com/pbreedt/ai-agent/history"

type Option func(b *Agent)

// WithHistory sets the storage to be used for keeping chat history
func WithHistory(history history.History) Option {
	return func(b *Agent) { b.history = history }
}
