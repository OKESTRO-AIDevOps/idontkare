package main

import (
	"fmt"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	"gopkg.in/yaml.v3"
)

func apix_get_then_verify() {

	mani, err := pkgapix.V1GetManifest()

	if err != nil {

		fmt.Println(err.Error())

		return

	}

	body := make(pkgresourceapix.V1Body)

	body["name"] = "hello"
	body["pass"] = "world"

	main_no_head := pkgresourceapix.V1Main{
		Kind: pkgresourceapix.V1KindClientRequest,
		Path: "/user/set",
		Body: body,
	}

	yb, _ := yaml.Marshal(main_no_head)

	_, err = pkgapix.V1GetMainByByte(yb, mani)

	if err != nil {

		fmt.Printf("failed: no head: %s\n", err.Error())

		return

	} else {

		fmt.Printf("success: no head\n")
	}

	body = make(pkgresourceapix.V1Body)

	body["name"] = "wrong"

	main_body_invalid1 := pkgresourceapix.V1Main{
		Kind: pkgresourceapix.V1KindClientRequest,
		Path: "/user/set",
		Body: body,
	}

	yb, _ = yaml.Marshal(main_body_invalid1)

	_, err = pkgapix.V1GetMainByByte(yb, mani)

	if err != nil {

		fmt.Printf("success: body invalid1: %s\n", err.Error())

	} else {

		fmt.Printf("failed: body invalid1\n")

		return
	}

	body = make(pkgresourceapix.V1Body)

	body["name"] = "wrong"
	body["bogus"] = "ney"

	main_body_invalid2 := pkgresourceapix.V1Main{
		Kind: pkgresourceapix.V1KindClientRequest,
		Path: "/user/set",
		Body: body,
	}

	yb, _ = yaml.Marshal(main_body_invalid2)

	_, err = pkgapix.V1GetMainByByte(yb, mani)

	if err != nil {

		fmt.Printf("success: body invalid2: %s\n", err.Error())

	} else {

		fmt.Printf("failed: body invalid2\n")

		return
	}

	body = make(pkgresourceapix.V1Body)

	body["name"] = "hello"
	body["pass"] = "world"

	main_path_invalid := pkgresourceapix.V1Main{
		Kind: pkgresourceapix.V1KindClientRequest,
		Path: "/user/settttt",
		Body: body,
	}

	yb, _ = yaml.Marshal(main_path_invalid)

	_, err = pkgapix.V1GetMainByByte(yb, mani)

	if err != nil {

		fmt.Printf("success: path invalid: %s\n", err.Error())

	} else {

		fmt.Printf("failed: path invalid\n")

		return
	}

}

func apix_get_main_from_args() {

	mani, err := pkgapix.V1GetManifest()

	if err != nil {

		fmt.Println(err.Error())

		return

	}

	kind := pkgresourceapix.V1KindAgentRequest

	args := []string{
		"cluster",
		"connect",
		"--name",
		"hello",
		"--username",
		"world",
	}

	_, err = pkgapix.V1GetMainFromArgs(kind, args, mani)

	if err != nil {
		fmt.Printf("failed: get main from args: %s\n", err.Error())

		return
	} else {

		fmt.Printf("success: get main from args\n")
	}

	kind = pkgresourceapix.V1KindClientRequest

	args = []string{
		"cluster",
		"ci",
		"ready",
		"--name",
		"wrong",
	}

	_, err = pkgapix.V1GetMainFromArgs(kind, args, mani)

	if err != nil {

		fmt.Printf("success: wrong kind and path: %s\n", err.Error())

	} else {

		fmt.Printf("failed: wrong kind and path\n")

		return
	}

	kind = pkgresourceapix.V1KindClientRequest

	args = []string{

		"user",
		"set",
		"--name",
		"--pass",
		"lllll",
	}

	_, err = pkgapix.V1GetMainFromArgs(kind, args, mani)

	if err != nil {
		fmt.Printf("success: imperfect arguments: %s\n", err.Error())

	} else {

		fmt.Printf("failed: imperfect arguments\n")

		return
	}

}

func main() {

	apix_get_then_verify()

	apix_get_main_from_args()

}
