package ai

import "github.com/pbreedt/ai-agent/storage"

type Option func(b *Agent)

// WithStorage sets the storage to be used for keeping chat history
func WithStorage(storage storage.Storage) Option {
	return func(b *Agent) { b.storage = storage }
}
