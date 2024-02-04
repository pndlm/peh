package peh3

import (
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

// The reason to use all this custom stuff with yaml.Node now is because
// the structs in the official compose-go project
// https://github.com/compose-spec/compose-go
// are oddly errant...  weird given that Docker itself is written in go

type Compose struct {
	Node *yaml.Node
}

func (compose *Compose) ServicesNode() *yaml.Node {
	// https://github.com/go-yaml/yaml/blob/v3/node_test.go
	for _, m1 := range compose.Node.Content {
		for i2, m2 := range m1.Content {
			if m2.Value == "services" {
				return m1.Content[i2+1]
			}
		}
	}
	panic(fmt.Errorf("could not identify services node"))
}

func (compose *Compose) DropService(serviceName string) {
	services := compose.ServicesNode()
	for i3, m3 := range services.Content {
		if m3.Value == serviceName {
			services.Content = slices.Delete(services.Content, i3, i3+2)
			return
		}
	}
	panic(fmt.Errorf("could not identify service %s", serviceName))
}

func ComposeFromFile(path string) *Compose {
	node := YamlNodeFromFile(path)
	if node.Kind != yaml.DocumentNode {
		panic(fmt.Errorf("expected yaml.DocumentNode"))
	}
	return &Compose{
		Node: node,
	}
}

// Parse yaml file to get *yaml.Node
func YamlNodeFromFile(path string) (doc *yaml.Node) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var node yaml.Node
	err = yaml.Unmarshal(b, &node)
	if err != nil {
		panic(err)
	}
	return &node
}

// Write yaml file from *yaml.Node
func YamlNodeToFile(doc *yaml.Node, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2) // match what has already been foisted upon us
	err = enc.Encode(doc)
	if err != nil {
		panic(err)
	}
}
