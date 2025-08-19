// Package lingo streamlines language translations.
//
// Features:
//   - Automatically finds and loads translation files from specified directories
//   - Follows BCP 47 language tags
//   - Graceful fallback to default language when translations are missing
//   - Support for dynamic content with template data
//   - Supports pluralization
//   - Defaults to https://github.com/nicksnyder/go-i18n for translation management
//   - Allows you to implement your own custom solution by implementing the LocalizerService interface
//
// Supported extensions (specific to go-i18n):
//   - JSON
//   - YAML / YML
//   - TOML
//
// For more details, see README.md.
package lingo
