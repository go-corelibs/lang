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
	"github.com/go-corelibs/slices"
)

type Placeholder struct {
	ID             string `json:"id"`
	String         string `json:"string"`
	Type           string `json:"type"`
	UnderlyingType string `json:"underlyingType"`
	ArgNum         int    `json:"argNum"`
	Expr           string `json:"expr"`
}

type Placeholders []*Placeholder

func (p Placeholders) Numeric() (found Placeholders) {
	for _, placeholder := range p {
		if slices.Present(placeholder.UnderlyingType, "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64") {
			found = append(found, placeholder)
		}
	}
	return
}

func (p Placeholders) Copy() (copied Placeholders) {
	for _, ph := range p {
		copied = append(copied, &(*ph))
	}
	return
}
