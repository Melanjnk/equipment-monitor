package dtos

import "strings"

func isWhiteSpace(character byte) bool {
	switch character {
		case ' ', '\t', '\v', '\n', '\r':
			return true
		default:
			return false
	}
}

func normalizeFieldName(fieldName *string) bool {
	for end := len(*fieldName); end > 0; {
		end--
		if !isWhiteSpace((*fieldName)[end]) {
			for start := 0; start < end; start++ {
				if !isWhiteSpace((*fieldName)[start]) {
					b := start <= end && (*fieldName)[start] == '-'
					if b {
						start++
					}
					*fieldName = strings.ToLower((*fieldName)[start:end + 1])
					return b
				}
			}
		}
	}
	*fieldName = ``
	return false
}
