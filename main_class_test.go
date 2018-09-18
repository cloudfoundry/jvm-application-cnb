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

package jvm_application_buildpack_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/jvm-application-buildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMainClass(t *testing.T) {
	spec.Run(t, "MainClass", testMainClass, spec.Report(report.Terminal{}))
}

func testMainClass(t *testing.T, when spec.G, it spec.S) {

	when("HasMainClass", func() {

		logger := libjavabuildpack.Logger{}

		it("returns false when no manifest", func() {
			root := test.ScratchDir(t, "main-class")

			ok, err := jvm_application_buildpack.HasMainClass(libbuildpack.Application{Root: root}, logger)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("HasMainClass = %t, expected false", ok)
			}
		})

		it("returns false when no Main-Class", func() {
			root := test.ScratchDir(t, "main-class")

			if err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(root, "META-INF", "MANIFEST.MF"), 0644); err != nil {
				t.Fatal(err)
			}

			ok, err := jvm_application_buildpack.HasMainClass(libbuildpack.Application{Root: root}, logger)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("HasMainClass = %t, expected false", ok)
			}
		})

		it("returns true when Main-Class exists", func() {
			root := test.ScratchDir(t, "main-class")

			if err := libjavabuildpack.WriteToFile(strings.NewReader("Main-Class: test-class"), filepath.Join(root, "META-INF", "MANIFEST.MF"), 0644); err != nil {
				t.Fatal(err)
			}

			ok, err := jvm_application_buildpack.HasMainClass(libbuildpack.Application{Root: root}, logger)
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Errorf("HasMainClass = %t, expected true", ok)
			}
		})

	})

	when("NewMainClass", func() {

		it("returns false when no manifest", func() {
			f := test.NewBuildFactory(t)

			_, ok, err := jvm_application_buildpack.NewMainClass(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("NewMainClass = %t, expected false", ok)
			}
		})

		it("returns false when no Main-Class", func() {
			f := test.NewBuildFactory(t)

			m := filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF")
			if err := libjavabuildpack.WriteToFile(strings.NewReader(""), m, 0644); err != nil {
				t.Fatal(err)
			}

			_, ok, err := jvm_application_buildpack.NewMainClass(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if ok {
				t.Errorf("NewMainClass = %t, expected false", ok)
			}
		})

		it("returns true when Main-Class exists", func() {
			f := test.NewBuildFactory(t)

			m := filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF")

			if err := libjavabuildpack.WriteToFile(strings.NewReader("Main-Class: test-class"), m, 0644); err != nil {
				t.Fatal(err)
			}

			_, ok, err := jvm_application_buildpack.NewMainClass(f.Build)
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

		m := filepath.Join(f.Build.Application.Root, "META-INF", "MANIFEST.MF")

		if err := libjavabuildpack.WriteToFile(strings.NewReader("Main-Class: test-class"), m, 0644); err != nil {
			t.Fatal(err)
		}

		c, _, err := jvm_application_buildpack.NewMainClass(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := c.Contribute(); err != nil {
			t.Fatal(err)
		}

		var actual libbuildpack.LaunchMetadata
		_, err = toml.DecodeFile(filepath.Join(f.Build.Launch.Root, "launch.toml"), &actual)
		if err != nil {
			t.Fatal(err)
		}

		command := fmt.Sprintf("java -cp %s $JAVA_OPTS test-class", f.Build.Application.Root)

		expected := libbuildpack.LaunchMetadata{
			Processes: libbuildpack.Processes{
				libbuildpack.Process{Type: "web", Command: command},
				libbuildpack.Process{Type: "task", Command: command},
			},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("launch.toml = %s, expected %s", actual, expected)
		}
	})
}
