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
	"fmt"

	"github.com/go-corelibs/x-text/language"
)

func ParseTag(input interface{}) (tag language.Tag, err error) {
	var text string
	switch t := input.(type) {
	case language.Tag:
		tag = t
		return
	case string:
		text = t
	default:
		text = fmt.Sprintf("%v", t)
	}
	tag, err = language.Parse(text)
	return
}
