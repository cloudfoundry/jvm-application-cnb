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

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/jvm-application-buildpack/mainclass"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/openjdk-buildpack/jre"
)

func main() {
	detect, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Detect: %s\n", err)
		os.Exit(101)
	}

	if err := detect.BuildPlan.Init(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Build Plan: %s\n", err)
		os.Exit(101)
	}

	if code, err := d(detect); err != nil {
		detect.Logger.Info(err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func d(detect detect.Detect) (int, error) {
	_, dep := detect.BuildPlan[jvmapplication.Dependency]

	mc, err := mainclass.HasMainClass(detect.Application, detect.Logger)
	if err != nil {
		return detect.Error(102), err
	}

	if dep || mc {
		return detect.Pass(buildplan.BuildPlan{
			jvmapplication.Dependency: buildplan.Dependency{},
			jre.Dependency: buildplan.Dependency{
				Metadata: buildplan.Metadata{jre.LaunchContribution: true},
			},
		})
	}

	return detect.Fail(), nil
}
