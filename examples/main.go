package main

import (
	"fmt"
	"log"

	"github.com/Zapharaos/lingo"
	"golang.org/x/text/language"
)

func main() {
	// Step 1: Initialize the localizer service
	fmt.Println("1. Initializing i18n service...")
	i18n, err := lingo.NewI18n(
		language.English, // default language
		"config/",        // translations directory
		"messages",       // file prefix
	)
	if err != nil {
		log.Fatalf("Failed to initialize localizer service: %v", err)
	}

	// Replace the global localizer service instance
	lingo.SetLocalizerService(i18n)
	fmt.Println("✓ Localizer service initialized successfully")

	// Step 2: Demonstrate basic translation usage
	fmt.Println("\n2. Basic Translation Examples:")
	demonstrateBasicTranslation()

	// Step 3: Demonstrate MustTranslate usage
	fmt.Println("\n3. MustTranslate Examples:")
	demonstrateMustTranslate()

	// Step 4: Demonstrate template variables
	fmt.Println("\n4. Template Variables Examples:")
	demonstrateTemplateVariables()

	// Step 5: Demonstrate pluralization
	fmt.Println("\n5. Pluralization Examples:")
	demonstratePluralization()

	// Step 6: Demonstrate fallback behavior
	fmt.Println("\n6. Fallback Behavior:")
	demonstrateFallback()
}

func demonstrateBasicTranslation() {
	languages := []language.Tag{language.English, language.French, language.Spanish, language.German}

	for _, lang := range languages {
		localizer, found, err := lingo.GetLocalizer(lang)
		if err != nil {
			fmt.Printf("Error getting localizer for %s: %v\n", lang, err)
			continue
		}

		result, success, err := lingo.Translate(localizer, lingo.NewMessage("hello_world"))

		foundStr := "✓"
		if !found {
			foundStr = "⚠ (fallback)"
		}

		if err != nil {
			fmt.Printf("  %s: Error - %v\n", lang, err)
		} else if success {
			fmt.Printf("  %s %s: %s\n", foundStr, lang, result)
		} else {
			fmt.Printf("  %s %s: %s (translation failed)\n", foundStr, lang, result)
		}
	}
}

func demonstrateMustTranslate() {
	languages := []language.Tag{language.English, language.French}

	for _, lang := range languages {
		localizer, _, err := lingo.GetLocalizer(lang)
		if err != nil {
			fmt.Printf("Error getting localizer for %s: %v\n", lang, err)
			continue
		}

		// Safe usage of MustTranslate - we know this key exists
		result := lingo.MustTranslate(localizer, lingo.NewMessage("goodbye"))
		fmt.Printf("  %s: %s\n", lang, result)
	}

	// Example of error handling with MustTranslate
	fmt.Println("\n  Demonstrating MustTranslate panic (commented out for safety):")
	fmt.Println("  // This would panic: lingo.MustTranslate(localizer, lingo.NewMessage(\"nonexistent_key\"))")
}

func demonstrateTemplateVariables() {
	localizer, _, _ := lingo.GetLocalizer(language.English)

	// Simple template variable
	message := lingo.NewMessage("welcome_user").
		WithData(map[string]interface{}{
			"Name": "John Doe",
		})

	result := lingo.MustTranslate(localizer, message)
	fmt.Printf("  English (with name): %s\n", result)

	// Multiple variables
	message = lingo.NewMessage("user_stats").
		WithData(map[string]interface{}{
			"Name":  "Alice",
			"Posts": 42,
			"Likes": 128,
		})
	result = lingo.MustTranslate(localizer, message)
	fmt.Printf("  English (with stats): %s\n", result)

	// French version
	localizerFr, _, _ := lingo.GetLocalizer(language.French)
	result = lingo.MustTranslate(localizerFr, message)
	fmt.Printf("  French (with stats): %s\n", result)
}

func demonstratePluralization() {
	localizer, _, _ := lingo.GetLocalizer(language.English)

	pluralCounts := []int{0, 1, 2, 5}

	for _, count := range pluralCounts {
		var message *lingo.Message

		// Handle zero case separately since English CLDR doesn't support "zero" category
		if count == 0 {
			message = lingo.NewMessage("no_items")
		} else {
			message = lingo.NewMessage("item_count").
				WithPluralCount(count).
				WithData(map[string]interface{}{
					"Count": count,
				})
		}

		result := lingo.MustTranslate(localizer, message)
		fmt.Printf("  %d items: %s\n", count, result)
	}

	// French pluralization (different rules)
	fmt.Println("\n  French pluralization:")
	localizerFr, _, _ := lingo.GetLocalizer(language.French)
	for _, count := range pluralCounts {
		var message *lingo.Message

		// Handle zero case for French too
		if count == 0 {
			message = lingo.NewMessage("no_items")
		} else {
			message = lingo.NewMessage("item_count").
				WithPluralCount(count).
				WithData(map[string]interface{}{
					"Count": count,
				})
		}

		result := lingo.MustTranslate(localizerFr, message)
		fmt.Printf("  %d items: %s\n", count, result)
	}
}

func demonstrateFallback() {
	// Try to get a localizer for a language that doesn't exist
	unsupportedLang := language.MustParse("ja") // Japanese - not in our examples files

	localizer, found, err := lingo.GetLocalizer(unsupportedLang)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !found {
		fmt.Printf("Japanese localizer not found, falling back to default (English)\n")
	}

	message := lingo.NewMessage("hello_world")
	result := lingo.MustTranslate(localizer, message)
	fmt.Printf("  Result: %s\n", result)

	// Try to translate a key that doesn't exist
	fmt.Println("\n  When translating a missing key:")
	localizerEn, _, _ := lingo.GetLocalizer(language.English)
	message = lingo.NewMessage("nonexistent_key")
	result, success, err := lingo.Translate(localizerEn, message)

	if err != nil {
		fmt.Printf("  Error translating missing key: %v\n", err)
	} else if !success {
		fmt.Printf("  Translation failed, returned: %s\n", result)
	}
}
