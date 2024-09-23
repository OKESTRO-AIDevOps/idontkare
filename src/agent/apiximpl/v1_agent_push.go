package apiximpl

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	pkgapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/apix"
	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
	pkgresourcelc "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/lifecycle"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
	"gopkg.in/yaml.v3"
)

func V1CiHandler(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		cilen := len(CI_OPTIONS_Q)

		term_list := []int{}

		for i := 0; i < cilen; i++ {

			if CI_OPTIONS_Q[i].Status == pkgresourceci.STATUS_READY {

				CI_OPTIONS_Q[i].Status = pkgresourceci.STATUS_RUNNING

				go V1CiHandler_BuildStart(i)

			} else if CI_OPTIONS_Q[i].Status == pkgresourceci.STATUS_RUNNING {

				V1CiHandler_BuildReport(i, mani, acon)

			} else {

				V1CiHandler_BuildReport(i, mani, acon)

				term_list = append(term_list, i)

			}

		}

		err := V1CiClear(term_list)

		if err != nil {

			log.Printf("failed to clear ci: %s", err.Error())
		}

	}

}

func V1CiHandler_BuildStart(aidx int) {

	bid, _ := pkgutils.RandomHex(8)

	v_repoaddr := pkgresourceci.BUILD_PROTOCOL + CI_OPTIONS_Q[aidx].Git
	repo_id := CI_OPTIONS_Q[aidx].GitId
	repo_pw := CI_OPTIONS_Q[aidx].GitPw
	v_regaddr := CI_OPTIONS_Q[aidx].Reg
	reg_id := CI_OPTIONS_Q[aidx].RegId
	reg_pw := CI_OPTIONS_Q[aidx].RegPw

	ns_bid := "npia-build-ns-" + bid
	sec_bid := "npia-build-secret-" + bid
	pod_bid := "npia-build-pod-" + bid

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	cmd := exec.Command("kubectl", "create", "namespace", ns_bid)

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	if err != nil {

		CI_OPTIONS_Q[aidx].Log += errBuf.String()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return

	} else {

		CI_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	outBuf = bytes.Buffer{}
	errBuf = bytes.Buffer{}

	docker_server := "--docker-server=" + v_regaddr
	docker_uname := "--docker-username=" + reg_id
	docker_pword := "--docker-password=" + reg_pw

	cmd = exec.Command("kubectl", "-n", ns_bid, "create", "secret", "docker-registry", sec_bid, docker_server, docker_uname, docker_pword)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CI_OPTIONS_Q[aidx].Log += errBuf.String()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return

	} else {

		CI_OPTIONS_Q[aidx].Log += outBuf.String()

	}
	outBuf = bytes.Buffer{}
	errBuf = bytes.Buffer{}

	kb := pkgresourceci.Builder{}

	kb_c := pkgresourceci.Builder_Container{}

	kb_c_vm := pkgresourceci.Builder_Container_VolumeMount{}

	kb_c_e1 := pkgresourceci.Builder_Container_Env{}
	kb_c_e2 := pkgresourceci.Builder_Container_Env{}

	kb_v := pkgresourceci.Builder_Volume{}

	kb_v_i := pkgresourceci.Builder_Volume_Item{}

	kb.APIVersion = "v1"
	kb.Kind = "Pod"
	kb.Metadata.Name = pod_bid
	kb.Spec.RestartPolicy = "Never"

	kb_c.Name = pod_bid
	kb_c.Image = pkgresourceci.BUILD_EXECUTOR
	kb_c.Args = append(kb_c.Args, "--dockerfile="+pkgresourceci.BUILD_FILE_DEAFULT)
	kb_c.Args = append(kb_c.Args, "--context="+v_repoaddr)
	kb_c.Args = append(kb_c.Args, "--destination="+v_regaddr)

	kb_c_vm.MountPath = "/kaniko/.docker"
	kb_c_vm.Name = "kaniko-secret"

	kb_c_e1.Name = "GIT_USERNAME"
	kb_c_e1.Value = repo_id

	kb_c_e2.Name = "GIT_PASSWORD"
	kb_c_e2.Value = repo_pw

	kb_v.Name = "kaniko-secret"
	kb_v.Secret.SecretName = sec_bid

	kb_v_i.Key = ".dockerconfigjson"
	kb_v_i.Path = "config.json"

	kb_v.Secret.Items = append(kb_v.Secret.Items, kb_v_i)

	kb_c.Env = append(kb_c.Env, kb_c_e1)
	kb_c.Env = append(kb_c.Env, kb_c_e2)

	kb_c.VolumeMounts = append(kb_c.VolumeMounts, kb_c_vm)

	kb.Spec.Containers = append(kb.Spec.Containers, kb_c)
	kb.Spec.Volumes = append(kb.Spec.Volumes, kb_v)

	yb, err := yaml.Marshal(kb)

	if err != nil {
		CI_OPTIONS_Q[aidx].Log = err.Error()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return
	}

	build_yaml, err := V1SaveToCache(bid+".yaml", yb)

	if err != nil {

		CI_OPTIONS_Q[aidx].Log = err.Error()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return

	}

	cmd = exec.Command("kubectl", "-n", ns_bid, "apply", "-f", build_yaml)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CI_OPTIONS_Q[aidx].Log += errBuf.String()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return

	} else {

		CI_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	for {

		outBuf = bytes.Buffer{}
		errBuf = bytes.Buffer{}

		cmd = exec.Command("kubectl", "-n", ns_bid, "get", "pod", pod_bid, "--no-headers")

		cmd.Stdout = &outBuf

		cmd.Stderr = &errBuf

		err = cmd.Run()

		if err != nil {

			CI_OPTIONS_Q[aidx].Log += errBuf.String()

			CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

			return
		}

		stdout_str := outBuf.String()
		_ = errBuf.String()

		if strings.Contains(stdout_str, "Completed") {

			break
		}

		if strings.Contains(stdout_str, "Error") {

			break
		}

		outBuf = bytes.Buffer{}

		errBuf = bytes.Buffer{}

		cmd = exec.Command("kubectl", "-n", ns_bid, "logs", pod_bid)

		cmd.Stdout = &outBuf

		cmd.Stderr = &errBuf

		err = cmd.Run()

		if err != nil {

			continue
		}

		stdout_str = outBuf.String()
		_ = errBuf.String()

		CI_OPTIONS_Q[aidx].BuildLog = stdout_str

		time.Sleep(time.Millisecond * 100)

	}

	outBuf = bytes.Buffer{}

	errBuf = bytes.Buffer{}

	cmd = exec.Command("kubectl", "-n", ns_bid, "logs", pod_bid)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CI_OPTIONS_Q[aidx].Log += errBuf.String()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return

	} else {

		CI_OPTIONS_Q[aidx].BuildLog = outBuf.String()

	}

	outBuf = bytes.Buffer{}
	errBuf = bytes.Buffer{}

	cmd = exec.Command("kubectl", "delete", "namespace", ns_bid)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CI_OPTIONS_Q[aidx].Log += errBuf.String()

		CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_ERROR

		return

	} else {

		CI_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	CI_OPTIONS_Q[aidx].Status = pkgresourceci.STATUS_COMPLETED

}

