/*
 * Copyright 2018-2020 the original author or authors.
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

package executablejar

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
	"github.com/mitchellh/mapstructure"
)

// Dependency indicates that an application is an executable JAR.
const Dependency = "executable-jar"

// ExecutableJAR represents an executable JAR JVM application.
type ExecutableJAR struct {
	// Metadata is metadata about the executable JAR application.
	Metadata Metadata

	layer  layers.Layer
	layers layers.Layers
}

// Contribute makes the contribution to launch.
func (e ExecutableJAR) Contribute() error {
	if err := e.layer.Contribute(e.Metadata, func(layer layers.Layer) error {
		return layer.PrependPathSharedEnv("CLASSPATH", strings.Join(e.Metadata.ClassPath, string(filepath.ListSeparator)))
	}, layers.Build, layers.Cache, layers.Launch); err != nil {
		return err
	}

	command := fmt.Sprintf("java -cp $CLASSPATH $JAVA_OPTS %s", e.Metadata.MainClass)

	return e.layers.WriteApplicationMetadata(layers.Metadata{
		Processes: layers.Processes{
			{Type: "executable-jar", Command: command},
			{Type: "task", Command: command},
			{Type: "web", Command: command},
		},
	})
}

// Plan returns the dependency information for this application.
func (e ExecutableJAR) Plan() (buildpackplan.Plan, error) {
	p := buildpackplan.Plan{
		Name:     Dependency,
		Metadata: make(buildpackplan.Metadata),
	}

	if err := mapstructure.Decode(e.Metadata, &p.Metadata); err != nil {
		return buildpackplan.Plan{}, err
	}

	return p, nil
}

// NewExecutableJAR creates a new ExecutableJAR instance.  OK is true if the build plan contains a "jvm-application"
// dependency and a "Main-Class" manifest key.
func NewExecutableJAR(build build.Build) (ExecutableJAR, bool, error) {
	md, ok, err := NewMetadata(build.Application, build.Logger)
	if err != nil {
		return ExecutableJAR{}, false, err
	}

	if !ok {
		return ExecutableJAR{}, false, nil
	}

	return ExecutableJAR{
		md,
		build.Layers.Layer(Dependency),
		build.Layers,
	}, true, nil
}
