package core

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/mesos/mesos-go/mesosproto"
    "github.com/icsnju/apt-mesos/core"
)

var (
	c *core.Core
)

func init() {
	frameworkName := "api-mesos test" 
	user := "tester" 
	frameworkInfo := &mesosproto.FrameworkInfo{Name: &frameworkName, User: &user}
	log	:= logrus.New()
	c = core.NewCore("127.0.0.1:3000", "192.168.33.10:5050", frameworkInfo, log)	
}

func TestRegisterFramework(t *testing.T) {
	err := c.RegisterFramework()
	assert.NoError(t, err)
}

