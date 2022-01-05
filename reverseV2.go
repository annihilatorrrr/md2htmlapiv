package md2htmlapi

import (
	"fmt"
	"html"
	"strings"
)

func ReverseV2(in string, bs []ButtonV2) (string, error) {
	return defaultConverterV2.Reverse(in, bs)
}

func (cv *ConverterV2) Reverse(in string, bs []ButtonV2) (string, error) {
	return cv.reverse([]rune(in), bs)
}

func (cv *ConverterV2) reverse(in []rune, buttons []ButtonV2) (string, error) {
	prev := 0
	out := strings.Builder{}
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case '<':
			c := getTagClose(in[i+1:])
			if c < 0 {
				// "no close tag"
				return "", fmt.Errorf("no closing '>' for opening bracket at %d", i)
			}
			closeTag := i + c + 1
			tagContent := string(in[i+1 : closeTag])
			tagFields := strings.Fields(tagContent)
			if len(tagFields) < 1 {
				return "", fmt.Errorf("no tag name for HTML tag started at %d", i)
			}
			tag := tagFields[0]

			co, cc := getClosingTag(in[closeTag+1:], tag)
			if co < 0 || cc < 0 {
				// "no closing open"
				return "", fmt.Errorf("no closing tag for HTML tag %q started at %d", tag, i)
			}
			closingOpen, closingClose := closeTag+1+co, closeTag+1+cc
			out.WriteString(html.UnescapeString(string(in[prev:i])))

			nested, err := cv.reverse(in[closeTag+1:closingOpen], nil)
			if err != nil {
				return "", err
			}

			switch tag {
			case "b", "strong":
				out.WriteString("*" + nested + "*")
			case "i", "em":
				out.WriteString("_" + nested + "_")
			case "u", "ins":
				out.WriteString("__" + nested + "__")
			case "s", "strike", "del":
				out.WriteString("~" + nested + "~")
			case "code":
				// code and pre don't look at nested values, because they're not parsed
				out.WriteString("`" + html.UnescapeString(string(in[closeTag+1:closingOpen])) + "`")
			case "pre":
				// code and pre don't look at nested values, because they're not parsed
				out.WriteString("```" + html.UnescapeString(string(in[closeTag+1:closingOpen])) + "```")
			case "span":
				// NOTE: All span tags are currently spoiler tags. This may change in the future.
				if len(tagFields) < 2 {
					return "", fmt.Errorf("span tag does not have enough fields %q", tagFields)
				}

				switch spanType := tagFields[1]; spanType {
				case "class=\"tg-spoiler\"":
					out.WriteString("||" + html.UnescapeString(string(in[closeTag+1:closingOpen])) + "||")
				default:
					return "", fmt.Errorf("unknown tag type %q", spanType)
				}
			case "a":
				if link.MatchString(tagContent) {
					matches := link.FindStringSubmatch(tagContent)
					out.WriteString("[" + nested + "](" + matches[1] + ")")
				} else {
					return "", fmt.Errorf("badly formatted anchor tag %q", tagContent)
				}
			default:
				return "", fmt.Errorf("unknown tag %q", tag)
			}

			prev = closingClose + 1
			i = closingClose

		case '\\', '_', '*', '~', '`', '[', ']', '(', ')': // these all need to be escaped to ensure we retain the same message
			out.WriteString(html.UnescapeString(string(in[prev:i])))
			out.WriteRune('\\')
			out.WriteRune(in[i])
			prev = i + 1
		}
	}
	out.WriteString(html.UnescapeString(string(in[prev:])))

	for _, btn := range buttons {
		out.WriteString("\n[" + btn.Name + "](" + cv.BtnPrefix + "//" + html.UnescapeString(btn.Content))
		if btn.SameLine {
			out.WriteString(cv.SameLineSuffix)
		}
		out.WriteString(")")
	}

	return out.String(), nil
}
