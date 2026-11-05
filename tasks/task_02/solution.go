package main

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	shift %= len(runes)
	if shift < 0 {
		shift += len(runes)
	}

	res := make([]rune, len(runes))
	for i, r := range runes {
		newIndex := i - shift
		if newIndex < 0 {
			newIndex += len(runes)
		}
		res[newIndex] = r
	}

	return string(res)
}
