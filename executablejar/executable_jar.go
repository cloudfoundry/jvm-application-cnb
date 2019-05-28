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

package executablejar

import (
	"fmt"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/manifest"
)

const (
	// Dependency indicates that an application is an executable JAR.
	Dependency = "executable-jar"

	// MainClass indicates the Main-Class of an executable JAR.
	MainClass = "main-class"
)

// ExecutableJAR represents the an executable JAR JVM application.
type ExecutableJAR struct {
	application application.Application
	class       string
	layer       layers.Layer
	layers      layers.Layers
	logger      logger.Logger
}

// Contribute makes the contribution to launch
func (e ExecutableJAR) Contribute() error {
	if err := e.layer.Contribute(marker(e.application.Root), func(layer layers.Layer) error {
		return layer.AppendPathSharedEnv("CLASSPATH", e.application.Root)
	}, layers.Build, layers.Cache, layers.Launch); err != nil {
		return err
	}

	command := fmt.Sprintf("java -cp $CLASSPATH $JAVA_OPTS %s", e.class)

	return e.layers.WriteApplicationMetadata(layers.Metadata{
		Processes: layers.Processes{
			{"executable-jar", command},
			{"task", command},
			{"web", command},
		},
	})
}

// String makes ExecutableJAR satisfy the Stringer interface.
func (e ExecutableJAR) String() string {
	return fmt.Sprintf("ExecutableJAR{ application: %s, class:%s, layer: %s, layers: %s, logger: %s }",
		e.application, e.class, e.layer, e.layers, e.logger)
}

// NewExecutableJAR creates a new ExecutableJAR instance.  OK is true if the build plan contains either a
// "executable-jar" dependency or a "jvm-application" dependency and a "Main-Class" manifest key.
func NewExecutableJAR(build build.Build) (ExecutableJAR, bool, error) {
	_, ok := build.BuildPlan[jvmapplication.Dependency]
	if !ok {
		return ExecutableJAR{}, false, nil
	}

	var class string

	e, ok := build.BuildPlan[Dependency]
	if ok {
		m, ok := e.Metadata[MainClass]
		if !ok {
			return ExecutableJAR{}, false, fmt.Errorf("executable-jar dependency must have main-class metadata")
		}

		if class, ok = m.(string); !ok {
			return ExecutableJAR{}, false, fmt.Errorf("main-class metadata must be a string")
		}
	} else {
		m, err := manifest.NewManifest(build.Application, build.Logger)
		if err != nil {
			return ExecutableJAR{}, false, err
		}

		md, ok := NewMetadata(m)
		if !ok {
			return ExecutableJAR{}, false, nil
		}

		class = md.MainClass
	}

	return ExecutableJAR{
		build.Application,
		class,
		build.Layers.Layer(Dependency),
		build.Layers,
		build.Logger,
	}, true, nil
}

type marker string

func (m marker) Identity() (string, string) {
	return "Executable JAR Classpath", ""
}
