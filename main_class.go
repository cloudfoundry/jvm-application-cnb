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

package jvm_application_buildpack

import (
	"fmt"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/magiconair/properties"
)

// MainClass represents the main class in a JVM application.
type MainClass struct {
	application libbuildpack.Application
	launch      libjavabuildpack.Launch
	logger      libjavabuildpack.Logger
	class       string
}

// Contribute makes the contribution to launch
func (m MainClass) Contribute() error {
	m.logger.FirstLine("Configuring Java Main Application")

	command := fmt.Sprintf("java -cp %s $JAVA_OPTS %s", m.application.Root, m.class)

	return m.launch.WriteMetadata(libbuildpack.LaunchMetadata{
		Processes: []libbuildpack.Process{
			{"web", command},
			{"task", command},
		},
	})
}

// String makes MainClass satisfy the Stringer interface.
func (m MainClass) String() string {
	return fmt.Sprintf("MainClass{ application: %s, launch: %s, logger: %s, class:%s }",
		m.application, m.launch, m.logger, m.class)
}

// HasMainClass returns true if the application contains a META-INF/MANIFEST.MF file with a "Main-Class" key in it,
// otherwise false.
func HasMainClass(application libbuildpack.Application, logger libjavabuildpack.Logger) (bool, error) {
	m, ok, err := newManifest(application, logger)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	_, ok = m.Get("Main-Class")
	return ok, nil
}

// NewMainCallas creates a new MainClass instance.  OK is true if the application contains a META-INF/MANIFEST.MF file
// with a "Main-Class" key in it.
func NewMainClass(build libjavabuildpack.Build) (MainClass, bool, error) {
	m, ok, err := newManifest(build.Application, build.Logger)
	if err != nil {
		return MainClass{}, false, err
	}
	if !ok {
		return MainClass{}, false, nil
	}

	c, ok := m.Get("Main-Class")
	if !ok {
		return MainClass{}, false, nil
	}

	return MainClass{
		build.Application,
		build.Launch,
		build.Logger,
		c,
	}, true, nil
}

func newManifest(application libbuildpack.Application, logger libjavabuildpack.Logger) (*properties.Properties, bool, error) {
	manifest := filepath.Join(application.Root, "META-INF", "MANIFEST.MF")

	exists, err := libjavabuildpack.FileExists(manifest)
	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, false, nil
	}

	p, err := properties.LoadFile(manifest, properties.UTF8)
	if err != nil {
		return nil, false, err
	}

	logger.Debug("Manifest: %s", p)
	return p, true, nil
}
