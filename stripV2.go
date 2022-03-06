package main

import (
	"html"
	"strings"
)

func StripMDV2(s string) string {
	return defaultConverterV2.StripMDV2(s)
}

func (cv *ConverterV2) StripMDV2(s string) string {
	text, _ := cv.MD2HTMLButtons(s)
	return cv.stripHTML([]rune(text))
}

func StripHTMLV2(s string) string {
	return defaultConverterV2.stripHTML([]rune(s))
}

func (cv *ConverterV2) StripHTMLV2(s string) string {
	return cv.stripHTML([]rune(s))
}

func (cv *ConverterV2) stripHTML(in []rune) string {
	out := strings.Builder{}
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case '<':
			tagClose := getTagClose(in[i+1:])
			if tagClose < 0 {
				// gone weird; just skip.
				continue
			}

			i += tagClose + 1 // skip to closing tag.
			continue

		default:
			out.WriteRune(in[i])
		}
	}
	return html.UnescapeString(out.String())
}
