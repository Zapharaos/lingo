[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/zapharaos/lingo)](https://pkg.go.dev/mod/github.com/zapharaos/lingo)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.24.1-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/Zapharaos/lingo)](https://goreportcard.com/report/github.com/Zapharaos/lingo)
![GitHub License](https://img.shields.io/github/license/zapharaos/lingo)

![GitHub Release](https://img.shields.io/github/v/release/zapharaos/lingo)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/zapharaos/lingo/golang.yml)
[![codecov](https://codecov.io/gh/Zapharaos/lingo/graph/badge.svg?token=48KQHAV93M)](https://codecov.io/gh/Zapharaos/lingo)

# lingo

Go library for streamlining language translation handling.

**lingo** automatically discovers translation configuration files in the input path, loads the translations, and sets up a localizer service. Users can either use the default go-i18n implementation or implement their own custom solution. It provides a simple, unified interface for managing internationalization in Go applications.

## Features

- **Auto-discovery**: Automatically finds and loads translation files from specified directories
- **Flexible implementation**: Use the built-in go-i18n implementation or create your own
- **Fallback support**: Graceful fallback to default language when translations are missing
- **Template variables**: Support for dynamic content with template data
- **Pluralization**: Built-in support for plural forms

## Supported File Formats

lingo supports the following translation file formats through go-i18n:
- **TOML** (`.toml`)
- **JSON** (`.json`)
- **YAML** (`.yaml`, `.yml`)

## Installation

```sh
go get github.com/Zapharaos/lingo
```

**Note:** Spit uses [Go Modules](https://go.dev/wiki/Modules) to manage dependencies.

## Usage Examples

#### 1. Create translation files

Create a translation file (e.g., `messages.en.toml`):

```toml
[hello_world]
other = "Hello, World!"
```

#### 2. Initialize and use lingo

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Zapharaos/lingo"
    "golang.org/x/text/language"
)

func main() {
    // Initialize the localizer service
    i18n, err := lingo.NewI18n(
        language.English, // default language
        "config/",        // translations directory
        "messages",       // file prefix
    )
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    
    // Set as global service
    lingo.SetLocalizerService(i18n)
    
    // Translate a message
    result := lingo.MustTranslate(language.English, &lingo.Message{
        ID: "hello_world",
    })
    fmt.Println(result) // Output: Hello, World!
}
```

For comprehensive usage examples including template variables, pluralization, and fallback behavior, see the [examples package](./examples/).

## Development

Install dependencies:
```shell
make dev-deps
```

Run unit tests and generate coverage report:
```shell
make test-unit
```

Run linters:

```shell
make lint
```

Some linter violations can automatically be fixed:

```shell
make fmt
```

## Contributing

We welcome contributions to the lingo library! If you have a bug fix, feature request, or improvement, please open an issue or pull request on GitHub. We appreciate your help in making lingo better for everyone. If you are interested in contributing to the lingo library, please check out our [contributing guidelines](CONTRIBUTING.md) for more information on how to get started.

## License

The project is licensed under the [MIT License](LICENSE).