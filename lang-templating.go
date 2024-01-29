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
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/go-corelibs/fmtstr"
	"github.com/go-corelibs/slices"
	clstrings "github.com/go-corelibs/strings"
)

func parseTmplSubStatements(statement string) (list []string) {

	var s scanner.Scanner
	s.Init(strings.NewReader(statement))
	s.Error = func(_ *scanner.Scanner, _ string) {}
	s.Filename = "input.tmpl"
	s.Mode ^= scanner.SkipComments
	s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<' '

	var isOpen bool
	var current string

	var stack []string

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		token := s.TokenText()

		stackSize := len(stack)
		if token == "(" {
			if isOpen {
				stack = append(stack, strings.TrimSpace(current))
				current = ""
			} else {
				isOpen = true
			}
			continue
		}
		if token == ")" {
			list = append(list, strings.TrimSpace(current))
			if stackSize > 0 {
				current = stack[stackSize-1]
				stack = stack[:stackSize-1]
			} else {
				current = ""
				isOpen = false
			}
			continue
		}

		if isOpen {
			current += token
		}
	}

	if len(stack) > 0 {
		list = append(list, stack...)
	}

	return
}

func parseTmplStatements(input string) (list []string) {

	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Error = func(_ *scanner.Scanner, _ string) {}
	s.Filename = "input.tmpl"
	s.Mode ^= scanner.SkipComments
	s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<' '

	var foundOpen, isOpen, foundClose bool
	var current string

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		token := s.TokenText()

		if token == "{" {
			if foundOpen {
				// found second opening curly-brace
				isOpen = true
				continue
			}
			// found first opening curly-brace
			foundOpen = true
			continue
		} else if foundOpen {
			foundOpen = false
		}

		if token == "}" {
			if foundClose {
				// found second closing curly-brace
				list = append(list, strings.TrimSpace(strings.Trim(current, "-")))
				if extras := parseTmplSubStatements(current); len(extras) > 0 {
					list = append(list, extras...)
				}
				current = ""
				foundOpen, isOpen, foundClose = false, false, false
				continue
			}
			// found first closing curly-brace
			foundClose = true
			continue
		} else if foundClose {
			foundClose = false
		}

		if isOpen {
			current += token
		}
	}

	return
}

func ParseMessagePlaceholders(key string, argv ...string) (replaced, labelled string, placeholders Placeholders) {
	var subs fmtstr.Variables
	replaced, labelled, subs, _ = fmtstr.Decompose(key, argv...)
	for _, sub := range subs {
		placeholders = append(placeholders, &Placeholder{
			ID:             sub.Label,
			String:         sub.String(),
			Type:           sub.Type,
			UnderlyingType: sub.Type,
			ArgNum:         sub.Pos,
			Expr:           "-",
		})
	}
	return
}

func MakeMessageFromKey(key, comment string, argv ...string) (m *Message) {
	replaced, labelled, placeholders := ParseMessagePlaceholders(key, argv...)
	m = &Message{
		ID:                labelled,
		Key:               key,
		Message:           replaced,
		Translation:       &Translation{String: replaced},
		TranslatorComment: CoalesceTranslatorComment(comment),
		Placeholders:      placeholders,
		Fuzzy:             true,
	}
	return
}

type parseTranlatorCommentState struct {
	skipOnce    bool
	skipMany    int
	comment     string
	commentOpen bool
	source      string
	sourceOpen  bool
	comments    []string
	sources     []string
}

func (state *parseTranlatorCommentState) processComment(idx, last int, this, next rune) (handled bool) {
	if handled = state.commentOpen; handled {
		if state.skipOnce = this == '*' && idx < last && next == '/'; state.skipOnce {
			state.comments = append(state.comments, strings.TrimSpace(state.comment))
			state.commentOpen = false
			state.comment = ""
			return
		}
		state.comment += string(this)
	}
	return
}

func (state *parseTranlatorCommentState) processSource(idx, last int, this, next rune) (handled bool) {
	if handled = state.sourceOpen; handled {
		if state.skipOnce = this == ',' && idx < last && next == ' '; state.skipOnce {
			state.sources = append(state.sources, strings.TrimSpace(state.source))
			state.source = ""
			return
		} else if this == ']' {
			state.sources = append(state.sources, strings.TrimSpace(state.source))
			state.sourceOpen = false
			state.source = ""
			return
		}
		state.source += string(this)
	}
	return
}

func (state *parseTranlatorCommentState) process(idx, last int, this, next rune, input string) {
	if state.processComment(idx, last, this, next) {
	} else if state.processSource(idx, last, this, next) {
	} else if state.skipOnce = this == '/' && idx < last && next == '*'; state.skipOnce {
		state.commentOpen = true
	} else if this == '[' && idx+6 < last && input[idx:idx+7] == "[from: " {
		state.sourceOpen = true
		state.skipMany = 6
	}
	return
}