func V1CiHandler_BuildReport(aidx int, mani *pkgresourceapix.V1Manifest, acon *V1AgentConn) {

	v1main, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindAgentPush, "/project/ci/log", mani)

	if err != nil {

		log.Printf("failed to report build: %s\n", err.Error())

		return
	}

	v1main.Body["name"] = CI_OPTIONS_Q[aidx].ProjectName
	v1main.Body["status"] = string(CI_OPTIONS_Q[aidx].Status)
	v1main.Body["log"] = CI_OPTIONS_Q[aidx].Log + pkgresourceci.BUILD_LOG_SEP + CI_OPTIONS_Q[aidx].BuildLog

	err = V1AgentPush(v1main, acon)

	if err != nil {

		log.Printf("failed to report build: agent push: %s", err.Error())
	}

	return

}

func V1CdHandler(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		cdlen := len(CD_OPTIONS_Q)

		term_list := []int{}

		for i := 0; i < cdlen; i++ {

			if CD_OPTIONS_Q[i].Status == pkgresourcecd.STATUS_READY {

				CD_OPTIONS_Q[i].Status = pkgresourcecd.STATUS_RUNNING

				go V1CdHandler_DeployStart(i)

			} else if CD_OPTIONS_Q[i].Status == pkgresourcecd.STATUS_RUNNING {

				V1CdHandler_DeployReport(i, mani, acon)

			} else {

				V1CdHandler_DeployReport(i, mani, acon)

				term_list = append(term_list, i)
			}

		}

		err := V1CdClear(term_list)

		if err != nil {

			log.Printf("failed to clear cd list: %s\n", err.Error())
		}

	}

}

