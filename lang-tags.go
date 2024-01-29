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

type Tags []language.Tag

func (list Tags) Has(tag language.Tag) (present bool) {
	for _, t := range list {
		if present = t == tag; present {
			return
		}
	}
	return
}

func (list Tags) Strings() (locales []string) {
	for _, t := range list {
		locales = append(locales, t.String())
	}
	return
}

func (list Tags) StringsWithDefault(tag language.Tag) (locales []string) {
	if list.Has(tag) {
		locales = append(locales, tag.String())
	}
	for _, t := range list {
		if t == tag {
			continue
		}
		locales = append(locales, t.String())
	}
	return
}
