Siberia Env
=================

[![Author](https://img.shields.io/badge/author-@siberia_projects-green.svg)](https://github.com/siberia-projects)
[![Source Code](https://img.shields.io/badge/source-siberia/main-blue.svg)](https://github.com/siberia-projects/siberia-env)
![Version](https://img.shields.io/badge/version-v2.0.0-orange.svg)
[![Coverage Status](https://coveralls.io/repos/github/siberia-projects/siberia-env/badge.svg?branch=main)](https://coveralls.io/github/siberia-projects/siberia-env?branch=main)

## What is it?
Siberia-env is a library written on clear go for working with environment variables

For now the library provides a method to expand all labeled (${ENV_VARIABLE_NAME})
environment variables in a content to their real values

## How to download?

```console
john@doe-pc:~$ go get github.com/siberia-projects/siberia-env
```

## How to use?
 - Label places you want to expand in your content with "${}" symbols sequence
 - Provide a name of a variable a value of which you want to extract
 - (Optional) Provide a default value using ":" symbol (could be empty)
 - Call the "env.ExpandEnvIn" method with the content

## Examples
```yaml
my:
  custom:
    properties:
      headers:
        - key: Cache-Control
          value: ${CACHE_CONTROL:no-cache}
        - key: Content-Type
          value: application/json
```

```go
    package main

    import (
		"github.com/siberia-projects/siberia-env/pkg/env"
		"io"
		"os"
    )

    func main() {
		file, _ := os.Open("path_to_the_yaml_file_above.yaml")
		defer file.Close()

		fileContent, _ := io.ReadAll(file)

		expandedFileContent, _ := env.ExpandEnvIn(fileContent)
		expandedFileContentString := string(expandedFileContent)
		
		println(expandedFileContentString)
    }
```

```console
my:
  custom:
    properties:
      headers:
        - key: Cache-Control
          value: no-cache
        - key: Content-Type
          value: application/json
```
