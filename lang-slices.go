// Copyright (c) 2024  The Go-Enjin Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lang

import (
	"github.com/go-enjin/golang-org-x-text/language"
)

func TagInTags(needle language.Tag, haystack ...language.Tag) (found bool) {
	for _, tag := range haystack {
		if found = language.Compare(needle, tag); found {
			return
		}
	}
	return
}

func TagInTagSlices(needle language.Tag, haystacks ...[]language.Tag) (found bool) {
	for _, haystack := range haystacks {
		for _, tag := range haystack {
			if found = language.Compare(needle, tag); found {
				return
			}
		}
	}
	return
}
