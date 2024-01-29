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
	"encoding/json"
)

type BaseMessage struct {
	ID                string       `json:"id"`
	Key               string       `json:"key"`
	Message           string       `json:"message"`
	Translation       interface{}  `json:"translation"`
	TranslatorComment string       `json:"translatorComment,omitempty"`
	Placeholders      Placeholders `json:"placeholders,omitempty"`
	Fuzzy             bool         `json:"fuzzy,omitempty"`
}

type Translation struct {
	String string  `json:"string,omitempty"`
	Select *Select `json:"select,omitempty"`
}

type Message struct {
	ID                string       `json:"id"`
	Key               string       `json:"key"`
	Message           string       `json:"message"`
	Translation       *Translation `json:"translation"`
	TranslatorComment string       `json:"translatorComment,omitempty"`
	Placeholders      Placeholders `json:"placeholders,omitempty"`
	Fuzzy             bool         `json:"fuzzy,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) (err error) {
	structMessage := &SelectMessage{}
	if e := json.Unmarshal(data, &structMessage); e == nil {
		*m = structMessage.Make()
		return
	}
	stringMessage := &StringMessage{}
	if e := json.Unmarshal(data, &stringMessage); e == nil {
		*m = stringMessage.Make()
		return
	} else {
		err = e
	}
	return
}

func (m *Message) MarshalJSON() (data []byte, err error) {
	if m.Translation.Select != nil {
		data, err = json.MarshalIndent(BaseMessage{
			ID:      m.ID,
			Key:     m.Key,
			Message: m.Message,
			Translation: map[string]interface{}{
				"select": m.Translation.Select,
			},
			TranslatorComment: m.TranslatorComment,
			Placeholders:      m.Placeholders[:],
			Fuzzy:             m.Fuzzy,
		}, "", "    ")
		return
	}
	data, err = json.MarshalIndent(BaseMessage{
		ID:                m.ID,
		Key:               m.Key,
		Message:           m.Message,
		Translation:       m.Translation.String,
		TranslatorComment: m.TranslatorComment,
		Placeholders:      m.Placeholders[:],
		Fuzzy:             m.Fuzzy,
	}, "", "    ")
	return
}

func (m *Message) Copy() (copied *Message) {
	copied = &Message{
		ID:                m.ID,
		Key:               m.Key,
		Message:           m.Message,
		Translation:       &Translation{},
		TranslatorComment: m.TranslatorComment,
		Placeholders:      m.Placeholders.Copy(),
		Fuzzy:             m.Fuzzy,
	}
	if m.Translation.Select == nil {
		copied.Translation.String = m.Translation.String
		return
	}
	copied.Translation.Select = &Select{
		Arg:     m.Translation.Select.Arg,
		Feature: m.Translation.Select.Feature,
		Cases:   make(map[string]SelectCase),
	}
	for k, sc := range m.Translation.Select.Cases {
		copied.Translation.Select.Cases[k] = SelectCase{
			Msg: sc.Msg,
		}
	}
	return
}
