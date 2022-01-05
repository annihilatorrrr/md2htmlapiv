package main

import (
	"html"
	"sort"
	"strings"
)

var defaultConverterV2 = ConverterV2{
	BtnPrefix:      btnPrefix,
	SameLineSuffix: sameLineSuffix,
}

type ButtonV2 struct {
	Name     string
	Content  string
	SameLine bool
}

type ConverterV2 struct {
	BtnPrefix      string
	SameLineSuffix string
}

func NewV2() *ConverterV2 {
	return &ConverterV2{
		BtnPrefix:      btnPrefix,
		SameLineSuffix: sameLineSuffix,
	}
}

func MD2HTMLV2(in string) string {
	return defaultConverterV2.MD2HTML(in)
}

func MD2HTMLButtonsV2(in string) (string, []ButtonV2) {
	return defaultConverterV2.md2html([]rune(html.EscapeString(in)), true)
}

var chars = map[string]string{
	"`":   "code",
	"```": "pre",
	"_":   "i",
	"*":   "b",
	"~":   "s",
	"__":  "u",
	"|":   "", // this is a placeholder for || to work
	"||":  "span class=\"tg-spoiler\"",
	"[":   "", // for links
	"]":   "", // for links
	"(":   "", // for links
	")":   "", // for links
	"\\":  "", // for escapes
}

var AllMarkdownV2Chars = func() []rune {
	var outString []string
	for k := range chars {
		outString = append(outString, k)
	}
	sort.Strings(outString)
	var out []rune
	for _, x := range outString {
		out = append(out, []rune(x)[0])
	}
	return out
}()

func (cv ConverterV2) MD2HTML(in string) string {
	text, _ := cv.md2html([]rune(html.EscapeString(in)), false)
	return text
}

func (cv ConverterV2) MD2HTMLButtons(in string) (string, []ButtonV2) {
	return cv.md2html([]rune(html.EscapeString(in)), true)
}

// TODO: add support for a map-like check of which items cannot be included.
//  Eg: `code` cannot be italic/bold/underline/strikethrough
//  however... this is currently implemented by server side by telegram, so not my problem :runs:
func (cv ConverterV2) md2html(in []rune, b bool) (string, []ButtonV2) {
	out := strings.Builder{}

	for i := 0; i < len(in); i++ {
		c := in[i]
		if _, ok := chars[string(c)]; !ok {
			out.WriteRune(c)
			continue
		}

		if !validStart(i, in) {
			if c == '\\' && i+1 < len(in) {
				if _, ok := chars[string(in[i+1])]; ok {
					out.WriteRune(in[i+1])
					i++
					continue
				}
			}
			out.WriteRune(c)
			continue
		}

		switch c {
		case '`', '*', '~', '_', '|': // '||', '__', and '```' are included here too
			item := string(c)
			if c == '|' { // support ||
				// if single |, ignore. We only care about double ||
				if i+1 >= len(in) || in[i+1] != '|' {
					out.WriteRune(c)
					continue
				}

				item = "||"
				i++
			} else if c == '_' && i+1 < len(in) && in[i+1] == '_' { // support __
				item = "__"
				i++
			} else if c == '`' && i+2 < len(in) && in[i+1] == '`' && in[i+2] == '`' { // support ```
				item = "```"
				i += 2
			}

			if i+1 >= len(in) {
				out.WriteString(item)
				continue
			}

			idx := getValidEnd(in[i+1:], item)
			if idx < 0 {
				// not found; write and move on.
				out.WriteString(item)
				continue
			}

			nStart, nEnd := i+1, i+idx+1

			var nestedT string
			var nestedB []ButtonV2
			if c == '`' {
				// ` and ``` dont support nested items, so don't parse children.
				nestedT = string(in[nStart:nEnd])
			} else {
				// internal wont have any interesting item closings
				nestedT, nestedB = cv.md2html(in[nStart:nEnd], b)
			}
			// nestedT, nestedB := cv.md2html(in[nStart:nEnd], b)
			followT, followB := cv.md2html(in[nEnd+len(item):], b)

			return out.String() + "<" + chars[item] + ">" + nestedT + "</" + closeSpans(chars[item]) + ">" + followT, append(nestedB, followB...)

		case '[':
			// find ]( and then )
			linkText, linkURL := findLinkSections(in[i:])
			if linkText < 0 || linkURL < 0 {
				out.WriteRune(c)
				continue
			}

			content := string(in[i+linkText+2 : i+linkURL])
			text := in[i+1 : i+linkText]
			end := i + linkURL + 1
			followT, followB := cv.md2html(in[end:], b)

			if b && strings.HasPrefix(content, cv.BtnPrefix) {
				content = strings.TrimLeft(content[len(cv.BtnPrefix):], "/")
				sameline := false
				if strings.HasSuffix(content, cv.SameLineSuffix) {
					sameline = true
					content = content[:len(content)-len(cv.SameLineSuffix)]
				}
				return out.String() + followT, append([]ButtonV2{{
					Name:     html.UnescapeString(string(text)),
					Content:  content,
					SameLine: sameline,
				}}, followB...)
			}
			nestedT, nestedB := cv.md2html(text, true)
			return out.String() + `<a href="` + content + `">` + nestedT + "</a>" + followT, append(nestedB, followB...)

		case ']', '(', ')':
			out.WriteRune(c)

		case '\\':
			if i+1 < len(in) {
				if _, ok := chars[string(in[i+1])]; ok {
					out.WriteRune(in[i+1])
					i++
					continue
				}
			}
			out.WriteRune(c)
		}
	}

	return out.String(), nil
}

func EscapeMarkdownV2(r []rune) string {
	out := strings.Builder{}
	for i, x := range r {
		if contains(x, AllMarkdownV2Chars) {
			if i == 0 || i == len(r)-1 || validEnd(i, r) || validStart(i, r) {
				out.WriteRune('\\')
			}
		}
		out.WriteRune(x)
	}
	return out.String()
}

// closeSpans gets the correct closing tags for spans.
// eg: closeSpans("span class=\"tg-spoiler\"") should return just "span"
// 	   closeSpans("pre") -> returns "pre"
func closeSpans(s string) string {
	if !strings.HasPrefix(s, "span") {
		return s
	}

	return "span"
}
