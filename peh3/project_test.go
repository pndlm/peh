package peh3_test

import (
	"testing"

	"github.com/pndlm/peh/peh3"
)

func TestRel(t *testing.T) {
	proj := &peh3.Project{
		Dir:  "/Users/pndlm/example-project",
		Name: "example-project",
	}
	// it even eliminates duplicate slashes!
	relPath := proj.RelPath("//docker/working/letsencrypt")
	expected := "/Users/pndlm/example-project/docker/working/letsencrypt"
	if relPath != expected {
		t.Fatalf("expected %s got %s", expected, relPath)
	}
}
