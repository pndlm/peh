package peh3_test

import (
	"bytes"
	"testing"

	"github.com/pndlm/peh/peh3"
	"gopkg.in/yaml.v3"
)

const TestDropComposeServiceA = `version: '3.8'
services:
  svc1:
    x: y
  svc2:
    x: y
`

const TestDropComposeServiceB = `version: '3.8'
services:
  svc2:
    x: y
`

func TestComposeDropService(t *testing.T) {
	var root yaml.Node
	err := yaml.Unmarshal([]byte(TestDropComposeServiceA), &root)
	if err != nil {
		panic(err)
	}
	// sanity checks
	if root.Content[0].Content[1].Value != "3.8" {
		t.Fatalf("got %s", root.Content[0].Content[1].Value)
	}
	if root.Content[0].Content[2].Value != "services" {
		t.Fatalf("got %s", root.Content[0].Content[2].Value)
	}
	if root.Content[0].Content[3].Content[0].Value != "svc1" {
		t.Fatalf("got %s", root.Content[0].Content[3].Content[0].Value)
	}
	compose := &peh3.Compose{
		Node: &root,
	}
	compose.DropService("svc1")
	if root.Content[0].Content[3].Content[0].Value != "svc2" {
		t.Fatalf("got %s", root.Content[0].Content[3].Content[0].Value)
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2) // match what has already been foisted upon us
	err = enc.Encode(&root)
	if err != nil {
		panic(err)
	}
	s := buf.String()
	if s != TestDropComposeServiceB {
		t.Fatalf("got %s", s)
	}
}
