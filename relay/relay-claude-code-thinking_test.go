package relay

import (
	"encoding/json"
	"testing"
)

func TestNormalizeInvalidThinkingKeepsClaudeSignatures(t *testing.T) {
	raw := json.RawMessage(`[
		{"role":"assistant","content":[
			{"type":"thinking","thinking":"plan a","signature":"EpIECkYIChgCKkBhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5eg=="},
			{"type":"tool_use","id":"toolu_1","name":"Read","input":{}}
		]},
		{"role":"assistant","content":[
			{"type":"thinking","thinking":"plan b","signature":"CAQSYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="},
			{"type":"text","text":"done"}
		]}
	]`)

	patched, changed := normalizeInvalidThinkingInMessagesRaw(raw)

	if changed {
		t.Fatalf("expected Claude thinking blocks to be kept, got %s", patched)
	}
}

func TestNormalizeInvalidThinkingConvertsNonClaudeSignatures(t *testing.T) {
	for name, signature := range map[string]string{
		"empty":       "",
		"uuid":        "5e3c1a2b-7d4f-4a6e-9b8c-0d1e2f3a4b5c",
		"placeholder": "skip_thought_signature_validator",
	} {
		t.Run(name, func(t *testing.T) {
			raw := json.RawMessage(`[{"role":"assistant","content":[
				{"type":"thinking","thinking":"reasoning","signature":"` + signature + `"},
				{"type":"text","text":"answer"}
			]}]`)

			patched, changed := normalizeInvalidThinkingInMessagesRaw(raw)

			if !changed {
				t.Fatalf("expected thinking block to be converted")
			}
			got := string(patched)
			if jsonContains(got, `"type":"thinking"`) || !jsonContains(got, `"text":"reasoning"`) {
				t.Fatalf("unexpected messages: %s", got)
			}
		})
	}
}
