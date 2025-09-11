package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	if inputStr == "" {
		return "", nil
	}

	// Предобработка строки в []any для более удобной работы

	var (
		textRune       = []rune(inputStr)
		processedSlice = make([]any, 0, len(inputStr))
		skipCount      int
	)
	for i, r := range textRune {
		// Пропуск
		if skipCount >= 1 {
			skipCount--
			continue
		}

		// runeToNum преобразует руну в int64.
		runeToNum := func(r rune) int64 {
			num, _ := strconv.ParseInt(string(r), 10, 64)
			return num
		}

		switch {
		// Символ экранирования "\"
		case r == '\\':
			// Если последний символ
			if i == len(textRune)-1 {
				return "", ErrInvalidString
			}
			skipCount++
			// Экранировать можно только цифру или слэш
			switch rNext := textRune[i+1]; {
			case rNext == '\\':
				processedSlice = append(processedSlice, string('\\'))
			case unicode.IsDigit(rNext):
				processedSlice = append(processedSlice, string(rNext))
			default:
				return "", ErrInvalidString
			}
		// Символ числа
		case unicode.IsDigit(r):
			processedSlice = append(processedSlice, runeToNum(r))
		// Любой другой символ
		default:
			processedSlice = append(processedSlice, string(r))
		}
	}

	var b strings.Builder
	for i, v := range processedSlice {
		if i == len(processedSlice)-1 {
			vStr, ok := v.(string)
			if ok {
				b.WriteString(vStr)
			}
			break
		}

		// Следующий символ и является ли он числом
		nextVInt, isNextInt := processedSlice[i+1].(int64)

		switch v := v.(type) {
		case int64:
			// Обработка некорретной строки
			if i == 0 || isNextInt {
				return "", ErrInvalidString
			}
		case string:
			var (
				vStr     = v
				strCount = 1
			)
			if isNextInt {
				strCount = int(nextVInt)
			}
			b.WriteString(strings.Repeat(vStr, strCount))
		}
	}

	fmt.Println(b.String())

	return b.String(), nil
}
