// Copyright (c) 2024 Nikolai Osipov <nao99.dev@gmail.com>
//
// All rights are reserved
// Information about license can be found in the LICENSE file

package env

import (
	"bytes"
	"fmt"
	"os"
)

const (
	environmentVariableBeginningSymbol byte = '$'
	environmentVariableEndingSymbol    byte = '}'

	environmentVariableKeyDefaultValueSeparator byte = ':'
)

// ExpandEnvIn finds all env variable mentions in passed content
// and replaces it with real env values
//
// Returns an error when env doesn't exist
// and default value is not presented
//
// e.g.
//
// -----------------------------------------------------
// CACHE_CONTROL = no-cache
// ${CACHE_CONTROL:yes-please-cache} => no-cache
// -----------------------------------------------------
//
// -----------------------------------------------------
// CACHE_CONTROL = no-cache
// ${CACHE_CONTROL} => no-cache
// -----------------------------------------------------
//
// -----------------------------------------------------
// ${CACHE_CONTROL:yes-please-cache} => yes-please-cache
// -----------------------------------------------------
//
// -----------------------------------------------------
// ${CACHE_CONTROL:} => ""
// -----------------------------------------------------
//
// -----------------------------------------------------
// ${CACHE_CONTROL} => error
// -----------------------------------------------------
func ExpandEnvIn(content []byte) ([]byte, error) {
	if content == nil {
		return nil, fmt.Errorf("unable to expand: a content is nil")
	}

	newContent := make([]byte, len(content))
	newContentIndex := 0

	for i := 0; i < len(content); i++ {
		if content[i] != environmentVariableBeginningSymbol {
			newContent[newContentIndex] = content[i]
			newContentIndex++

			continue
		}

		envKeyBytes := make([]byte, 0)

		j := i + 2
		for content[j] != environmentVariableEndingSymbol {
			envKeyBytes = append(envKeyBytes, content[j])
			j++
		}

		envKeyByteParts := bytes.Split(envKeyBytes, []byte{environmentVariableKeyDefaultValueSeparator})

		envKey := string(envKeyByteParts[0])
		envValue := os.Getenv(envKey)

		if envValue == "" {
			if len(envKeyByteParts) == 2 {
				envValue = string(envKeyByteParts[1])
			} else {
				return nil, fmt.Errorf("unable to expand \"%s\" variable: it doesn't exist in env", envKey)
			}
		}

		envValueBytes := []byte(envValue)
		isNewContentIncreasingNeeded := false

		for k := 0; k < len(envValueBytes); k++ {
			if newContentIndex == len(newContent) {
				isNewContentIncreasingNeeded = true
			}

			if isNewContentIncreasingNeeded {
				newContent = append(newContent, envValueBytes[k])

				isNewContentIncreasingNeeded = false
				newContentIndex++

				continue
			}

			newContent[newContentIndex] = envValueBytes[k]
			newContentIndex++
		}

		i += len(envKeyBytes) + 2
	}

	return newContent[:newContentIndex], nil
}
