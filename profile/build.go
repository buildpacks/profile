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
	_ "embed"
	"os"

	"github.com/buildpacks/libcnb"
)

//go:embed execd_wrapper.sh
var execDScript string

func Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	// NOTE: the logger is not passed into this function, that will likely be a change in libcnbv2
	result := libcnb.NewBuildResult()

	layer, err := context.Layers.Layer(profileName)

	if err != nil {
		return result, err
	}

	execPath := layer.Exec.FilePath(execDScriptName)

	err = os.MkdirAll(layer.Exec.Path, os.ModePerm)
	if err != nil {
		return result, err
	}

	f, err := os.Create(execPath)
	if err != nil {
		return result, err
	}
	defer f.Close()

	_, err = f.WriteString(execDScript)
	if err != nil {
		return result, err
	}

	err = os.Chmod(execPath, 0755)
	if err != nil {
		return result, err
	}

	layer.Launch = true
	layer.Cache = true
	result.Layers = append(result.Layers, layer)

	return result, nil
}
