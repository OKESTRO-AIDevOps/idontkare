package apix

import (
	"fmt"
	"os"

	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	"gopkg.in/yaml.v3"
)

func V1GetManifest() (*pkgresourceapix.V1Manifest, error) {

	var v1man pkgresourceapix.V1Manifest

	file_b, err := os.ReadFile("apix/v1/manifest.yaml")

	if err != nil {

		return nil, fmt.Errorf("failed to get manifest: %s", err.Error())
	}

	err = yaml.Unmarshal(file_b, &v1man)

	if err != nil {

		return nil, fmt.Errorf("failed to get manifest: unmarshal: %s", err.Error())
	}

	return &v1man, nil

}

func V1GetMain(main_b []byte, mani *pkgresourceapix.V1Manifest) (*pkgresourceapix.V1Main, error) {

	var v1main pkgresourceapix.V1Main

	err := yaml.Unmarshal(main_b, &v1main)

	if err != nil {

		return nil, fmt.Errorf("invalid format: %s", err.Error())

	}

	realMain := mani.Main

	rmlen := len(realMain)

	found := -1

	targetPath := v1main.Kind + v1main.Path

	for i := 0; i < rmlen; i++ {

		thisPath := realMain[i].Kind + realMain[i].Path

		if targetPath == thisPath {

			found = i

			break
		}

	}

	if found == -1 {

		return nil, fmt.Errorf("invalid path: %s", targetPath)

	}

	realBody := realMain[found].Body

	rblen := len(realBody)

	blen := len(v1main.Body)

	if rblen != blen {

		return nil, fmt.Errorf("invalid body len: %d", blen)
	}

	for k, _ := range v1main.Body {

		_, okay := realBody[k]

		if !okay {
			return nil, fmt.Errorf("invalid key: %s", k)
		}

		rblen -= 1

	}

	if rblen != 0 {

		return nil, fmt.Errorf("invalid key count: left: %d", rblen)
	}

	return &v1main, nil
}
