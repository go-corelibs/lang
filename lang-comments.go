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
	"regexp"
)

var (
	//rxGoComments                 = regexp.MustCompile(`(?ms)\{\{\s*/\*.+?\*/\s*\}\}`)
	rxTranslatorInlineComments   = regexp.MustCompile(`(?ms)\((\s*_\s+.+?\s*)/\*.+?\*/\s*\)`)
	rxTranslatorPipelineComments = regexp.MustCompile(`(?ms)\{\{(-??\s*_\s+.+?\s*)/\*.+?\*/(\s*-??)}}`)
)

func PruneTranslatorComments(raw string) (clean string) {
	clean = rxTranslatorInlineComments.ReplaceAllString(raw, `(${1})`)
	clean = rxTranslatorPipelineComments.ReplaceAllString(clean, `{{${1}${2}}}`)
	//clean = rxGoComments.ReplaceAllString(clean, "")
	return
}