func ParseTranslatorComment(input string) (comments, sources []string) {
	// /* screen reader only *//* screen reader only */\n[from: layouts/partials/footer.tmpl, path/file.name...]

	state := &parseTranlatorCommentState{}

	last := len(input) - 1
	for idx, this := range input {

		if state.skipOnce {
			state.skipOnce = false
			continue
		} else if state.skipMany > 0 {
			state.skipMany -= 1
			continue
		} else if this == '\n' {
			continue
		}

		var next rune
		if idx < last {
			next = rune(input[idx+1])
		}

		state.process(idx, last, this, next, input)

	}

	return
}

func CoalesceTranslatorComment(input string) (coalesced string) {
	if input == "" {
		return
	}
	comments, sources := ParseTranslatorComment(input)
	var comment, source string
	if len(comments) > 0 {
		comments = slices.Unique(comments)
		comment = strings.Join(comments, "; ")
	}
	if len(sources) > 0 {
		lookup := slices.DuplicateCounts(sources)
		sources = slices.Unique(sources)
		for idx, src := range sources {
			if count, present := lookup[src]; present {
				sources[idx] += "=" + strconv.Itoa(count)
			}
		}
		source = strings.Join(sources, ", ")
	}
	if comment != "" {
		coalesced += "/* " + comment + " */"
	}
	if source != "" {
		if coalesced != "" {
			coalesced += "\n"
		}
		coalesced += "[from: " + source + "]"
	}
	return
}

type parseMessageState struct {
	format  string
	argv    []string
	comment string
}

func (state *parseMessageState) process(s scanner.Scanner) (ok bool) {
	ok = true
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		token := s.TokenText()
		if state.format == "" {
			if size := len(token); size > 2 {
				if clstrings.IsQuoted(token) {
					state.format = clstrings.TrimQuotes(token)
				} else {
					// variable translation
					ok = false
					return
				}
			}
		} else if strings.HasPrefix(token, "/*") {
			if state.comment != "" {
				state.comment += "\n"
			}
			state.comment += token
		} else if clstrings.IsQuoted(token) {
			// support quoted string arguments
			state.argv = append(state.argv, clstrings.TrimQuotes(token))
		} else {
			if argc := len(state.argv); argc > 0 {
				switch state.argv[argc-1] {
				case "$", ".":
					state.argv[argc-1] += token
				default:
					state.argv = append(state.argv, token)
				}
			} else {
				state.argv = append(state.argv, token)
			}
		}
	}
	return
}

func pruneTemplateMessages(input string) (pruned []string) {
	for _, item := range parseTmplStatements(input) {
		if strings.HasPrefix(item, "_ ") {
			item = item[2:]
			if pIdx := strings.Index(item, "|"); pIdx > -1 {
				item = strings.TrimSpace(item[:pIdx])
			}
			pruned = append(pruned, item)
		}
	}
	return
}

func parseMessageStateList(pruned []string) (list []*parseMessageState) {
	for _, item := range pruned {
		var s scanner.Scanner
		s.Init(strings.NewReader(item))
		s.Error = func(_ *scanner.Scanner, _ string) {}
		s.Filename = "input.tmpl"
		s.Mode ^= scanner.SkipComments
		//s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<' '
		s.IsIdentRune = func(ch rune, i int) bool {
			if i == 0 {
				if ch == '$' || ch == '.' {
					return true
				}
				// all template identifiers start with $ or .
				return false
			}
			return ch == '.' || ch == '_' || unicode.IsLetter(ch) || (unicode.IsDigit(ch) && i > 1)
		}

		state := &parseMessageState{}
		if state.process(s) {
			list = append(list, state)
		}
	}

	return
}

func parseUniqueOrderMessages(list []*parseMessageState) (unique map[string][]*parseMessageState, order []string) {
	unique = make(map[string][]*parseMessageState)
	for _, item := range list {
		if _, present := unique[item.format]; !present {
			order = append(order, item.format)
		}
		unique[item.format] = append(unique[item.format], item)
	}
	return
}

func ParseTemplateMessages(input string) (msgs []*Message, err error) {

	// TODO: implement support for inline template statements
	//       example: _ "the message: %[1]s" (printf "not working yet")
	//       error: placeholders are "Printf" and "NotWorkingYet", even though there is only one replacement verb

	pruned := pruneTemplateMessages(input)
	list := parseMessageStateList(pruned)
	unique, order := parseUniqueOrderMessages(list)

	for _, key := range order {
		items, _ := unique[key]
		item := items[0]
		comment := item.comment
		if count := len(items); count > 1 {
			for idx, itm := range items {
				if idx == 0 {
					continue
				}
				if itm.comment != "" {
					var dupe bool
					for jdx, other := range items {
						if dupe = idx != jdx && other.comment == itm.comment; dupe {
							break
						}
					}
					if !dupe {
						if comment != "" {
							comment += "\n"
						}
						comment += itm.comment
					}
				}
			}
		}
		msg := MakeMessageFromKey(item.format, comment, item.argv...)
		if argc := len(item.argv); argc > 0 {
			for idx, placeholder := range msg.Placeholders {
				index := placeholder.ArgNum - 1
				if argc > index {
					if msg.Placeholders[idx].Expr == "-" {
						msg.Placeholders[idx].Expr = item.argv[index]
					} else {
						msg.Placeholders[idx].Expr += ", " + item.argv[index]
					}
				}
			}
		}
		msgs = append(msgs, msg)
	}

	return
}
