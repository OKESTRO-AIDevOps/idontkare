package apiximpl

import (
	"fmt"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
)

func V1UserSet(name string, pass string) error {

	err := pkgdbquery.SetUser(name, pass)

	if err != nil {

		return fmt.Errorf("failed to set user: %s", err.Error())
	}

	return nil

}
