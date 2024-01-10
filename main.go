/*
 *   Copyright (c) 2023 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/google/uuid"
	"github.com/intel/trustauthority-client/go-connector"
	"github.com/intel/trustauthority-client/go-sgx"
	"github.com/pkg/errors"
)

// #cgo CFLAGS: -I/opt/intel/sgxsdk/include -fstack-protector-strong
// #cgo LDFLAGS: -lsgx_urts -lutils -Lenclave
// #include "sgx_urts.h"
// #include "enclave/Enclave_u.h"
// #include "enclave/utils.h"
import "C"

type Config struct {
	TrustAuthorityUrl    string `json:"trustauthority_url"`
	TrustAuthorityApiUrl string `json:"trustauthority_api_url"`
	TrustAuthorityApiKey string `json:"trustauthority_api_key"`
}

func main() {
	var policyId string
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Config file containing trustauthority details in JSON format")
	flag.StringVar(&policyId, "pid", "", "Policy id for verification")
	flag.Parse()

	configJson, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		panic(err)
	}

	if config.TrustAuthorityUrl == "" || config.TrustAuthorityApiUrl == "" || config.TrustAuthorityApiKey == "" {
		fmt.Println("Either Trust Authority URL, API URL or API Key is missing in config")
		os.Exit(1)
	}

	cfg := &connector.Config{
		TlsCfg: &tls.Config{
			InsecureSkipVerify: true,
		},
		BaseUrl: config.TrustAuthorityUrl,
		ApiUrl:  config.TrustAuthorityApiUrl,
		ApiKey:  config.TrustAuthorityApiKey,
	}

	trustAuthorityConnector, err := connector.New(cfg)
	if err != nil {
		panic(err)
	}

	eid, err := createSgxEnclave("enclave/enclave.signed.so")
	if err != nil {
		panic(err)
	}

	pubBytes, err := loadPublicKey(eid)
	if err != nil {
		panic(err)
	}

	adapter, err := sgx.NewEvidenceAdapter(eid, pubBytes, unsafe.Pointer(C.enclave_create_report))
	if err != nil {
		panic(err)
	}

	var policyIds []uuid.UUID
	if policyId != "" {
		policyIds = append(policyIds, uuid.MustParse(policyId))
	}

	req := connector.AttestArgs{
		Adapter:   adapter,
		PolicyIds: policyIds,
	}
	resp, err := trustAuthorityConnector.Attest(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nTOKEN: %s\n", string(resp.Token))

	token, err := trustAuthorityConnector.VerifyToken(string(resp.Token))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nCLAIMS: %+v\n", token.Claims)
}

func loadPublicKey(eid uint64) ([]byte, error) {
	// keySize holds the length of the key byte array returned from enclave
	var keySize C.uint32_t

	// keyBuf holds the bytes array of the key returned from enclave
	var keyBuf *C.uint8_t

	ret := C.get_public_key(C.ulong(eid), &keyBuf, &keySize)
	if ret != 0 {
		return nil, errors.New("failed to retrieve key from sgx enclave")
	}

	key := C.GoBytes(unsafe.Pointer(keyBuf), C.int(keySize))
	C.free_public_key(keyBuf)

	return key, nil
}

func createSgxEnclave(enclavePath string) (uint64, error) {
	var status C.sgx_status_t
	eid := C.sgx_enclave_id_t(0)
	updated := C.int(0)
	token := C.sgx_launch_token_t{}

	status = C.sgx_create_enclave(C.CString(enclavePath),
		0,
		&token,
		&updated,
		&eid,
		nil)

	if status != 0 {
		return 0, errors.Errorf("Failed to create enclave: %x", status)
	}

	return uint64(eid), nil
}
