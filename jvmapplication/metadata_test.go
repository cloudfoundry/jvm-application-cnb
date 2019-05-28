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

package jvmapplication_test

import (
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/openjdk-cnb/jre"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMetadata(t *testing.T) {
	spec.Run(t, "Metadata", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("returns build plan entries", func() {
			g.Expect(jvmapplication.Metadata{}.BuildPlan(buildplan.BuildPlan{})).To(Equal(buildplan.BuildPlan{
				jvmapplication.Dependency: buildplan.Dependency{},
				jre.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{jre.LaunchContribution: true},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
