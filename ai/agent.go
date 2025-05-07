package ai

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/pbreedt/ai-agent/calendar"
	"github.com/pbreedt/ai-agent/contacts"
	"github.com/pbreedt/ai-agent/history"
)

// Requires:
// export GEMINI_API_KEY=<key>
type Agent struct {
	genkit      *genkit.Genkit
	chatHistory history.ChatHistory
	calendar    calendar.Calendar
	contactsDB  contacts.ContactsDB
}

func NewAgent(opts ...Option) *Agent {
	// Initialize Genkit with the Google AI plugin and Gemini 2.0 Flash.
	g, err := genkit.Init(context.Background(),
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatalf("could not initialize Genkit: %v", err)
	}

	a := &Agent{
		genkit:      g,
		chatHistory: history.NewNoHistory(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Agent) RespondToPrompt(ctx context.Context, prompt string) (string, error) {

	// dte := time.Now().Format(time.RFC3339)

	resp, err := genkit.Generate(ctx, a.genkit,
		ai.WithSystem(`
		You are acting as a helpful AI chatbot called BreedtBot. You can answer general questions about the provided information.
		You can remember the recent history of the conversation. This chat history is provided to you as context. You can also be asked to forget or remove items from lists.
		You should first attempt to answer the question using one of the provided tools.`),
		//Use only the context provided to answer the question. If you don't know, do not	make up an answer.
		ai.WithPrompt(prompt),
		ai.WithConfig(&googlegenai.GeminiConfig{
			MaxOutputTokens: 500,
		}),
		ai.WithMessages(a.chatHistory.GetLast(10)...), // TODO: Make this configurable
		ai.WithTools(
			DoBasicArithmeticTool(a.genkit),
			GetContactTool(a), StoreContactTool(a),
			GetCalendarEventsTool(a), CreateCalendarEventsTool(a),
		),
	)
	if err != nil {
		log.Println("AI returned error: ", err.Error())
		return "", err
	}

	a.chatHistory.Store(resp.Message)

	return resp.Text(), nil
}
