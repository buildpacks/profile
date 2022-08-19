/*
 * Copyright 2018-2022 the original author or authors.
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

package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/buildpacks/profile/profile"

	"github.com/buildpacks/libcnb"
)

type BuildTest struct {
	test    *testing.T
	context libcnb.BuildContext
	expect  ExpectFunc
}

func NewBuildTestBuilder(t *testing.T) BuildTest {
	return BuildTest{
		test:   t,
		expect: NewGomegaWithT(t).Expect,
	}
}

func (b BuildTest) SetupWorkspace() BuildTest {
	var err error

	b.context.Buildpack.Path = "../"
	b.context.ApplicationPath = b.test.TempDir()
	b.context.Layers.Path = b.test.TempDir()

	err = os.MkdirAll(b.context.Layers.Path, os.ModePerm)
	b.expect(err).ToNot(HaveOccurred())

	profilePath := filepath.Join(b.context.ApplicationPath, ".profile")
	f, err := os.Create(profilePath)
	b.expect(err).NotTo(HaveOccurred())
	defer f.Close()

	_, err = f.WriteString(
		`
echo "Hello world"
HELLO="world"
export HELLO
		`,
	)
	b.expect(err).NotTo(HaveOccurred())

	return b
}

func (b BuildTest) Build() (libcnb.BuildContext, ExpectFunc) {
	return b.context, b.expect
}

func TestBuildExecutes(t *testing.T) {
	ctx, Expect := NewBuildTestBuilder(t).SetupWorkspace().Build()

	Expect(profile.Build(ctx)).To(Equal(libcnb.BuildResult{
		Layers: []libcnb.Layer{{
			LayerTypes: libcnb.LayerTypes{
				Build:  false,
				Launch: true,
			},
			BuildEnvironment:  libcnb.Environment{},
			LaunchEnvironment: libcnb.Environment{},
			SharedEnvironment: libcnb.Environment{},
			Name:              "profile",
			Path:              filepath.Join(ctx.Layers.Path, "profile"),
			Profile:           libcnb.Profile{},
			Exec: libcnb.Exec{
				Path: filepath.Join(ctx.Layers.Path, "profile", "exec.d"),
			},
		}},
		PersistentMetadata: map[string]interface{}{},
	}))

	Expect(filepath.Join(ctx.Layers.Path, "profile", "exec.d", profile.ExecDScriptName)).To(BeAnExistingFile())
}
