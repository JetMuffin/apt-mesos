package server

import (
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/bitly/go-simplejson"
	
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/server"
	. "github.com/icsnju/apt-mesos/registry"
)

func init() {
	r := registry.NewRegistry()
	go server.ListenAndServe(":3000", r)	
}

func TestHandshake(t *testing.T) {
	res, err := http.Get("http://localhost:3000/api/handshake")
	defer res.Body.Close()
	assert.NoError(t, err)

    body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	object, err := simplejson.NewJson(body)
	assert.NoError(t, err)
	assert.Equal(t, true, object.Get("success").MustBool())
}

func TestSubmitTask(t *testing.T) {
	task := &Task{
		Command: 		"touch /data/volt",
    	Cpus: 			0.1,
    	Mem: 			32,
    	DockerImage: 	"busybox",
	}
	b, err := json.Marshal(task)
	assert.NoError(t, err)

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://localhost:3000/api/tasks", "application/json;charset=utf-8", body)
	assert.NoError(t, err)

	result, err := ioutil.ReadAll(res.Body)
  	defer res.Body.Close()
  	assert.NoError(t, err)

	object, err := simplejson.NewJson(result)
	assert.NoError(t, err)
	assert.Equal(t, true, object.Get("success").MustBool())
}

func TestListTasks(t *testing.T) {
	res, err := http.Get("http://localhost:3000/api/tasks")
	assert.NoError(t, err)

    body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	object, err := simplejson.NewJson(body)
	assert.NoError(t, err)
	assert.Equal(t, true, object.Get("success").MustBool())	
	assert.Equal(t, 1, len(object.Get("result").MustArray()))
}

// TODO test remove task, update task