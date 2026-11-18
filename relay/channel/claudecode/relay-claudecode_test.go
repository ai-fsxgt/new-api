package claudecode

import (
	"encoding/json"
	"one-api/dto"
	"testing"
)

func newSamplingBody() map[string]json.RawMessage {
	return map[string]json.RawMessage{
		"temperature": json.RawMessage(`0.2`),
		"top_p":       json.RawMessage(`0.9`),
		"top_k":       json.RawMessage(`40`),
		"thinking":    json.RawMessage(`{"type":"enabled","budget_tokens":2048}`),
	}
}

func TestRemoveUnsupportedModelParamsForFixedSamplingModels(t *testing.T) {
	for _, model := range []string{"claude-haiku-5-5", "claude-sonnet-5", "claude-opus-4-7", "claude-opus-5-5", "claude-fable-5-2"} {
		t.Run(model, func(t *testing.T) {
			body := newSamplingBody()

			RemoveUnsupportedModelParams(body, model)

			for _, key := range []string{"temperature", "top_p", "top_k"} {
				if _, ok := body[key]; ok {
					t.Fatalf("expected %s to be removed", key)
				}
			}
			if string(body["thinking"]) != `{"type":"adaptive"}` {
				t.Fatalf("expected adaptive thinking, got %s", body["thinking"])
			}
		})
	}
}

func TestRemoveUnsupportedModelParamsKeepsOtherThinkingFields(t *testing.T) {
	body := map[string]json.RawMessage{"thinking": json.RawMessage(`{"type":"enabled","budget_tokens":2048,"display":"summarized"}`)}

	RemoveUnsupportedModelParams(body, "claude-opus-4-7")

	if string(body["thinking"]) != `{"display":"summarized","type":"adaptive"}` {
		t.Fatalf("expected display to be kept, got %s", body["thinking"])
	}
}

func TestRemoveUnsupportedModelParamsKeepsOtherThinkingTypes(t *testing.T) {
	body := map[string]json.RawMessage{"thinking": json.RawMessage(`{"type":"disabled"}`)}

	RemoveUnsupportedModelParams(body, "claude-haiku-5-5")

	if string(body["thinking"]) != `{"type":"disabled"}` {
		t.Fatalf("expected thinking to be kept, got %s", body["thinking"])
	}
}

func TestRemoveUnsupportedModelParamsKeepsOlderModels(t *testing.T) {
	for _, model := range []string{"claude-sonnet-4-6", "claude-opus-4-6", "claude-haiku-4-5-20251001"} {
		t.Run(model, func(t *testing.T) {
			body := newSamplingBody()

			RemoveUnsupportedModelParams(body, model)

			if len(body) != 4 || string(body["thinking"]) != `{"type":"enabled","budget_tokens":2048}` {
				t.Fatalf("expected body to be unchanged, got %v", body)
			}
		})
	}
}

func TestRequestOpenAI2ClaudeMessageDropsSamplingForFixedSamplingModels(t *testing.T) {
	temperature := 0.2
	request := dto.GeneralOpenAIRequest{
		Model:       "claude-haiku-5-5",
		Temperature: &temperature,
		TopP:        0.9,
		TopK:        40,
		Messages:    []dto.Message{{Role: "user", Content: json.RawMessage(`"hi"`)}},
	}

	claudeRequest, err := RequestOpenAI2ClaudeMessage(request, nil)
	if err != nil {
		t.Fatal(err)
	}

	if claudeRequest.Temperature != nil || claudeRequest.TopP != 0 || claudeRequest.TopK != 0 {
		t.Fatalf("expected sampling params to be removed, got %v %v %v", claudeRequest.Temperature, claudeRequest.TopP, claudeRequest.TopK)
	}
}
