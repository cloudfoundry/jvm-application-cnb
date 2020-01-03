/*
 * Copyright 2019-2020 the original author or authors.
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

package executablejar_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/jvm-application-cnb/executablejar"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestExecutableJAR(t *testing.T) {
	spec.Run(t, "ExecutableJAR", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		when("NewExecutableJAR", func() {

			it("returns false when no Main-Class", func() {
				test.WriteFile(t, filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF"), "")

				_, ok, err := executablejar.NewExecutableJAR(f.Build)
				g.Expect(ok).To(gomega.BeFalse())
				g.Expect(err).NotTo(gomega.HaveOccurred())
			})

			it("returns true when Main-Class exists", func() {
				test.WriteFile(t, filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

				_, ok, err := executablejar.NewExecutableJAR(f.Build)
				g.Expect(ok).To(gomega.BeTrue())
				g.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		it("contributes command", func() {
			test.WriteFile(t, filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF"), "Main-Class: test-class")

			e, ok, err := executablejar.NewExecutableJAR(f.Build)
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(e.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("executable-jar")
			g.Expect(layer).To(test.HaveLayerMetadata(true, true, true))
			g.Expect(layer).To(test.HavePrependPathSharedEnvironment("CLASSPATH", f.Build.Application.Root))

			command := "java -cp $CLASSPATH $JAVA_OPTS test-class"
			g.Expect(f.Build.Layers).To(test.HaveApplicationMetadata(layers.Metadata{
				Processes: []layers.Process{
					{Type: "executable-jar", Command: command},
					{Type: "task", Command: command},
					{Type: "web", Command: command},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
