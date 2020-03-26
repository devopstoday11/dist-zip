/*
 * Copyright 2018-2020 the original author or authors.
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

package distzip

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var DefaultPattern = filepath.Join("*", "bin", `*`)

type ScriptResolver struct {
	ApplicationPath string
}

func (s *ScriptResolver) Resolve() (string, bool, error) {
	var (
		candidates []string
		err        error
		ok         bool
		pattern    string
	)

	if pattern, ok = os.LookupEnv("BP_APPLICATION_SCRIPT"); ok {
		file := filepath.Join(s.ApplicationPath, pattern)
		candidates, err = filepath.Glob(file)
		if err != nil {
			return "", false, fmt.Errorf("unable to find files with %s\n%w", pattern, err)
		}
	} else {
		pattern = DefaultPattern
		file := filepath.Join(s.ApplicationPath, pattern)
		candidates, err = filepath.Glob(file)
		if err != nil {
			return "", false, fmt.Errorf("unable to find files with %s\n%w", pattern, err)
		}

		i := 0
		for i < len(candidates) {
			if strings.HasSuffix(candidates[i], ".bat") {
				candidates = append(candidates[:i], candidates[i+1:]...)
			}
			i++
		}
	}

	switch len(candidates) {
	case 0:
		return "", false, nil
	case 1:
		return candidates[0], true, nil
	default:
		sort.Strings(candidates)
		return "", false, fmt.Errorf("unable to find application script in %s, candidates: %s", pattern, candidates)
	}
}
