package utils

import (
	"homeserver/pkg/parser"
	"strings"
)

var (
	LoadIcon = parser.LoadIcon
)

const (
	IconTagStart          = "<icon "
	MalformedIconTagStart = "<icon>"
	IconTagEnd            = "</icon>"
	ErrorMalformedTag     = "Malformed icon tag in response\n"
	ErrorMissingNameAttr  = "Missing name attribute in icon tag\n"
)

func IconsInjector(b []byte) []byte {
	str := string(b)

	for i := 0; i <= len(str)-7; i++ {
		if str[i:i+6] == MalformedIconTagStart {
			panic(ErrorMissingNameAttr)
		}

		if str[i:i+6] != IconTagStart {
			continue
		}

		cut_str := str[i:]

		end_index := strings.Index(cut_str, IconTagEnd)
		if end_index == -1 {
			panic(ErrorMalformedTag)
		}

		end_index += len(IconTagEnd)

		cut_str = cut_str[:end_index]

		name_start_idx := strings.Index(cut_str, "\"")
		if name_start_idx == -1 {
			panic(ErrorMissingNameAttr)
		}

		name_end_idx := strings.Index(cut_str[name_start_idx+1:], "\"")
		if name_end_idx == -1 {
			panic(ErrorMalformedTag)
		}

		name_end_idx += name_start_idx + 1

		iconName := cut_str[name_start_idx+1 : name_end_idx]

		iconSVG := LoadIcon(iconName)

		str = str[:i] + iconSVG + str[i+end_index:]
		i += len(iconSVG) - 1
	}

	return []byte(str)
}
