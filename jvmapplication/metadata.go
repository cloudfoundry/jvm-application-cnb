/*
 * Copyright 2018-2019 the original author or authors.
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

package jvmapplication

import (
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/openjdk-cnb/jre"
)

type Metadata struct{}

// BuildPlan returns a BuildPlan that contains dependencies for jvm-application and openjdk-jre.
func (m Metadata) BuildPlan(buildPlan buildplan.BuildPlan) buildplan.BuildPlan {
	return buildplan.BuildPlan{
		Dependency:     buildPlan[Dependency],
		jre.Dependency: m.jre(buildPlan),
	}
}

// String makes Metadata satisfy the Stringer interface.
func (Metadata) String() string {
	return "Metadata{}"
}

func (Metadata) jre(buildPlan buildplan.BuildPlan) buildplan.Dependency {
	d := buildPlan[jre.Dependency]

	if d.Metadata == nil {
		d.Metadata = make(buildplan.Metadata)
	}

	d.Metadata[jre.LaunchContribution] = true

	return d
}
