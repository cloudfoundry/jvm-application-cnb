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
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack/jre"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("fails without Main-Class", func() {
			g.Expect(d(f.Detect)).To(Equal(detect.FailStatusCode))
		})

		it("passes with jvm-application", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})

			g.Expect(d(f.Detect)).To(Equal(detect.PassStatusCode))
			g.Expect(f.Output).To(Equal(buildplan.BuildPlan{
				jvmapplication.Dependency: buildplan.Dependency{},
				jre.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{jre.LaunchContribution: true},
				},
			}))
		})

		it("passes with Main-Class", func() {
			test.WriteFile(t, filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

			g.Expect(d(f.Detect)).To(Equal(detect.PassStatusCode))
			g.Expect(f.Output).To(Equal(buildplan.BuildPlan{
				jvmapplication.Dependency: buildplan.Dependency{},
				jre.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{jre.LaunchContribution: true},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
