// Copyright 2015 Google Inc.
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

package yang

var knownWords = map[string]string{
	"Ietf": "IETF",
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore or dash followed by a lower case letter,
// drop the underscore or dash and convert the letter to upper case.  There is a
// remote possibility of this rewrite causing a name collision, but it's so
// remote we're prepared to pretend it's nonexistent - since the C++ generator
// lowercases names, it's extremely unlikely to have two fields with different
// capitalizations.  In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	fix := func(c byte) byte {
		if c == '/' || c == '-' || c == ':' {
			return '_'
		}
		return c
	}
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if fix(s[0]) == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := fix(s[i])
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		start := len(t)
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
		if kn := knownWords[string(t[start:])]; kn != "" {
			t = append(t[:start], []byte(kn)...)
		}
	}
	return string(t)
}
