package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/localvec"
)

func InitContacts() (*genkit.Genkit, *ai.Indexer, *ai.Retriever) {
	ctx := context.Background()

	// Used Google AI instead of Vertex AI
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = localvec.Init(); err != nil {
		log.Fatal(err)
	}

	// Also used Google AI here instead of Vertex AI
	indexer, retriever, err := localvec.DefineIndexerAndRetriever(g, "contacts-indexer-retreiever",
		localvec.Config{
			Dir:      "./data/contacts",
			Embedder: googlegenai.GoogleAIEmbedder(g, "text-embedding-004"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return g, &indexer, &retriever
}

func Index(g *genkit.Genkit, indexer *ai.Indexer) *core.Flow[Contact, any, struct{}] {

	f := genkit.DefineFlow(
		g, "contactIndex",
		func(ctx context.Context, contact Contact) (any, error) {
			log.Printf("Indexing contact %s", contact)

			var docs []*ai.Document
			docs = append(docs, ai.DocumentFromText(contact.String(), contact.ToMap()))

			err := ai.Index(ctx, *indexer, ai.WithDocs(docs...))
			if err != nil {
				log.Printf("Error saving index: %v", err)
				return nil, err
			}

			log.Println("Done indexing menu")
			return nil, err
		},
	)

	return f

	// _, err = indexPDFFlow.Run(ctx, "./rag/menu.pdf")
	// if err != nil {
	// 	log.Printf("Error running flow: %v", err)
	// 	log.Fatal(err)
	// }

}

func Retrieve(g *genkit.Genkit, retriever *ai.Retriever) *core.Flow[string, string, struct{}] {

	f := genkit.DefineFlow(
		g, "ragRetrieve",
		func(ctx context.Context, question string) (string, error) {
			// Retrieve text relevant to the user's question.
			resp, err := ai.Retrieve(ctx, *retriever, ai.WithTextDocs(question))

			if err != nil {
				return "", err
			}

			// Call Generate, including the menu information in your prompt.
			return genkit.GenerateText(ctx, g,
				ai.WithModelName("googleai/gemini-2.0-flash"),
				ai.WithDocs(resp.Documents...),
				ai.WithSystem(`
You are acting as a helpful AI assistant that can answer questions about the provided information.
Use only the context provided to answer the question. If you don't know, do not make up an answer.`),
				ai.WithPrompt(question),
			)
		})
	return f

	// res, err := retrieveFlow.Run(ctx, "What are the specials on Monday?")
	// if err != nil {
	// 	log.Printf("Error running flow: %v", err)
	// 	log.Fatal(err)
	// }

	// log.Println(res)
}

func DogFacts() {
	ctx := context.Background()

	// Initialize Genkit
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(
			&googlegenai.GoogleAI{},
		),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatalf("Genkit initialization error: %v", err)
	}

	// Dummy retriever that always returns the same facts
	// Retriever executes each time the flow is run
	dummyRetrieverFunc := func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
		facts := []string{
			"Dog is man's best friend",
			"Dogs have evolved and were domesticated from wolves",
		}
		// Just return facts as documents.
		var docs []*ai.Document
		for _, fact := range facts {
			docs = append(docs, ai.DocumentFromText(fact, nil))
		}
		log.Printf("Retrieved %d dog facts for request: %+v\n", len(docs), req.Query.Content)
		for _, part := range req.Query.Content {
			log.Printf("Part: %s\n", part.Text)
		}
		return &ai.RetrieverResponse{Documents: docs}, nil
	}
	factsRetriever := genkit.DefineRetriever(g, "local", "dogFacts", dummyRetrieverFunc)

	m := googlegenai.GoogleAIModel(g, "gemini-2.0-flash")
	if m == nil {
		log.Fatal("failed to find model")
	}

	// A simple question-answering flow
	genkit.DefineFlow(g, "dogFacts", func(ctx context.Context, query string) (string, error) {
		factDocs, err := ai.Retrieve(ctx, factsRetriever, ai.WithTextDocs(query))
		if err != nil {
			return "", fmt.Errorf("retrieval failed: %w", err)
		}
		llmResponse, err := genkit.Generate(ctx, g,
			ai.WithModelName("googleai/gemini-2.0-flash"),
			ai.WithPrompt("Answer this question with the given context: %s", query),
			ai.WithDocs(factDocs.Documents...),
		)
		if err != nil {
			return "", fmt.Errorf("generation failed: %w", err)
		}
		return llmResponse.Text(), nil
	})
}
