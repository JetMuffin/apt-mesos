package impl

import (
	"errors"
	"time"

	"github.com/icsnju/apt-mesos/registry"
)

// ErrNodeNotExists defined errors
var (
	ErrNodeNotExists = errors.New("Specific task not exist")
)

// RegisterNode register node to registry
func (core *Core) RegisterNode(id string, node *registry.Node) error {
	return core.nodes.Add(id, node)
}

// GetNode return node with id
func (core *Core) GetNode(id string) (*registry.Node, error) {
	if node := core.nodes.Get(id); node != nil {
		return node.(*registry.Node), nil
	}
	return nil, ErrNodeNotExists
}

// UpdateNode update node information
func (core *Core) UpdateNode(id string, node *registry.Node) error {
	node.LastUpdateTime = time.Now().Unix()
	if err := core.nodes.Update(id, node); err != nil {
		return err
	}
	return nil
}

// DeleteNode delete node with given id
func (core *Core) DeleteNode(id string) error {
	return core.nodes.Delete(id)
}

// ExistsNode return if node exists in registry
func (core *Core) ExistsNode(id string) bool {
	return core.nodes.Exists(id)
}

// GetAllNodes return all nodes
func (core *Core) GetAllNodes() []*registry.Node {
	rawList := core.nodes.List()
	nodes := make([]*registry.Node, len(rawList))

	for i, v := range rawList {
		nodes[i] = v.(*registry.Node)
	}
	return nodes
}
