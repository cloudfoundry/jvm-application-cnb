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

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/manifest"
)

// Metadata describes the metadata
type Metadata struct {
	jvmapplication.Metadata

	// MainClass is the Main-Class of the executable JAR.
	MainClass string
}

// BuildPlan returns a BuildPlan that contains all the dependencies of jvmapplication.Metadata as well as
// executable-jar.
func (m Metadata) BuildPlan(buildPlan buildplan.BuildPlan) buildplan.BuildPlan {
	bp := m.Metadata.BuildPlan(buildPlan)
	bp[Dependency] = m.executableJar(buildPlan)
	return bp
}

// String makes Metadata satisfy the Stringer interface.
func (m Metadata) String() string {
	return fmt.Sprintf("Metadata{ Metadata: %s, MainClass: %s }", m.Metadata, m.MainClass)
}

func (m Metadata) executableJar(buildPlan buildplan.BuildPlan) buildplan.Dependency {
	d := buildPlan[Dependency]

	if d.Metadata == nil {
		d.Metadata = make(buildplan.Metadata)
	}

	d.Metadata[MainClass] = m.MainClass

	return d
}

// NewMetadata creates a new Metadata returning false if Main-Class is not defined.
func NewMetadata(manifest manifest.Manifest) (Metadata, bool) {
	m, ok := manifest.Get("Main-Class")
	if !ok {
		return Metadata{}, false
	}

	return Metadata{
		jvmapplication.Metadata{},
		m,
	}, true
}
