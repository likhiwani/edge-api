/*
	Copyright NetFoundry Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"context"
	"crypto/x509"
	"time"
	"ztna-core/edge-api/rest_management_api_client/config"
	"ztna-core/edge-api/rest_model"
	"ztna-core/edge-api/rest_util"

	log "github.com/sirupsen/logrus"
)

// main is an example usage of the rest_util package. This example is missing the concepts of enrollment.
// Part of enrollment is being delivered a signed JWT that can be used to verify the controller server certificate
// That step is missing from this example. The CA bundle from the well-known endpoint is verified as a sanity
// check against the controller. However, this does not add any extra security, just sanity.
func main() {
	ctrlAddress := "https://localhost:8444"
	caCerts, err := rest_util.GetControllerWellKnownCas(ctrlAddress)

	if err != nil {
		log.Fatal(err)
	}

	caPool := x509.NewCertPool()

	for _, ca := range caCerts {
		caPool.AddCert(ca)
		// pemCert := pem.EncodeToMemory(&pem.Block{
        //     Type:  "CERTIFICATE",
        //     Bytes: ca.Raw,
        // })
        // println(string(pemCert))
	}

	ok, err := rest_util.VerifyController(ctrlAddress, caPool)

	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("controller failed CA validation")
	}

	client, err := rest_util.NewEdgeManagementClientWithUpdb("likhiwani", "likhiwani", ctrlAddress, caPool)

	if err != nil {
		log.Fatal(err)
	}

	// params := &identity.ListIdentitiesParams{
	// 	Context: context.Background(),
	// }
	// {
	// 	"forwardAddress": true,
	// 	"allowedAddresses": [
	// 		"*"
	// 	],
	// 	"forwardPort": true,
	// 	"allowedPortRanges": [
	// 		{
	// 			"low": 5000,
	// 			"high": 2048
	// 		}
	// 	],
	// 	"forwardProtocol": false,
	// 	"allowedProtocols": [
	// 		"tcp",
	// 		"udp"
	// 	]
	// }

	protocols := []string{"tcp", "udp"}
    addresses := []string{"*"}
    portRangeLow := 1000
    portRangeHigh := 2000

	interceptV1Data := map[string]interface{}{
		"protocols": protocols,
		"addresses": addresses,
		"portRanges": []map[string]interface{}{
			{
				"low":  portRangeLow,
				"high": portRangeHigh,
			},
		},
	}

	confCreate := &rest_model.ConfigCreate{
		ConfigTypeID: ptr("g7cIWbcGg"),
		Name: ptr("edge-api-config-1"),
		Data: &interceptV1Data,
	}
	configParams := &config.CreateConfigParams{
		Context: context.Background(),
		Config: confCreate,
	}

	configParams.SetTimeout(30 * time.Second)
	resp, err := client.Config.CreateConfig(configParams, nil)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not create Intercept.V1 service config :(")
	}

	println("\n=== New Config is created ===", resp)
}

func ptr(s string) *string {
	return &s
}
