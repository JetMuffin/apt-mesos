package core

import (
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
)

// Core of sher
type Core interface {
	Run() error
	GetAddr() string
	GetListenIPAndPort() (string, string, error)
	GetAgentLisenPort() string

	HandleFrameworkRegisteredMessage(message *mesosproto.FrameworkRegisteredMessage)
	HandleResourceOffersMessage(message *mesosproto.ResourceOffersMessage)
	HandleStatusUpdateMessage(message *mesosproto.StatusUpdateMessage) error

	RequestOffers() ([]*mesosproto.Offer, error)
	LaunchTask(task *registry.Task, offer *mesosproto.Offer, offers []*mesosproto.Offer) error

	// Task manage
	AddTask(id string, task *registry.Task) error
	GetAllTasks() []*registry.Task
	GetTask(id string) (*registry.Task, error)
	DeleteTask(id string) error
	KillTask(id string) error
	GetUnScheduledTask() []*registry.Task

	// Node manage
	RegisterNode(id string, node *registry.Node) error
	GetNode(id string) (*registry.Node, error)
	UpdateNode(id string, node *registry.Node) error
	DeleteNode(id string) error
	ExistsNode(id string) bool
	GetAllNodes() []*registry.Node

	ReadFile(id string, filename string) (string, error)
	GetSystemUsage() *registry.Metrics
	MergePorts(ports []*registry.Port) *mesosproto.Resource
}
