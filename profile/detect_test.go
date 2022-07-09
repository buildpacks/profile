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
type AfterFunc func()

type DetectTest struct {
	context libcnb.DetectContext
	expect  ExpectFunc
}

func (d DetectTest) SetupGomega(t *testing.T) DetectTest {
	d.expect = NewGomegaWithT(t).Expect
	return d
}

func (d DetectTest) SetupWorkspace() DetectTest {
	var err error
	d.context.ApplicationPath, err = os.MkdirTemp("", "profile")
	d.expect(err).NotTo(HaveOccurred())
	return d
}

func (d DetectTest) RemoveWorkspace() DetectTest {
	d.expect(os.RemoveAll(d.context.ApplicationPath)).To(Succeed())
	return d
}

func (d DetectTest) Build() (libcnb.DetectContext, ExpectFunc, AfterFunc) {
	return d.context, d.expect, func() { d.RemoveWorkspace() }
}

func TestDetectFailsWithoutProfileScript(t *testing.T) {
	ctx, Expect, After := DetectTest{}.SetupGomega(t).SetupWorkspace().Build()
	defer After()

	Expect(profile.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
}

func TestDetectPassesWithProfileScript(t *testing.T) {
	ctx, Expect, After := DetectTest{}.SetupGomega(t).SetupWorkspace().Build()
	defer After()

	Expect(os.WriteFile(filepath.Join(ctx.ApplicationPath, ".profile"), []byte(`echo "Hello World!"`), 0600))

	Expect(profile.Detect(ctx)).To(Equal(libcnb.DetectResult{
		Pass: false, // TODO: this should be true, after implemented
	}))
}
