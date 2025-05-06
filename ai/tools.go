package ai

import (
	"fmt"
	"log"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type CalendarInput struct {
	From time.Time `json:"from-datetime"`
	To   time.Time `json:"to-datetime"`
}

func GetCalendarEventsTool(a *Agent) ai.Tool {

	getCalEvents := genkit.LookupTool(a.genkit, "getCalEvents")

	if getCalEvents != nil {
		return getCalEvents
	}

	getCalEvents = genkit.DefineTool(
		a.genkit, "getCalEvents", "Gets the calendar events within a given time range. If no time range is provided, then assume start time 00:00 on the start date to end time 24:00 on the end date.",
		func(ctx *ai.ToolContext, input CalendarInput) (string, error) {
			if input.From.IsZero() {
				return "Please provide at least a start date.", nil
			}
			if input.To.IsZero() {
				input.To = input.From.Add(time.Hour * 24)
			}

			log.Printf("[GetCalendarEventsTool] Getting events from %s to %s\n", input.From, input.To)
			log.Printf("[GetCalendarEventsTool] Agent: %v, Calendar: %v\n", a, a.calendar)

			events, err := a.calendar.GetEvents(input.From, input.To)
			if err != nil {
				return fmt.Sprintf("The calendar entries could not be found due to the following error: %s", err.Error()), err
			}

			var eventStrings []string
			for _, event := range events {
				eventStrings = append(eventStrings, event.String())
			}
			return fmt.Sprintf("The calendar entries from %s to %s are: %s", input.From, input.To, eventStrings), nil
		})

	return getCalEvents
}

type WeatherInput struct {
	Location string `json:"location"`
}

func GetWeatherTool(g *genkit.Genkit) ai.Tool {

	getWeatherTool := genkit.LookupTool(g, "getWeather")

	if getWeatherTool != nil {
		return getWeatherTool
	}

	getWeatherTool = genkit.DefineTool(
		g, "getWeather", "Gets the current weather in a given location",
		func(ctx *ai.ToolContext, input WeatherInput) (string, error) {
			// Here, we would typically make an API call or database query.
			// For this example, we just return a fixed value.
			return fmt.Sprintf("The current weather in %s is 63°F and sunny.", input.Location), nil
		})

	return getWeatherTool
}

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
