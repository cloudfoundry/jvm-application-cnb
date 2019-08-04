/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-cnb/jre"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("indeterminant", func() {
			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Requires: []buildplan.Required{
					{Name: jre.Dependency, Metadata: buildplan.Metadata{jre.LaunchContribution: true}},
					{Name: jvmapplication.Dependency},
				},
			}))
		})

		it("classes", func() {
			test.TouchFile(t, f.Detect.Application.Root, "test.class")

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jvmapplication.Dependency},
				},
				Requires: []buildplan.Required{
					{Name: jre.Dependency, Metadata: buildplan.Metadata{jre.LaunchContribution: true}},
					{Name: jvmapplication.Dependency},
				},
			}))
		})

		it("groovies", func() {
			test.TouchFile(t, f.Detect.Application.Root, "test.groovy")

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jvmapplication.Dependency},
				},
				Requires: []buildplan.Required{
					{Name: jre.Dependency, Metadata: buildplan.Metadata{jre.LaunchContribution: true}},
					{Name: jvmapplication.Dependency},
				},
			}))
		})

		it("jars", func() {
			test.TouchFile(t, f.Detect.Application.Root, "test.jar")

			g.Expect(d(f.Detect)).To(gomega.Equal(detect.PassStatusCode))
			g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: jvmapplication.Dependency},
				},
				Requires: []buildplan.Required{
					{Name: jre.Dependency, Metadata: buildplan.Metadata{jre.LaunchContribution: true}},
					{Name: jvmapplication.Dependency},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
