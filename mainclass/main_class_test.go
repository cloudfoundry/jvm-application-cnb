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

package mainclass_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/jvm-application-buildpack/mainclass"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMainClass(t *testing.T) {
	spec.Run(t, "MainClass", testMainClass, spec.Report(report.Terminal{}))
}

func testMainClass(t *testing.T, when spec.G, it spec.S) {

	when("GetMainClass", func() {

		it("returns false when no manifest", func() {
			f := test.NewDetectFactory(t)

			_, ok, err := mainclass.GetMainClass(f.Detect.Application, f.Detect.Logger)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("GetMainClass = %t, expected false", ok)
			}
		})

		it("returns false when no Main-Class", func() {
			f := test.NewDetectFactory(t)

			if err := layers.WriteToFile(strings.NewReader(""), filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"), 0644); err != nil {
				t.Fatal(err)
			}

			_, ok, err := mainclass.GetMainClass(f.Detect.Application, f.Detect.Logger)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("GetMainClass = %t, expected false", ok)
			}
		})

		it("returns true when Main-Class exists", func() {
			f := test.NewDetectFactory(t)

			if err := layers.WriteToFile(strings.NewReader("Main-Class: test-class"), filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"), 0644); err != nil {
				t.Fatal(err)
			}

			class, ok, err := mainclass.GetMainClass(f.Detect.Application, f.Detect.Logger)
			if err != nil {
				t.Fatal(err)
			}

			if class != "test-class" {
				t.Errorf("GetMainClass = %s, expected test-class", class)
			}

			if !ok {
				t.Errorf("GetMainClass = %t, expected true", ok)
			}
		})
	})

	when("NewMainClass", func() {

		it("returns false when no jvm-application", func() {
			f := test.NewBuildFactory(t)

			_, ok, err := mainclass.NewMainClass(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("NewMainClass = %t, expected false", ok)
			}
		})

		it("returns false when no main-Class", func() {
			f := test.NewBuildFactory(t)
			f.AddBuildPlan(t, jvmapplication.Dependency, buildplan.Dependency{})

			_, ok, err := mainclass.NewMainClass(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("NewMainClass = %t, expected false", ok)
			}
		})

		it("returns true when main-Class exists", func() {
			f := test.NewBuildFactory(t)

			f.AddBuildPlan(t, jvmapplication.Dependency, buildplan.Dependency{Metadata: buildplan.Metadata{
				mainclass.MainClassContribution: "test-class",
			}})

			_, ok, err := mainclass.NewMainClass(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Errorf("NewMainClass = %t, expected true", ok)
			}
		})
	})

	it("contributes command", func() {
		f := test.NewBuildFactory(t)
		f.AddBuildPlan(t, jvmapplication.Dependency, buildplan.Dependency{Metadata: buildplan.Metadata{
			mainclass.MainClassContribution: "test-class",
		}})

		c, _, err := mainclass.NewMainClass(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := c.Contribute(); err != nil {
			t.Fatal(err)
		}

		command := fmt.Sprintf("java -cp %s $JAVA_OPTS test-class", f.Build.Application.Root)

		test.BeLaunchMetadataLike(t, f.Build.Layers, layers.Metadata{
			Processes: []layers.Process{
				{"web", command},
				{"task", command},
			},
		})
	})
}
