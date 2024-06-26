/*
 * Copyright 2018-2024 the original author or authors.
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
	_, shErr := exec.LookPath("bash")

	profilePath := filepath.Join(context.ApplicationPath, scriptName)
	if _, err := os.Stat(profilePath); shErr == nil && !os.IsNotExist(err) {
		context.Logger.Debug("PASSED: a .profile application exists and bash is available")
		return libcnb.DetectResult{Pass: true, Plans: []libcnb.BuildPlan{}}, nil
	}

	context.Logger.Debug("SKIPPED: a .profile application does not exist or bash is not available")
	return libcnb.DetectResult{Pass: false}, nil
}
