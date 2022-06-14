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
	return b
}

func (b BuildTest) RemoveWorkspace() BuildTest {
	b.expect(os.RemoveAll(b.context.ApplicationPath)).To(Succeed())
	return b
}

func (b BuildTest) Build() (libcnb.BuildContext, ExpectFunc, AfterFunc) {
	return b.context, b.expect, func() { b.RemoveWorkspace() }
}

func TestBuildDoesNothingWithoutPlanEntry(t *testing.T) {
	ctx, Expect, After := BuildTest{}.SetupGomega(t).SetupWorkspace().Build()
	defer After()

	Expect(profile.Build(ctx)).To(Equal(libcnb.NewBuildResult()))
}

func TestBuildExecutesWithPlanEntry(t *testing.T) {
	ctx, Expect, After := BuildTest{}.SetupGomega(t).SetupWorkspace().Build()
	defer After()

	// TODO: test setup for a working build

	// TODO: test validation for a working build
	Expect(profile.Build(ctx)).To(Equal(libcnb.NewBuildResult()))
}
