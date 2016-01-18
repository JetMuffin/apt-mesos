package core

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/icsnju/apt-mesos/mesosproto"
    "github.com/icsnju/apt-mesos/core"
    "github.com/icsnju/apt-mesos/server"
    "github.com/icsnju/apt-mesos/registry"
)

var (
	c *core.Core
)

func init() {
	frameworkName := "api-mesos test" 
	user := "tester" 
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	log	:= logrus.New()
	log.Level = logrus.DebugLevel
	r := registry.NewRegistry()
	c = core.NewCore("192.168.33.1:3000", "192.168.33.10:5050", frameworkInfo, log)
	server.ListenAndServe("192.168.33.1:3000", r, c)	
}

func TestRegisterFramework(t *testing.T) {
	err := c.RegisterFramework()
	assert.NoError(t, err)
	var event *mesosproto.Event
	select{
		case event = <-c.GetEvent(mesosproto.Event_REGISTERED):
	}
	assert.NotNil(t, event)
}

func TestRequestOffers(t *testing.T) {
	resources := c.BuildResources(1, 16, 10)
	_, err := c.RequestOffers(resources)
	assert.NoError(t, err)
}

