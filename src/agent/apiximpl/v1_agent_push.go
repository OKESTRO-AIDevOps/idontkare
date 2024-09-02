package apiximpl

import (
	"log"

	pkgresourceapix "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/apix"
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
	pkgresourceci "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/ci"
)

func V1CiHandler(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		cilen := len(CI_OPTIONS_Q)

		term_list := []int{}

		for i := 0; i < cilen; i++ {

			if CI_OPTIONS_Q[i].Status == pkgresourceci.STATUS_READY {

				go V1CiHandler_BuildStart(&CI_OPTIONS_Q[i])

			} else if CI_OPTIONS_Q[i].Status == pkgresourceci.STATUS_RUNNING {

				V1CiHandler_BuildReport(&CI_OPTIONS_Q[i])

			} else {

				term_list = append(term_list, i)

			}

		}

		err := V1CiClear(term_list)

		if err != nil {

			log.Printf("failed to clear ci: %s", err.Error())
		}

	}

}

func V1CiHandler_BuildStart(agent_ci *V1AgentCi) {

}

func V1CiHandler_BuildReport(agent_ci *V1AgentCi) {

}

func V1CdHandler(acon *V1AgentConn, mani *pkgresourceapix.V1Manifest) {

	for {

		cdlen := len(CD_OPTIONS_Q)

		term_list := []int{}

		for i := 0; i < cdlen; i++ {

			if CD_OPTIONS_Q[i].Status == pkgresourcecd.STATUS_READY {

				go V1CdHandler_DeployStart(&CD_OPTIONS_Q[i])

			} else if CD_OPTIONS_Q[i].Status == pkgresourcecd.STATUS_RUNNING {

				V1CdHandler_DeployReport(&CD_OPTIONS_Q[i])

			} else {

				term_list = append(term_list, i)
			}

		}

		err := V1CdClear(term_list)

		if err != nil {

			log.Printf("failed to clear cd list: %s\n", err.Error())
		}

	}

}

func V1CdHandler_DeployStart(agent_cd *V1AgentCd) {

}

func V1CdHandler_DeployReport(agent_cd *V1AgentCd) {

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

				go V1CiHandler_BuildStart(CICD_PIPE_Q[i].Ci)

			} else if CICD_PIPE_Q[i].Ci.Status == pkgresourceci.STATUS_RUNNING {

				V1CiHandler_BuildReport(CICD_PIPE_Q[i].Ci)

			} else if CICD_PIPE_Q[i].Cd.Status == pkgresourcecd.STATUS_READY {

				go V1CdHandler_DeployStart(CICD_PIPE_Q[i].Cd)

			} else if CICD_PIPE_Q[i].Cd.Status == pkgresourcecd.STATUS_RUNNING {

				V1CdHandler_DeployReport(CICD_PIPE_Q[i].Cd)

			} else {

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

		lclen := len(LC_MANIFEST_Q)

		for i := 0; i < lclen; i++ {

		}

	}

}
