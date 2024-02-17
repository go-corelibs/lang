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
	"strings"

	clStrings "github.com/go-corelibs/strings"
)

func startsWithUnderscore(input string) (ok bool) {
	if v := strings.TrimSpace(input); len(v) >= 2 {
		ok = v[:2] == "_ "
	}
	return
}

func endsWithCloseComment(input string) (ok bool) {
	v := strings.TrimSpace(input)
	if size := len(v); size >= 2 {
		ok = v[size-2:] == "*/"
	}
	return
}

func PruneCommandTranslatorComments(raw string) (clean string) {
	remainder := raw
	var stack []string
	for {
		if before, middle, after, found := clStrings.ScanCarve(remainder, "{{", "}}"); found {
			stack = append(stack, before)
			var sDash, eDash string
			if strings.HasPrefix(middle, "-") {
				sDash = "-"
				middle = middle[1:]
			}
			if strings.HasSuffix(middle, "-") {
				eDash = "-"
				middle = middle[:len(middle)-1]
			}
			if startsWithUnderscore(middle) && endsWithCloseComment(middle) {
				if idx := strings.Index(middle, "/*"); idx >= 0 {
					stack = append(stack, "{{"+sDash+middle[:idx]+eDash+"}}")
				} else {
					stack = append(stack, "{{"+sDash+middle+eDash+"}}")
				}
			} else {
				stack = append(stack, "{{"+sDash+middle+eDash+"}}")
			}
			remainder = after
			continue
		}
		stack = append(stack, remainder)
		break
	}
	clean = strings.Join(stack, "")
	return
}

func PruneInlineTranslatorComments(raw string) (clean string) {
	remainder := raw
	var stack []string
	for {
		if before, middle, after, found := clStrings.ScanCarve(remainder, "{{", "}}"); found {
			stack = append(stack, before+"{{")
			var mStack []string
			mRemainder := middle[:]
			for {
				if b, m, a, f := clStrings.ScanCarve(mRemainder, "(", ")"); f {
					mStack = append(mStack, b)
					if startsWithUnderscore(m) && endsWithCloseComment(m) {
						if idx := strings.Index(m, "/*"); idx >= 0 {
							mStack = append(mStack, "("+m[:idx]+")")
						} else {
							mStack = append(mStack, "("+m+")")
						}
					} else {
						mStack = append(mStack, "("+m+")")
					}
					mRemainder = a
					continue
				}
				mStack = append(mStack, mRemainder)
				break
			}
			stack = append(stack, strings.Join(mStack, ""))
			stack = append(stack, "}}")
			remainder = after
			continue
		}
		break
	}
	clean = strings.Join(stack, "") + remainder
	return
}

func PruneTranslatorComments(input string) (clean string) {
	return PruneCommandTranslatorComments(
		PruneInlineTranslatorComments(input),
	)
}
