package ai

import (
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type ArithmeticInput struct {
	Number1   float64 `json:"num1"`
	Operation string  `json:"operation"`
	Number2   float64 `json:"num2"`
}

func DoBasicArithmeticTool(g *genkit.Genkit) ai.Tool {

	arithmeticTool := genkit.LookupTool(g, "doBasicArithmetic")

	if arithmeticTool != nil {
		return arithmeticTool
	}

	arithmeticTool = genkit.DefineTool(
		g, "doBasicArithmetic", `Do basic arithmetic on two numbers. For example:
		1 + 2 = 3
		3 * 2 = 6
		10 / 5 = 2
		17 - 4 = 13
		`,
		func(ctx *ai.ToolContext, input ArithmeticInput) (string, error) {
			switch input.Operation {
			case "+", "add", "sum", "plus", "increase":
				return fmt.Sprintf("%f + %f = %f", input.Number1, input.Number2, input.Number1+input.Number2), nil
			case "-", "subtract", "minus", "sub", "reduce":
				return fmt.Sprintf("%f - %f = %f", input.Number1, input.Number2, input.Number1-input.Number2), nil
			case "*", "multiply", "times":
				return fmt.Sprintf("%f x %f = %f", input.Number1, input.Number2, input.Number1*input.Number2), nil
			case "/", "divide", "divide by", "divided by":
				return fmt.Sprintf("%f / %f = %f", input.Number1, input.Number2, input.Number1/input.Number2), nil
			}
			return fmt.Sprintf("Sorry, I cannot handle the operation %s, I can only do plus, minus, multiply, and divide on two numbers.", input.Operation), nil
		})

	return arithmeticTool
}
