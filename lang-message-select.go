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
	"strconv"
	"strings"
)

type Select struct {
	Arg     string                `json:"arg"`
	Feature string                `json:"feature"`
	Cases   map[string]SelectCase `json:"cases"`
}

type SelectCase struct {
	Msg string `json:"msg"`
}

type SelectMessage struct {
	ID                string       `json:"id"`
	Key               string       `json:"key"`
	Message           string       `json:"message"`
	Translation       *Translation `json:"translation"`
	TranslatorComment string       `json:"translatorComment,omitempty"`
	Placeholders      Placeholders `json:"placeholders,omitempty"`
	Fuzzy             bool         `json:"fuzzy,omitempty"`
}

func (s *SelectMessage) Make() (m Message) {
	m = Message{
		ID:                s.ID,
		Key:               s.Key,
		Message:           s.Message,
		Translation:       s.Translation,
		TranslatorComment: s.TranslatorComment,
		Placeholders:      s.Placeholders[:],
		Fuzzy:             s.Fuzzy,
	}
	return
}

func ParsePluralCaseKey(key string) (real string) {
	real = strings.ToLower(strings.TrimSpace(key))
	switch real {
	case "zero", "one", "two", "few", "many", "other":
		return
	}
	if len(real) <= 1 {
		real = ""
		return
	}
	switch real[0] {
	case '=':
		if v, ee := strconv.Atoi(real[1:]); ee != nil {
			real = ""
		} else {
			real = fmt.Sprintf("=%d", v)
		}
	case '<':
		if v, ee := strconv.Atoi(real[1:]); ee != nil {
			real = ""
		} else {
			real = fmt.Sprintf("<%d", v)
		}
	default:
		real = ""
	}
	return
}
