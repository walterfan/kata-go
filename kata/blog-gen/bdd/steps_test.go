package bdd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

type testWorld struct {
	llmServer *httptest.Server
	cmdOutput string
	cmdErr    error
	exitCode  int
	today     string
	blogPath  string
}

func (w *testWorld) reset() {
	w.cmdOutput = ""
	w.cmdErr = nil
	w.exitCode = 0
	w.today = time.Now().Format("2006-01-02")
	w.blogPath = filepath.Join(".", fmt.Sprintf("blog-%s.md", w.today))
	_ = os.Remove(w.blogPath)
}

func (w *testWorld) aValidEnvironmentWithAPIKey(apiKey string) error {
	os.Setenv("LLM_API_KEY", apiKey)
	os.Setenv("LLM_BASE_URL", "http://127.0.0.1:18080/v1")
	os.Setenv("LLM_MODEL", "gpt-4o-mini")
	os.Unsetenv("WEATHER_API_URL")
	return nil
}

func (w *testWorld) anEnvironmentWithoutAPIKey() error {
	os.Unsetenv("LLM_API_KEY")
	os.Setenv("LLM_BASE_URL", "http://127.0.0.1:18080/v1")
	os.Setenv("LLM_MODEL", "gpt-4o-mini")
	os.Unsetenv("WEATHER_API_URL")
	return nil
}

func (w *testWorld) aMockLLMServerIsRunning() error {
	w.llmServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasSuffix(req.URL.Path, "/chat/completions") {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if os.Getenv("LLM_API_KEY") == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		// Return a minimal OpenAI-like response
		response := `{"choices":[{"index":0,"message":{"role":"assistant","content":"# my blog at ` + w.today + `\n\n## \"WebRTC recorder with Pion\"\n\nWeather: Sunny"},"finish_reason":"stop"}]}`
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(response))
	}))
	// point base url to mock
	os.Setenv("LLM_BASE_URL", w.llmServer.URL)
	return nil
}

func (w *testWorld) weatherAPIIsNotConfigured() error {
	os.Unsetenv("WEATHER_API_URL")
	return nil
}

func (w *testWorld) iRun(cmdline string) error {
	parts := strings.Fields(cmdline)
	cmd := exec.Command(parts[0], parts[1:]...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	w.cmdOutput = outBuf.String() + errBuf.String()
	w.cmdErr = err
	if exitErr, ok := err.(*exec.ExitError); ok {
		w.exitCode = exitErr.ExitCode()
	} else if err != nil {
		w.exitCode = 1
	}
	return nil
}

func (w *testWorld) aBlogMarkdownFileForTodayShouldBeCreated() error {
	if _, err := os.Stat(w.blogPath); err != nil {
		return fmt.Errorf("expected blog file not found: %s", w.blogPath)
	}
	return nil
}

func (w *testWorld) theBlogContentShouldIncludeIdea(idea string) error {
	b, err := os.ReadFile(w.blogPath)
	if err != nil {
		return err
	}
	if !strings.Contains(string(b), idea) {
		return fmt.Errorf("blog does not contain idea: %s", idea)
	}
	return nil
}

func (w *testWorld) theBlogContentShouldInclude(text string) error {
	b, err := os.ReadFile(w.blogPath)
	if err != nil {
		return err
	}
	if !strings.Contains(string(b), text) {
		return fmt.Errorf("blog does not contain: %s", text)
	}
	return nil
}

func (w *testWorld) theBlogTitleShouldBe(expected string) error {
	b, err := os.ReadFile(w.blogPath)
	if err != nil {
		return err
	}
	if expected == "my blog at <today>" {
		expected = fmt.Sprintf("my blog at %s", w.today)
	}
	if !strings.Contains(string(b), "# "+expected) {
		return fmt.Errorf("blog title mismatch")
	}
	return nil
}

func (w *testWorld) theProcessShouldExitNonZero() error {
	if w.exitCode == 0 {
		return fmt.Errorf("expected non-zero exit code, got 0")
	}
	return nil
}

func (w *testWorld) theErrorOutputShouldContain(msg string) error {
	if !strings.Contains(w.cmdOutput, msg) {
		return fmt.Errorf("expected error output to contain %q, got %q", msg, w.cmdOutput)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	w := &testWorld{}

	ctx.Before(func(context.Context, *godog.Scenario) (context.Context, error) {
		w.reset()
		return context.Background(), nil
	})
	ctx.After(func(context.Context, *godog.Scenario, error) (context.Context, error) {
		if w.llmServer != nil {
			w.llmServer.Close()
		}
		return context.Background(), nil
	})

	ctx.Step(`^a valid environment with API key "([^"]*)"$`, w.aValidEnvironmentWithAPIKey)
	ctx.Step(`^an environment without API key$`, w.anEnvironmentWithoutAPIKey)
	ctx.Step(`^a mock LLM API server is running$`, w.aMockLLMServerIsRunning)
	ctx.Step(`^weather API is not configured$`, w.weatherAPIIsNotConfigured)
	ctx.Step(`^I run "([^"]*)"$`, w.iRun)
	ctx.Step(`^a blog markdown file for today should be created$`, w.aBlogMarkdownFileForTodayShouldBeCreated)
	ctx.Step(`^the blog content should include the idea "([^"]*)"$`, w.theBlogContentShouldIncludeIdea)
	ctx.Step(`^the blog content should include "([^"]*)"$`, w.theBlogContentShouldInclude)
	ctx.Step(`^the blog title should be "([^"]*)"$`, w.theBlogTitleShouldBe)
	ctx.Step(`^the process should exit with a non-zero status$`, w.theProcessShouldExitNonZero)
	ctx.Step(`^the error output should contain "([^"]*)"$`, w.theErrorOutputShouldContain)
}

func TestGodog(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"."},
		},
	}
	if suite.Run() != 0 {
		t.Fail()
	}
}
