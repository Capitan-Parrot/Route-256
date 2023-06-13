package stringUtils

import "strconv"

// ConvertSlice конвертирует слайс uint64 в string
func ConvertSlice(integers []*uint64) []string {
	result := make([]string, len(integers))
	for index, integer := range integers {
		result[index] = strconv.FormatUint(*integer, 10)
	}
	return result
}
