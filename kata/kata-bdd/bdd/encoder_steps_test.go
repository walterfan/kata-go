package bdd

import (
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/waltfy/kata-bdd/encoder/pkg/encoder"
)

type encoderWorld struct {
	text   string
	result string
	Err    error
}

func (w *encoderWorld) iHaveTheText(txt string) error {
	w.text = txt
	return nil
}

func (w *encoderWorld) iEncodeUsing(t string) error {
	switch t {
	case "base64":
		w.result = encoder.EncodeBase64(w.text)
	case "hex":
		w.result = encoder.EncodeHex(w.text)
	case "url":
		w.result = encoder.EncodeURL(w.text)
	default:
		w.Err = fmt.Errorf("unknown type: %s", t)
	}
	return nil
}

func (w *encoderWorld) iDecodeUsing(t string) error {
	switch t {
	case "base64":
		w.result, w.Err = encoder.DecodeBase64(w.text)
	case "hex":
		w.result, w.Err = encoder.DecodeHex(w.text)
	case "url":
		w.result, w.Err = encoder.DecodeURL(w.text)
	default:
		w.Err = fmt.Errorf("unknown type: %s", t)
	}
	return nil
}

func (w *encoderWorld) theResultShouldBe(expected string) error {
	if w.Err != nil {
		return w.Err
	}
	if w.result != expected {
		return fmt.Errorf("expected %q, got %q", expected, w.result)
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	world := &encoderWorld{}

	sc.Step(`^I have the text "([^"]*)"$`, world.iHaveTheText)
	sc.Step(`^I encode using "([^"]*)"$`, world.iEncodeUsing)
	sc.Step(`^I decode using "([^"]*)"$`, world.iDecodeUsing)
	sc.Step(`^the result should be "([^"]*)"$`, world.theResultShouldBe)
}

func TestGodog(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "encoder",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  []string{"."},
			Strict: true,
		},
	}

	if suite.Run() != 0 {
		t.Fail()
	}
}
