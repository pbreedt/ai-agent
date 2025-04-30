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
			// log.Println("AI got prompt:", p)
			resp := a.RespondToPrompt(ctx, p)
			// log.Println("AI response:", resp)
			responseChan <- resp
		case <-ctx.Done():
			// log.Println("Context done:", ctx.Err())
			responseChan <- "AI shut down"
			return
		}
	}
}

func (a *Agent) RespondToPrompt(ctx context.Context, prompt string) string {

	resp, err := genkit.Generate(ctx, a.genkit,
		ai.WithPrompt(prompt),
		ai.WithConfig(&googlegenai.GeminiConfig{
			MaxOutputTokens: 500,
		}),
		ai.WithMessages(a.storage.GetHistory()...),
		ai.WithTools(GetWeatherTool(a.genkit), DoBasicArithmeticTool(a.genkit)),
	)
	if err != nil {
		log.Println("AI returned error: ", err.Error())
		return err.Error()
	}

	a.storage.StoreMessage(resp.Message)

	return resp.Text()
}
