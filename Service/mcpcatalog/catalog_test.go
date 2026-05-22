package mcpcatalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSupportedOpsJSONEnvelope(t *testing.T) {
	raw := SupportedOpsJSON()
	var env struct {
		Version      int                  `json:"version"`
		Ops          []string             `json:"ops"`
		Capabilities []BridgeOpCapability `json:"capabilities"`
	}
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		t.Fatal(err)
	}
	if env.Version != 1 {
		t.Fatalf("version: %d", env.Version)
	}
	if len(env.Ops) == 0 || len(env.Capabilities) == 0 {
		t.Fatalf("empty ops=%d caps=%d", len(env.Ops), len(env.Capabilities))
	}
	if len(env.Ops) != len(env.Capabilities) {
		t.Fatalf("ops len %d vs caps len %d", len(env.Ops), len(env.Capabilities))
	}
	for i := range env.Ops {
		if env.Ops[i] != env.Capabilities[i].Op {
			t.Fatalf("mismatch at %d: %q vs %q", i, env.Ops[i], env.Capabilities[i].Op)
		}
		if env.Capabilities[i].Description == "" {
			t.Fatalf("empty description for %q", env.Capabilities[i].Op)
		}
		if strings.TrimSpace(env.Capabilities[i].Returns) == "" {
			t.Fatalf("empty returns for %q", env.Capabilities[i].Op)
		}
	}
}

func TestMainRowNoteSetInputSchemaHasRowIDs(t *testing.T) {
	schema := ToolInputSchema("main_row_note_set")
	if schema == nil {
		t.Fatal("expected schema for main_row_note_set")
	}
	b, err := json.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	raw := string(b)
	if !strings.Contains(raw, "rowIds") {
		t.Fatalf("schema missing rowIds: %s", raw)
	}
	if !strings.Contains(raw, "note") {
		t.Fatalf("schema missing note: %s", raw)
	}
}

func TestBridgeMCPToolsMainRowNoteSetDescription(t *testing.T) {
	tools := BridgeMCPTools()
	var desc string
	for _, tt := range tools {
		if tt.Op == "main_row_note_set" {
			desc = tt.Description
			break
		}
	}
	if desc == "" {
		t.Fatal("main_row_note_set tool not found")
	}
	if !strings.Contains(desc, "禁止") || !strings.Contains(desc, "rowIds") {
		t.Fatalf("batch guidance missing: %s", desc)
	}
}

func TestBridgeMCPTools(t *testing.T) {
	tools := BridgeMCPTools()
	if len(tools) != len(SupportedBridgeOps) {
		t.Fatalf("tools %d vs ops %d", len(tools), len(SupportedBridgeOps))
	}
	for i, tt := range tools {
		want := BridgeMCPHTTPPrefix + SupportedBridgeOps[i]
		if tt.MCPName != want {
			t.Fatalf("idx %d: mcp name %q want %q", i, tt.MCPName, want)
		}
		if tt.Op != SupportedBridgeOps[i] {
			t.Fatalf("idx %d: op %q want %q", i, tt.Op, SupportedBridgeOps[i])
		}
		if tt.Description == "" {
			t.Fatalf("empty description for %q", tt.MCPName)
		}
	}
}
