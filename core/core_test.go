package core

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/icsnju/apt-mesos/core"
	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/server"
	"github.com/stretchr/testify/assert"
)

var (
	c *core.Core
)

func init() {
	frameworkName := "api-mesos test"
	user := "tester"
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	log := logrus.New()
	log.Level = logrus.DebugLevel
	r := registry.NewTaskRegistry()
	c = NewCore("192.168.33.1:3000", "192.168.33.10:5050", frameworkInfo, log)
	server.ListenAndServe("192.168.33.1:3000", r, c)
}

func TestRegisterFramework(t *testing.T) {
	err := cc.RegisterFramework()
	assert.NoError(t, err)
	var event *mesosproto.Event
	select {
	case event = <-cc.GetEvent(mesosproto.Event_REGISTERED):
	}
	assert.NotNil(t, event)
}

func TestLaunchTask(t *testing.T) {
	commands := "echo hello"
	task := &registry.Task{
		ID:          "1",
		Command:     commands,
		DockerImage: "mini",
		Volumes:     nil,
	}
	resources := cc.BuildResources(1, 16, 10)
	offers, err := cc.RequestOffers(resources)
	assert.NoError(t, err)
	err = cc.LaunchTask(offers[0], resources, task)
	assert.NoError(t, err)
}

func TestKillTask(t *testing.T) {
	err := cc.KillTask("1")
	assert.NoError(t, err)
}
