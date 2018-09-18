/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/jvm-application-buildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/openjdk-buildpack"
)

func main() {
	detect, err := libjavabuildpack.DefaultDetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Detect: %s\n", err.Error())
		os.Exit(101)
	}

	_, j := detect.BuildPlan[jvm_application_buildpack.JVMApplication]

	m, err := jvm_application_buildpack.HasMainClass(detect.Application, detect.Logger)
	if err != nil {
		detect.Error(102)
		return
	}

	if j || m {
		detect.Pass(libbuildpack.BuildPlan{
			jvm_application_buildpack.JVMApplication: libbuildpack.BuildPlanDependency{},
			openjdk_buildpack.JREDependency: libbuildpack.BuildPlanDependency{
				Metadata: libbuildpack.BuildPlanDependencyMetadata{
					"version":                            "1.*",
					openjdk_buildpack.LaunchContribution: true,
				},
			},
		})
		return
	}

	detect.Fail()
	return
}
