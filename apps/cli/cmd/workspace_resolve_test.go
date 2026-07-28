package cmd

import (
	"testing"

	"github.com/zain-23/local-vault/apps/cli/internal/api"
)

func memberships(ids ...string) []api.WorkspaceMembership {
	out := make([]api.WorkspaceMembership, len(ids))
	for i, id := range ids {
		out[i] = api.WorkspaceMembership{Workspace: api.Workspace{ID: id, Name: id}}
	}
	return out
}

func TestResolveWorkspaceFlag(t *testing.T) {
	id, err := resolveWorkspaceID("ws_2", memberships("ws_1", "ws_2"), nil)
	if err != nil || id != "ws_2" {
		t.Fatalf("got %q %v", id, err)
	}
}

func TestResolveWorkspaceFlagNotMember(t *testing.T) {
	_, err := resolveWorkspaceID("ws_x", memberships("ws_1"), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveWorkspaceAutoOne(t *testing.T) {
	id, err := resolveWorkspaceID("", memberships("ws_1"), nil)
	if err != nil || id != "ws_1" {
		t.Fatalf("got %q %v", id, err)
	}
}

func TestResolveWorkspaceZero(t *testing.T) {
	_, err := resolveWorkspaceID("", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveWorkspacePick(t *testing.T) {
	id, err := resolveWorkspaceID("", memberships("ws_1", "ws_2"), func() (string, error) { return "2", nil })
	if err != nil || id != "ws_2" {
		t.Fatalf("got %q %v", id, err)
	}
}

func TestResolveWorkspacePickBad(t *testing.T) {
	_, err := resolveWorkspaceID("", memberships("ws_1", "ws_2"), func() (string, error) { return "9", nil })
	if err == nil {
		t.Fatal("expected error")
	}
}
