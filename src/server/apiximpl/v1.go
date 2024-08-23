package apiximpl

import (
	"fmt"
	"log"

	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	"github.com/gorilla/websocket"
)

func V1ClientRequestCtl(v1main *pkgresourceapix.V1Main, c *websocket.Conn) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	route := v1main.Path

	switch route {

	case "/user/set":

		name := v1main.Body["name"]
		pass := v1main.Body["pass"]

		err := V1UserSet(name, pass)

		if err != nil {

			log.Printf("failed to set user: %s: %s", name, err.Error())

			resp.Output = fmt.Sprintf("failed to set user: %s", name)

		} else {

			resp.Output = fmt.Sprintf("successfully set user: %s", name)

		}

	case "/cluster/set":

		_ = v1main.Body["username"]
		_ = v1main.Body["name"]

		resp.Output = "not yet implemented: " + route

	case "/project/set":

		_ = v1main.Body["username"]
		_ = v1main.Body["name"]
		_ = v1main.Body["git"]
		_ = v1main.Body["gitid"]
		_ = v1main.Body["gitpw"]
		_ = v1main.Body["reg"]
		_ = v1main.Body["regid"]
		_ = v1main.Body["regpw"]

		resp.Output = "not yet implemented: " + route

	case "/project/ci/option/set":

		_ = v1main.Body["username"]
		_ = v1main.Body["name"]
		_ = v1main.Body["path"]

		resp.Output = "not yet implemented: " + route

	case "/project/cd/option/set":

		_ = v1main.Body["username"]
		_ = v1main.Body["name"]
		_ = v1main.Body["path"]

		resp.Output = "not yet implemented: " + route

	default:

		resp.Output = "no such path: " + route

	}

	return &resp, nil

}

func V1AgentRequestCtl(v1main *pkgresourceapix.V1Main, c *websocket.Conn) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	return &resp, nil

}

func V1AgentRoundTrip(v1main *pkgresourceapix.V1Main, c *websocket.Conn) (*pkgresourceapix.V1ResultData, error) {

	var resp pkgresourceapix.V1ResultData

	return &resp, nil
}
