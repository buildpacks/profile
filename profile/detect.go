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

package profile

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/buildpacks/libcnb/v2"
)

func Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	// NOTE: the logger is not passed into this function, that will likely be a change in libcnbv2

	_, shErr := exec.LookPath("bash")

	profilePath := filepath.Join(context.ApplicationPath, scriptName)
	if _, err := os.Stat(profilePath); shErr == nil && !os.IsNotExist(err) {
		return libcnb.DetectResult{Pass: true, Plans: []libcnb.BuildPlan{}}, nil
	}

	return libcnb.DetectResult{Pass: false}, nil
}
