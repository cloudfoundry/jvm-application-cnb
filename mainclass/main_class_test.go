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

package mainclass_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/jvm-application-buildpack/mainclass"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMainClass(t *testing.T) {
	spec.Run(t, "MainClass", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		when("HasMainClass", func() {

			var f *test.DetectFactory

			it.Before(func() {
				f = test.NewDetectFactory(t)
			})

			it("returns false when no manifest", func() {
				g.Expect(mainclass.HasMainClass(f.Detect.Application, f.Detect.Logger)).To(BeFalse())
			})

			it("returns false when no Main-Class", func() {
				test.TouchFile(t, f.Detect.Application.Root, "META-INF", "MANIFEST.MF")

				g.Expect(mainclass.HasMainClass(f.Detect.Application, f.Detect.Logger)).To(BeFalse())
			})

			it("returns true when Main-Class exists", func() {
				test.WriteFile(t, filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

				g.Expect(mainclass.HasMainClass(f.Detect.Application, f.Detect.Logger)).To(BeTrue())
			})
		})

		when("NewMainClass", func() {

			it("returns false when no jvm-application", func() {
				test.WriteFile(t, filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

				_, ok, err := mainclass.NewMainClass(f.Build)
				g.Expect(ok).To(BeFalse())
				g.Expect(err).NotTo(HaveOccurred())
			})

			it("returns false when no main-Class", func() {
				f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})

				_, ok, err := mainclass.NewMainClass(f.Build)
				g.Expect(ok).To(BeFalse())
				g.Expect(err).NotTo(HaveOccurred())
			})

			it("returns true when main-Class exists", func() {
				f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
				test.WriteFile(t, filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

				_, ok, err := mainclass.NewMainClass(f.Build)
				g.Expect(ok).To(BeTrue())
				g.Expect(err).NotTo(HaveOccurred())
			})
		})

		it("contributes command", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			test.WriteFile(t, filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

			c, _, err := mainclass.NewMainClass(f.Build)
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(c.Contribute()).To(Succeed())

			command := fmt.Sprintf("java -cp %s $JAVA_OPTS test-class", f.Build.Application.Root)

			g.Expect(f.Build.Layers).To(test.HaveApplicationMetadata(layers.Metadata{
				Processes: []layers.Process{
					{"task", command},
					{"web", command},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