func V1CdHandler_DeployStart(aidx int) {

	namespace := CD_OPTIONS_Q[aidx].ProjectName

	project_name := CD_OPTIONS_Q[aidx].ProjectName

	service := CD_OPTIONS_Q[aidx].CdOption.Service

	deployment := CD_OPTIONS_Q[aidx].CdOption.Deployment

	v_regaddr := CD_OPTIONS_Q[aidx].Reg
	reg_id := CD_OPTIONS_Q[aidx].RegId
	reg_pw := CD_OPTIONS_Q[aidx].RegPw

	service_b, err := yaml.Marshal(service)

	deployment_b, err := yaml.Marshal(deployment)

	service_file_name := fmt.Sprintf("service-%s.yaml", namespace)

	deployment_file_name := fmt.Sprintf("deployment-%s.yaml", namespace)

	service_yaml, err := V1SaveToCache(service_file_name, service_b)

	if err != nil {

		CD_OPTIONS_Q[aidx].Log += err.Error()

		CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_ERROR

		return

	}

	deployment_yaml, err := V1SaveToCache(deployment_file_name, deployment_b)

	if err != nil {

		CD_OPTIONS_Q[aidx].Log += err.Error()

		CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_ERROR

		return

	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	cmd := exec.Command("kubectl", "create", "namespace", namespace)

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CD_OPTIONS_Q[aidx].Log += errBuf.String()

		CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_ERROR

		return

	} else {

		CD_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	outBuf = bytes.Buffer{}
	errBuf = bytes.Buffer{}

	docker_server := "--docker-server=" + v_regaddr
	docker_uname := "--docker-username=" + reg_id
	docker_pword := "--docker-password=" + reg_pw

	cmd = exec.Command("kubectl", "-n", namespace, "create", "secret", "docker-registry", project_name, docker_server, docker_uname, docker_pword)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CD_OPTIONS_Q[aidx].Log += errBuf.String()

		CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_ERROR

		return

	} else {

		CD_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	outBuf = bytes.Buffer{}
	errBuf = bytes.Buffer{}

	cmd = exec.Command("kubectl", "-n", namespace, "apply", "-f", service_yaml)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CD_OPTIONS_Q[aidx].Log += errBuf.String()

		CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_ERROR

		return

	} else {

		CD_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	outBuf = bytes.Buffer{}
	errBuf = bytes.Buffer{}

	cmd = exec.Command("kubectl", "-n", namespace, "apply", "-f", deployment_yaml)

	cmd.Stdout = &outBuf

	cmd.Stderr = &errBuf

	err = cmd.Run()

	if err != nil {

		CD_OPTIONS_Q[aidx].Log += errBuf.String()

		CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_ERROR

		return

	} else {

		CD_OPTIONS_Q[aidx].Log += outBuf.String()

	}

	CD_OPTIONS_Q[aidx].Status = pkgresourcecd.STATUS_COMPLETED

}

func V1CdHandler_DeployReport(aidx int, mani *pkgresourceapix.V1Manifest, acon *V1AgentConn) {

	v1main, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindAgentPush, "/project/cd/log", mani)

	if err != nil {

		log.Printf("failed to report deploy: %s\n", err.Error())

		return
	}

	v1main.Body["name"] = CD_OPTIONS_Q[aidx].ProjectName
	v1main.Body["status"] = string(CD_OPTIONS_Q[aidx].Status)
	v1main.Body["log"] = CD_OPTIONS_Q[aidx].Log

	err = V1AgentPush(v1main, acon)

	if err != nil {

		log.Printf("failed to report deploy: agent push: %s", err.Error())
	}

	return

}

