// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package main

import (
	"fmt"
	"time"
)

type SetupPipeline struct {
	kube_client       *KubeClient
	create_devcluster bool
	err               error
}

type setupStep func(sp *SetupPipeline)

func (sp *SetupPipeline) step(f setupStep) {

	if sp.err != nil {
		return
	}
	f(sp)
	time.Sleep(2 * time.Second)
}

func setupNuvolaris(args []string) error {
	sp := SetupPipeline{}
	if len(args) > 0 {
		if args[0] == "--devcluster" {
			sp.create_devcluster = true
		} else {
			fmt.Println("did you mean nuv setup --devcluster?")
			return nil
		}
	}
	sp.step(assertNuvolarisClusterConfig)
	sp.step(createNuvolarisNamespace)
	sp.step(deployWhiskCrd)
	sp.step(deployServiceAccount)
	sp.step(deployClusterRoleBinding)
	sp.step(setupWskProperties)
	sp.step(runNuvolarisOperatorPod)
	sp.step(deployOperatorObject)
	sp.step(waitForOpenWhiskReady)
	return sp.err
}

func assertNuvolarisClusterConfig(sp *SetupPipeline) {
	sp.kube_client, sp.err = initClients(sp.create_devcluster)
}

func createNuvolarisNamespace(sp *SetupPipeline) {
	sp.err = sp.kube_client.createNuvNamespace()
}

func deployWhiskCrd(sp *SetupPipeline) {
	sp.err = sp.kube_client.deployCRD()
}

func deployServiceAccount(sp *SetupPipeline) {
	sp.err = sp.kube_client.createServiceAccount()
}

func deployClusterRoleBinding(sp *SetupPipeline) {
	sp.err = sp.kube_client.createClusterRoleBinding()
}

func runNuvolarisOperatorPod(sp *SetupPipeline) {
	sp.err = sp.kube_client.createOperatorPod()
}

func setupWskProperties(sp *SetupPipeline) {
	sp.err = writeWskPropertiesFile()
}

func deployOperatorObject(sp *SetupPipeline) {
	sp.err = createWhiskOperatorObject(sp.kube_client.cfg)
}

func waitForOpenWhiskReady(sp *SetupPipeline) {
	sp.err = readinessProbe()
}
