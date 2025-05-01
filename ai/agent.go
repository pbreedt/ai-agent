package ai

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/pbreedt/ai-agent/storage"
)

// Requires:
// export GEMINI_API_KEY=<key>
type Agent struct {
	genkit  *genkit.Genkit
	storage storage.Storage
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
		genkit:  g,
		storage: storage.NewNoStorage(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Read from prompt channel and respond
func (a *Agent) Run(ctx context.Context, promptChan <-chan string, responseChan chan<- string) {
	defer close(responseChan)

	for {
		select {
		case p := <-promptChan:
			log.Println("AI got prompt:", p)
			resp, err := a.RespondToPrompt(ctx, p)
			if err != nil {
				resp = err.Error()
			}
			log.Println("AI response:", resp)
			responseChan <- resp
		case <-ctx.Done():
			log.Println("Context done:", ctx.Err())
			responseChan <- "AI shut down"
			return
		}
	}
}

func (a *Agent) RespondToPrompt(ctx context.Context, prompt string) (string, error) {

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
		ai.WithMessages(a.storage.GetHistory()...),
		ai.WithTools(GetWeatherTool(a.genkit), DoBasicArithmeticTool(a.genkit)),
	)
	if err != nil {
		log.Println("AI returned error: ", err.Error())
		return "", err
	}

	a.storage.StoreMessage(resp.Message)

	return resp.Text(), nil
}
