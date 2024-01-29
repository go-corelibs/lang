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

type StringMessage struct {
	ID                string       `json:"id"`
	Key               string       `json:"key"`
	Message           string       `json:"message"`
	Translation       string       `json:"translation"`
	TranslatorComment string       `json:"translatorComment,omitempty"`
	Placeholders      Placeholders `json:"placeholders,omitempty"`
	Fuzzy             bool         `json:"fuzzy,omitempty"`
}

func (s *StringMessage) Make() (m Message) {
	m = Message{
		ID:      s.ID,
		Key:     s.Key,
		Message: s.Message,
		Translation: &Translation{
			String: s.Translation,
		},
		TranslatorComment: s.TranslatorComment,
		Placeholders:      s.Placeholders[:],
		Fuzzy:             s.Fuzzy,
	}
	return
}
