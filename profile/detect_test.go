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

type ExpectFunc func(actual interface{}, extra ...interface{}) types.Assertion

type DetectTest struct {
	test    *testing.T
	context libcnb.DetectContext
	expect  ExpectFunc
}

func NewDetectTestBuilder(t *testing.T) DetectTest {
	return DetectTest{
		test:   t,
		expect: NewGomegaWithT(t).Expect,
	}
}

func (d DetectTest) SetupWorkspace() DetectTest {
	d.context.ApplicationPath = d.test.TempDir()
	return d
}

func (d DetectTest) Build() (libcnb.DetectContext, ExpectFunc) {
	return d.context, d.expect
}

func TestDetectFailsWithoutProfileScript(t *testing.T) {
	ctx, Expect := NewDetectTestBuilder(t).SetupWorkspace().Build()

	Expect(profile.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
}

func TestDetectPassesWithProfileScript(t *testing.T) {
	ctx, Expect := NewDetectTestBuilder(t).SetupWorkspace().Build()

	Expect(os.WriteFile(filepath.Join(ctx.ApplicationPath, ".profile"), []byte(`echo "Hello World!"`), 0600))

	Expect(profile.Detect(ctx)).To(Equal(libcnb.DetectResult{
		Pass:  true,
		Plans: []libcnb.BuildPlan{},
	}))
}
