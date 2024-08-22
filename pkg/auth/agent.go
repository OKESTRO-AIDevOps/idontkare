package challenge

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	pkgdbquery "github.com/OKESTRO-AIDevOps/idontkare/pkg/dbquery"
	pkgresourceauth "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/auth"
	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
)

func AgentAuthGenerateChallenge(username string, clustername string) (string, string, error) {

	var chalKey string

	userRecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return "", "", fmt.Errorf("agent auth: failed to get user: %s", err.Error())

	}

	if userRecord == nil {

		return "", "", fmt.Errorf("agent auth: no such user: %s", username)

	}

	clusterRecord, err := pkgdbquery.GetClustersByUserId(userRecord.UserId)

	if err != nil {

		return "", "", fmt.Errorf("agent auth: failed to get cluster: %s", err.Error())

	}

	if clusterRecord == nil {

		return "", "", fmt.Errorf("agent auth: no associated cluster for: %s", username)
	}

	cLen := len(clusterRecord)

	found := -1

	for i := 0; i < cLen; i++ {

		if clusterRecord[i].ClusterName == clustername {

			found = i

			break
		}

	}

	if found == -1 {

		return "", "", fmt.Errorf("agent auth: no such cluster: %s", clustername)
	}

	pKeyStr := clusterRecord[found].ClusterPub

	pKey, _ := pkgutils.BytesToPublicKey([]byte(pKeyStr))

	chalKey, _ = pkgutils.RandomHex(32)

	chalData := pkgresourceauth.ChallengeData{

		Key: chalKey,
	}

	chalData_b, _ := json.Marshal(chalData)

	chalDataEnc, _ := pkgutils.EncryptWithPublicKey(chalData_b, pKey)

	chalDataHex := hex.EncodeToString(chalDataEnc)

	return chalKey, chalDataHex, nil

}

func AgentAuthVerifyChallenge(username string, clustername string, chal_data_hex string, key string) (string, error) {

	var newSessionKey string

	var userPass string

	userRecord, err := pkgdbquery.GetUserByName(username)

	if err != nil {

		return "", fmt.Errorf("agent auth verify: failed to get user: %s", err.Error())

	}

	if userRecord == nil {

		return "", fmt.Errorf("agent auth verify: no such user: %s", username)

	}

	userPass = userRecord.UserPass

	key_b := []byte(key)

	chal_data_enc, err := hex.DecodeString(chal_data_hex)

	if err != nil {

		return "", fmt.Errorf("agent auth verify: failed to decode hex: %s", err.Error())
	}

	chalData_b, err := pkgutils.DecryptWithSymmetricKey(key_b, chal_data_enc)

	if err != nil {

		return "", fmt.Errorf("agent auth verify: failed to decrypt: %s", err.Error())
	}

	var chalData pkgresourceauth.ChallengeData

	err = json.Unmarshal(chalData_b, &chalData)

	if err != nil {

		return "", fmt.Errorf("agent auth verify: failed to unmarshal: %s", err.Error())

	}

	if chalData.Pass != userPass {

		return "", fmt.Errorf("agent auth verify: pass doesn't match")
	}

	newSessionKey = chalData.Key

	err = VerifySessionKey(newSessionKey)

	if err != nil {

		return "", fmt.Errorf("agent auth verify: invalid new key: %s", err.Error())
	}

	return newSessionKey, nil

}

func VerifySessionKey(key string) error {

	keylen := len(key)

	if keylen < 32 {

		return fmt.Errorf("keylen too short: %d", keylen)

	} else if keylen > 1024 {

		return fmt.Errorf("keylen too long: %d", keylen)

	}

	for _, c := range key {

		if !pkgutils.CheckIfSliceContains[rune](pkgutils.HEX_RUNES, c) {

			return fmt.Errorf("invalid key")

		}

	}

	return nil
}
