package mcpbridge

import (
	"testing"
)

func TestArgRowIDs(t *testing.T) {
	t.Parallel()
	ids, err := argRowIDs(map[string]any{
		"rowId":  "1",
		"rowIds": []any{"2", "1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != "1" || ids[1] != "2" {
		t.Fatalf("got %v", ids)
	}
	_, err = argRowIDs(map[string]any{})
	if err == nil {
		t.Fatal("expected error for empty")
	}
}

func TestArgTheologyList(t *testing.T) {
	t.Parallel()
	th, err := argTheologyList(map[string]any{
		"theology": float64(12),
		"rowIds":   []any{"3"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(th) < 2 {
		t.Fatalf("got %v", th)
	}
}