func V1CiCdHandler(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		pipelen := len(CICD_PIPE_Q)

		term_list := []int{}

		for i := 0; i < pipelen; i++ {

			if CICD_PIPE_Q[i].Ci == nil || CICD_PIPE_Q[i].Cd == nil {

				continue

			}

			if CICD_PIPE_Q[i].Ci.Status == pkgresourceci.STATUS_READY {

				CICD_PIPE_Q[i].Ci.Status = pkgresourceci.STATUS_RUNNING

				go V1CiHandler_BuildStart(i)

			} else if CICD_PIPE_Q[i].Ci.Status == pkgresourceci.STATUS_RUNNING {

				V1CiHandler_BuildReport(i, mani, acon)

			} else if CICD_PIPE_Q[i].Cd.Status == pkgresourcecd.STATUS_READY {

				CICD_PIPE_Q[i].Cd.Status = pkgresourcecd.STATUS_RUNNING

				V1CiHandler_BuildReport(i, mani, acon)

				go V1CdHandler_DeployStart(i)

			} else if CICD_PIPE_Q[i].Cd.Status == pkgresourcecd.STATUS_RUNNING {

				V1CdHandler_DeployReport(i, mani, acon)

			} else {

				V1CdHandler_DeployReport(i, mani, acon)

				term_list = append(term_list, i)
			}

		}

		err := V1CiCdClear(term_list)

		if err != nil {

			log.Printf("failed to clear cicd: %s", err.Error())
		}

	}

}

func V1LifecycleHandler(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		LC_LOCK.Lock()

		lclen := len(LC_MANIFEST_Q)

		for i := 0; i < lclen; i++ {

			lc_mani := LC_MANIFEST_Q[i]

			var report pkgresourcelc.LifecycleReport

			report.Obsolete = false
			report.SentTimestamp = time.Now()
			report.Process = lc_mani.Process

			report_b, err := yaml.Marshal(report)

			if err != nil {

				log.Printf("failed to report lifecycle: %s\n", err.Error())

				continue

			}

			v1main, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindAgentPush, "/project/lifecycle/report", mani)

			if err != nil {

				log.Printf("failed to report lifecycle: v1 main: %s\n", err.Error())

				continue
			}

			v1main.Body["report"] = string(report_b)

			err = V1AgentPush(v1main, acon)

			if err != nil {

				log.Printf("failed to report lifecycle: agent push: %s", err.Error())

				continue

			}

		}

		LC_LOCK.Unlock()

	}

}

func V1LifecycleTerminator(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		LC_TERM_LOCK.Lock()

		lctermlen := len(LC_TERMINATOR_Q)

		for i := 0; i < lctermlen; i++ {

			report := LC_TERMINATOR_Q[i]

			report_b, err := yaml.Marshal(report)

			if err != nil {

				log.Printf("failed to report free lc: %s\n", err.Error())

				continue
			}

			v1main, err := pkgapix.V1GetMainCopyByAddress(pkgresourceapix.V1KindAgentPush, "/project/lifecycle/report", mani)

			if err != nil {

				log.Printf("failed to report free lifecycle: v1 main: %s\n", err.Error())

				continue
			}

			v1main.Body["report"] = string(report_b)

			err = V1AgentPush(v1main, acon)

			if err != nil {

				log.Printf("failed to report free lifecycle: agent push: %s", err.Error())

				continue

			}

		}

		LC_TERMINATOR_Q = make([]pkgresourcelc.LifecycleReport, 0)

		LC_TERM_LOCK.Unlock()
	}

}
