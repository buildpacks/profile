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
	"github.com/onsi/gomega/types"

	"github.com/buildpacks/profile/profile"

	"github.com/buildpacks/libcnb"
)

type BuildTest struct {
	context libcnb.BuildContext
	expect  func(actual interface{}, extra ...interface{}) types.Assertion
}

func (b BuildTest) SetupGomega(t *testing.T) BuildTest {
	b.expect = NewGomegaWithT(t).Expect
	return b
}

func (b BuildTest) SetupWorkspace() BuildTest {
	var err error
	b.context.ApplicationPath, err = os.MkdirTemp("", "profile")
	b.expect(err).NotTo(HaveOccurred())
	profilePath := filepath.Join(b.context.ApplicationPath, ".profile")

	f, err := os.Create(profilePath)
	b.expect(err).NotTo(HaveOccurred())
	defer f.Close()

	b.expect(err).To(BeNil())

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

func (b BuildTest) RemoveWorkspace() BuildTest {
	b.expect(os.RemoveAll(b.context.ApplicationPath)).To(Succeed())
	return b
}

func (b BuildTest) Build() (libcnb.BuildContext, ExpectFunc, AfterFunc) {
	outputPath, _ := os.MkdirTemp("", "dotprofile_out")
	os.MkdirAll(outputPath, os.ModePerm)
	b.context.Layers.Path = outputPath
	return b.context, b.expect, func() {
		b.RemoveWorkspace()
		os.RemoveAll(outputPath)
	}
}

func TestBuildExecutes(t *testing.T) {
	ctx, Expect, After := BuildTest{}.SetupGomega(t).SetupWorkspace().Build()
	defer After()
	Expect(profile.Build(ctx)).To(Equal(libcnb.BuildResult{
		Layers: []libcnb.Layer{{
			LayerTypes: libcnb.LayerTypes{
				Build:  false,
				Launch: true,
				Cache:  true,
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

	Expect(filepath.Join(ctx.Layers.Path, "profile", "exec.d", "dotprofile.sh")).To(BeAnExistingFile())
}
