package scan

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// NamedQuery returns the query renamed
func NamedQuery(query string) (string, []string) {
	var (
		buffer = &bytes.Buffer{}
		next   int
	)

	for index := strings.Index(query, "?"); index != -1; index = strings.Index(query, "?") {
		fmt.Fprint(buffer, query[:index])
		fmt.Fprint(buffer, ":")
		fmt.Fprint(buffer, fmt.Sprintf("arg%d", next))

		query = query[index+1:]
		next++
	}

	fmt.Fprint(buffer, query)

	query = buffer.String()

	var (
		params     = []string{}
		underscore = '_'
		runes      = []*unicode.RangeTable{
			unicode.Letter,
			unicode.Digit,
		}
	)

	scanner := bufio.NewScanner(buffer)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		for index := 0; index < len(data); index++ {
			if data[index] == ':' {
				return index + 1, data[:index+1], nil
			}
		}

		if !atEOF {
			return 0, nil, nil
		}

		return 0, data, bufio.ErrFinalToken
	})

	for scanner.Scan() {
		part := scanner.Text()
		tail := part[len(part)-1]

		for tail == ':' {
			if !scanner.Scan() {
				break
			}

			part = scanner.Text()
			tail = part[len(part)-1]

			param := &bytes.Buffer{}

			for _, ch := range part {
				if !unicode.IsOneOf(runes, ch) && ch != underscore {
					break
				}

				fmt.Fprint(param, string(ch))
			}

			params = append(params, param.String())
		}
	}

	return query, params
}
