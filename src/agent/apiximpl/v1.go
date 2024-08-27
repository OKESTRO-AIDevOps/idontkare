package apiximpl

import (
	"fmt"
	"time"

	pkgcomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/comm"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourcecomm "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/comm"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

var V1TIMEOUT_MS int = 5000
var interval_ms int = 10

func V1ServerWriteCtl(v1main *pkgresourceapix.V1Main) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	route := v1main.Path

	switch route {

	}

	return &resp, nil

}

func V1ServerReadCtl(v1main *pkgresourceapix.V1Main) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	route := v1main.Path

	switch route {

	}

	return &resp, nil

}

func V1AgentRequestCtl(v1main *pkgresourceapix.V1Main, c *websocket.Conn, sess_key []byte) (*pkgresourceapix.V1ResultData, error) {

	route := v1main.Path

	switch route {

	}

	resp, err := V1RoundTrip(v1main, c, sess_key)

	if err != nil {

		return resp, err
	}

	return resp, nil

}

func V1RoundTrip(v1main *pkgresourceapix.V1Main, c *websocket.Conn, sess_key []byte) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	data, err := yaml.Marshal(v1main)

	if err != nil {

		return nil, fmt.Errorf("roundtrip: marshal: %s", err.Error())
	}

	if sess_key != nil {

		data_str, err := pkgcomm.CommDataEncrypt(data, sess_key)

		if err != nil {

			return nil, fmt.Errorf("roundtrip encrypt: %s", err.Error())
		}

		data = []byte(data_str)

	}
	commjson := pkgresourcecomm.CommJSON{
		Status: pkgresourcecomm.COMM_STATUS_REQUEST,
		Data:   data,
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(interval_ms))

	timelimit_count := 0
	timelimit := V1TIMEOUT_MS / interval_ms

	read_wait_loop := 1

	error_channel := make(chan error)
	read_channel := make(chan *pkgresourcecomm.CommJSON)

	go func() {

		var read_commjson_in pkgresourcecomm.CommJSON

		err := c.ReadJSON(&read_commjson_in)

		if err != nil {

			error_channel <- err
			return
		}

		read_channel <- &read_commjson_in

	}()

	err = c.WriteJSON(commjson)

	if err != nil {
		return nil, fmt.Errorf("routdtrip: write: %s", err.Error())
	}

	var read_commjson *pkgresourcecomm.CommJSON

	for read_wait_loop == 1 {

		select {

		case <-ticker.C:

			timelimit_count += 1

			if timelimit_count >= timelimit {

				ticker.Stop()

				return nil, fmt.Errorf("roundtrip: timeout")

			}

		case read_error := <-error_channel:

			return nil, fmt.Errorf("roundtrip: read: %s", read_error.Error())

		case read_commjson = <-read_channel:

			read_wait_loop = 0

		}
	}

	if read_commjson.Status != pkgresourcecomm.COMM_STATUS_SUCCESS {

		return nil, fmt.Errorf("roundtrip: comm failed: %s", read_commjson.Message)

	}

	resp_data := read_commjson.Data

	if sess_key != nil {

		resp_data, err = pkgcomm.CommDataDecrypt(string(resp_data), sess_key)

		if err != nil {
			return nil, fmt.Errorf("roundtrip: decrypt failed: %s", err.Error())
		}

	}

	err = yaml.Unmarshal(resp_data, &resp)

	if err != nil {

		return nil, fmt.Errorf("roundtrip: data unmarshal: %s", err.Error())
	}

	return &resp, nil
}
