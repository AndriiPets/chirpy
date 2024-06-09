package utils

import "strings"

var badWorsdMap = map[string]int{
	"kerfuffle": 0,
	"sharbert":0,
	"fornax":0,
}

func CleanInput(text string) string {
	str := strings.Split(text, " ") 
	var clean []string
	for _,word := range str {
		_, ok := badWorsdMap[strings.ToLower(word)]
		if ok {
			clean = append(clean, strings.Repeat("*", len(word)))
			continue
		}
		clean = append(clean, word)
	}

	return strings.Join(clean, " ")
}