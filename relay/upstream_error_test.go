package relay

import (
	"net/http"
	"testing"

	"one-api/dto"
)

func newContextLimitUpstreamError(message string) *dto.OpenAIErrorWithStatusCode {
	return &dto.OpenAIErrorWithStatusCode{
		Error:      dto.OpenAIError{Message: message, Type: "upstream_error", Code: "bad_response_status_code"},
		StatusCode: http.StatusBadGateway,
	}
}

func TestNormalizeClaudeContextLimitErrorRewritesMaxTokensOverflow(t *testing.T) {
	openaiErr := newContextLimitUpstreamError("This model configuration accepts at most 262140 combined input and output tokens. However, your request has 235926 input tokens and asks for 32000 output tokens (267926 tokens total). Please reduce the input length or requested output length.")

	NormalizeClaudeContextLimitError(openaiErr)

	want := "input length and `max_tokens` exceed context limit: 235926 + 32000 > 262140, decrease input length or `max_tokens` and try again"
	if openaiErr.Error.Message != want {
		t.Fatalf("unexpected message: %s", openaiErr.Error.Message)
	}
	if openaiErr.Error.Type != "invalid_request_error" || openaiErr.StatusCode != http.StatusBadRequest || !openaiErr.LocalError {
		t.Fatalf("unexpected type, status or local flag: %s %d %v", openaiErr.Error.Type, openaiErr.StatusCode, openaiErr.LocalError)
	}
}

func TestNormalizeClaudeContextLimitErrorRewritesPromptTooLong(t *testing.T) {
	openaiErr := newContextLimitUpstreamError("This model configuration accepts at most 262140 combined input and output tokens. However, your request has 270000 input tokens and asks for 32000 output tokens (302000 tokens total). Please reduce the input length or requested output length.")

	NormalizeClaudeContextLimitError(openaiErr)

	want := "prompt is too long: 270000 tokens > 262140 maximum"
	if openaiErr.Error.Message != want {
		t.Fatalf("unexpected message: %s", openaiErr.Error.Message)
	}
	if openaiErr.Error.Type != "invalid_request_error" || openaiErr.StatusCode != http.StatusBadRequest || !openaiErr.LocalError {
		t.Fatalf("unexpected type, status or local flag: %s %d %v", openaiErr.Error.Type, openaiErr.StatusCode, openaiErr.LocalError)
	}
}

func TestNormalizeClaudeContextLimitErrorKeepsOtherErrors(t *testing.T) {
	openaiErr := newContextLimitUpstreamError("upstream overloaded")

	NormalizeClaudeContextLimitError(openaiErr)

	if openaiErr.Error.Message != "upstream overloaded" || openaiErr.Error.Type != "upstream_error" || openaiErr.StatusCode != http.StatusBadGateway || openaiErr.LocalError {
		t.Fatalf("unrelated error was modified: %#v", openaiErr)
	}
}
