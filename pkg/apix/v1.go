package apix

import (
	"fmt"
	"os"
	"strings"

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

func V1GetMainByByte(main_b []byte, mani *pkgresourceapix.V1Manifest) (*pkgresourceapix.V1Main, error) {

	var v1main pkgresourceapix.V1Main

	err := yaml.Unmarshal(main_b, &v1main)

	if err != nil {

		return nil, fmt.Errorf("invalid format: %s", err.Error())

	}

	rmlen := len(mani.Main)

	found := -1

	targetPath := v1main.Kind + v1main.Path

	for i := 0; i < rmlen; i++ {

		thisPath := mani.Main[i].Kind + mani.Main[i].Path

		if targetPath == thisPath {

			found = i

			break
		}

	}

	if found == -1 {

		return nil, fmt.Errorf("invalid path: %s", targetPath)

	}

	realHead := mani.Main[found].Head

	if realHead.FromFile != nil {

		fromFilePath, v1mainokay := v1main.Body[pkgresourceapix.V1HeadFromFile]

		if v1mainokay {

			file_b, err := os.ReadFile(fromFilePath)

			if err != nil {

				return nil, fmt.Errorf("from-file flag: file doens't exist: %s", err.Error())
			}

			var fileBody pkgresourceapix.V1Body

			err = yaml.Unmarshal(file_b, &fileBody)

			if err != nil {

				return nil, fmt.Errorf("from-file flag: invalid format: %s", err.Error())

			}

			v1main.Body = fileBody

		}

	}

	realBody := mani.Main[found].Body

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

func V1GetMainCopyByAddress(kind string, path string, mani *pkgresourceapix.V1Manifest) (*pkgresourceapix.V1Main, error) {

	var retMain pkgresourceapix.V1Main

	realMain := mani.Main

	rmlen := len(realMain)

	found := -1

	targetPath := kind + path

	for i := 0; i < rmlen; i++ {

		thisPath := realMain[i].Kind + realMain[i].Path

		if targetPath == thisPath {

			found = i

			break
		}

	}

	if found == -1 {

		return nil, fmt.Errorf("not found: %s", targetPath)

	}

	retMain = realMain[found]

	return &retMain, nil

}

func V1GetMainFromArgs(kind string, args []string, mani *pkgresourceapix.V1Manifest) (*pkgresourceapix.V1Main, error) {

	pathString := ""

	possibleBody := make(pkgresourceapix.V1Body)

	idx := 0

	allen := len(args)

	for {

		if idx >= allen {

			break
		}

		if strings.HasPrefix(args[idx], "--") {

			arg := strings.ReplaceAll(args[idx], "--", "")

			_, okay := possibleBody[arg]

			if okay {
				return nil, fmt.Errorf("duplicate key: %s", arg)
			}

			if idx+1 >= allen {
				return nil, fmt.Errorf("value not provided for: %s", arg)
			}

			idx += 1

			if strings.HasPrefix(args[idx], "--") {
				return nil, fmt.Errorf("value not provided for: %s", arg)
			}

			val := args[idx]

			possibleBody[arg] = val

			idx += 1

		} else {

			arg := args[idx]

			pathString += "/" + arg

			idx += 1

		}

	}

	realMain, err := V1GetMainCopyByAddress(kind, pathString, mani)

	if err != nil {
		return nil, fmt.Errorf("failed to get main template for: %s", kind+pathString)
	}

	if realMain == nil {

		return nil, fmt.Errorf("main template not found for: %s", kind+pathString)
	}

	realMain.Kind = kind
	realMain.Path = pathString
	realMain.Body = possibleBody

	rm_b, err := yaml.Marshal(realMain)

	if err != nil {

		return nil, fmt.Errorf("failed to marshal into byte: %s", err.Error())
	}

	finalMain, err := V1GetMainByByte(rm_b, mani)

	if err != nil {

		return nil, fmt.Errorf("failed to finalize main: %s", err.Error())

	}

	return finalMain, nil

}
