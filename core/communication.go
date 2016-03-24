package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
)

// MesosMasterReachable test whether mesos-master is connectable or not
func (core *Core) MesosMasterReachable() bool {
	_, err := http.Get("http://" + core.master + "/health")
	if err != nil {
		core.Log.Errorf("Failed to connect to mesos %v Error: %v\n", core.master, err)
		return false
	}
	return true
}

// SendMessageToMesos is the api to send proto message to mesos-master
func (core *Core) SendMessageToMesos(msg proto.Message, path string) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/master/%s", core.master, path)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-type", "application/octet-stream")
	req.Header.Add("Libprocess-From", fmt.Sprintf("core@%s", core.addr))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code %d received while posting to: %s", resp.StatusCode, url)
	}
	return nil
}

// InitEndpoints init mesos endpoints of core
func (core *Core) InitEndpoints() {
	core.Endpoints = map[string]map[string]func(w http.ResponseWriter, r *http.Request) error{
		"POST": {
			"/core/mesos.internal.FrameworkRegisteredMessage": core.FrameworkRegisteredMessage,
			"/core/mesos.internal.ResourceOffersMessage":      core.ResourceOffersMessage,
			"/core/mesos.internal.StatusUpdateMessage":        core.StatusUpdateMessage,
		},
	}
}

// FrameworkRegisteredMessage is the api that called when framework register to mesos-master
func (core *Core) FrameworkRegisteredMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	message := new(mesosproto.FrameworkRegisteredMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}

	core.Log.WithField("frameworkId", message.FrameworkId).Debug("receive framworkId")
	core.frameworkInfo.Id = message.FrameworkId

	eventType := mesosproto.Event_REGISTERED
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Registered: &mesosproto.Event_Registered{
			FrameworkId: message.FrameworkId,
			MasterInfo:  message.MasterInfo,
		},
	})
	w.WriteHeader(http.StatusOK)
	return nil
}

// ResourceOffersMessage is the api that called when framework receive offers from master
func (core *Core) ResourceOffersMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	message := new(mesosproto.ResourceOffersMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}
	eventType := mesosproto.Event_OFFERS
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Offers: &mesosproto.Event_Offers{
			Offers: message.Offers,
		},
	})
	w.WriteHeader(http.StatusOK)
	return nil
}

// StatusUpdateMessage called when slave's status updated
func (core *Core) StatusUpdateMessage(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	message := new(mesosproto.StatusUpdateMessage)
	if err := proto.Unmarshal(data, message); err != nil {
		return err
	}

	if err := core.SendMessageToMesos(&mesosproto.StatusUpdateAcknowledgementMessage{
		FrameworkId: core.frameworkInfo.Id,
		SlaveId:     message.Update.Status.SlaveId,
		TaskId:      message.Update.Status.TaskId,
		Uuid:        message.Update.Uuid,
	}, "mesos.internal.StatusUpdateAcknowledgementMessage"); err != nil {
		return err
	}

	eventType := mesosproto.Event_UPDATE
	core.AddEvent(eventType, &mesosproto.Event{
		Type: &eventType,
		Update: &mesosproto.Event_Update{
			Uuid:   message.Update.Uuid,
			Status: message.Update.Status,
		},
	})

	w.WriteHeader(http.StatusOK)
	return nil
}
