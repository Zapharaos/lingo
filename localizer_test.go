package lingo

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

// TestMessage tests the Message struct and its methods
func TestMessage(t *testing.T) {
	t.Run("NewMessage creates message with ID", func(t *testing.T) {
		msg := NewMessage("test.id")
		assert.Equal(t, "test.id", msg.ID)
		assert.Nil(t, msg.Data)
		assert.Nil(t, msg.PluralCount)
	})

	t.Run("WithData sets data and returns message", func(t *testing.T) {
		msg := NewMessage("test.id")
		data := map[string]string{"name": "John"}

		result := msg.WithData(data)

		assert.Equal(t, msg, result) // Should return the same instance
		assert.Equal(t, data, msg.Data)
	})

	t.Run("WithPluralCount sets plural count and returns message", func(t *testing.T) {
		msg := NewMessage("test.id")
		count := 5

		result := msg.WithPluralCount(count)

		assert.Equal(t, msg, result) // Should return the same instance
		assert.Equal(t, count, msg.PluralCount)
	})

	t.Run("Method chaining works correctly", func(t *testing.T) {
		data := map[string]string{"name": "John"}
		count := 3

		msg := NewMessage("test.id").
			WithData(data).
			WithPluralCount(count)

		assert.Equal(t, "test.id", msg.ID)
		assert.Equal(t, data, msg.Data)
		assert.Equal(t, count, msg.PluralCount)
	})

	t.Run("WithData can handle different data types", func(t *testing.T) {
		tests := []struct {
			name string
			data interface{}
		}{
			{"string", "hello"},
			{"int", 42},
			{"map", map[string]interface{}{"key": "value"}},
			{"slice", []string{"a", "b", "c"}},
			{"nil", nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				msg := NewMessage("test").WithData(tt.data)
				assert.Equal(t, tt.data, msg.Data)
			})
		}
	})

	t.Run("WithPluralCount can handle different count types", func(t *testing.T) {
		tests := []struct {
			name  string
			count interface{}
		}{
			{"int", 42},
			{"int64", int64(42)},
			{"float64", 42.5},
			{"string", "42"},
			{"nil", nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				msg := NewMessage("test").WithPluralCount(tt.count)
				assert.Equal(t, tt.count, msg.PluralCount)
			})
		}
	})
}

// TestSetLocalizerService tests the SetLocalizerService function
// It verifies that the global service can be replaced and restored correctly.
func TestSetLocalizerService(t *testing.T) {
	t.Run("Basic replacement and restoration", func(t *testing.T) {
		// Mock
		ctrl := gomock.NewController(t)
		m := NewMockLocalizerService(ctrl)
		defer ctrl.Finish()

		// Replace the global service with a mock service
		restore := SetLocalizerService(m)

		// Ensure the global service is replaced
		assert.Equal(t, m, GetLocalizerService())

		// Restore the previous global service
		restore()
		assert.NotEqual(t, m, GetLocalizerService())
	})

	t.Run("Multiple replacements work correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock1 := NewMockLocalizerService(ctrl)
		mock2 := NewMockLocalizerService(ctrl)
		defer ctrl.Finish()

		// Set first service
		restore1 := SetLocalizerService(mock1)
		assert.Equal(t, mock1, GetLocalizerService())

		// Set second service
		restore2 := SetLocalizerService(mock2)
		assert.Equal(t, mock2, GetLocalizerService())

		// Restore to first service
		restore2()
		assert.Equal(t, mock1, GetLocalizerService())

		// Restore to original
		restore1()
		assert.NotEqual(t, mock1, GetLocalizerService())
		assert.NotEqual(t, mock2, GetLocalizerService())
	})

	t.Run("Restore function can be called multiple times safely", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := NewMockLocalizerService(ctrl)
		defer ctrl.Finish()

		originalService := GetLocalizerService()
		restore := SetLocalizerService(mock)

		assert.Equal(t, mock, GetLocalizerService())

		// Call restore multiple times
		restore()
		assert.Equal(t, originalService, GetLocalizerService())

		restore() // Should not panic or cause issues
		assert.Equal(t, originalService, GetLocalizerService())
	})
}

// TestGetLocalizerService tests the GetLocalizerService function
// It verifies that the global service can be accessed correctly.
func TestGetLocalizerService(t *testing.T) {
	t.Run("Returns the current global service", func(t *testing.T) {
		// Mock
		ctrl := gomock.NewController(t)
		m := NewMockLocalizerService(ctrl)
		defer ctrl.Finish()

		// Replace the global service with a mock service
		restore := SetLocalizerService(m)
		defer restore()

		// Access the global service
		service := GetLocalizerService()
		assert.Equal(t, m, service)
	})

	t.Run("Returns nil when no service is set", func(t *testing.T) {
		// Temporarily clear the global service
		restore := SetLocalizerService(nil)
		defer restore()

		service := GetLocalizerService()
		assert.Nil(t, service)
	})
}

// TestGlobalFunctions tests the global wrapper functions
func TestGlobalFunctions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockLocalizerService(ctrl)
	defer ctrl.Finish()

	restore := SetLocalizerService(mockService)
	defer restore()

	t.Run("GetLocalizer calls service GetLocalizer", func(t *testing.T) {
		lang := language.English
		expectedLocalizer := "test-localizer"
		expectedFound := true
		expectedErr := assert.AnError

		mockService.EXPECT().
			GetLocalizer(lang).
			Return(expectedLocalizer, expectedFound, expectedErr)

		localizer, found, err := GetLocalizer(lang)

		assert.Equal(t, expectedLocalizer, localizer)
		assert.Equal(t, expectedFound, found)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Translate calls service Translate", func(t *testing.T) {
		localizer := "test-localizer"
		message := NewMessage("test.message")
		expectedTranslation := "Test Translation"
		expectedFound := true
		expectedErr := assert.AnError

		mockService.EXPECT().
			Translate(localizer, message).
			Return(expectedTranslation, expectedFound, expectedErr)

		translation, found, err := Translate(localizer, message)

		assert.Equal(t, expectedTranslation, translation)
		assert.Equal(t, expectedFound, found)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("MustTranslate calls service MustTranslate", func(t *testing.T) {
		localizer := "test-localizer"
		message := NewMessage("test.message")
		expectedTranslation := "Test Translation"

		mockService.EXPECT().
			MustTranslate(localizer, message).
			Return(expectedTranslation)

		translation := MustTranslate(localizer, message)

		assert.Equal(t, expectedTranslation, translation)
	})
}

// TestConcurrentAccess tests thread safety of the global service
func TestConcurrentAccess(t *testing.T) {
	t.Run("Concurrent reads are safe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := NewMockLocalizerService(ctrl)
		defer ctrl.Finish()

		restore := SetLocalizerService(mockService)
		defer restore()

		const numReaders = 100
		var wg sync.WaitGroup
		results := make([]LocalizerService, numReaders)

		wg.Add(numReaders)
		for i := 0; i < numReaders; i++ {
			go func(idx int) {
				defer wg.Done()
				results[idx] = GetLocalizerService()
			}(i)
		}

		wg.Wait()

		// All reads should return the same service
		for i, result := range results {
			assert.Equal(t, mockService, result, "Reader %d got unexpected result", i)
		}
	})

	t.Run("Concurrent writes are safe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const numWriters = 10
		var wg sync.WaitGroup
		mocks := make([]*MockLocalizerService, numWriters)
		restoreFuncs := make([]func(), numWriters)

		// Create mock services
		for i := 0; i < numWriters; i++ {
			mocks[i] = NewMockLocalizerService(ctrl)
		}

		wg.Add(numWriters)
		for i := 0; i < numWriters; i++ {
			go func(idx int) {
				defer wg.Done()
				restoreFuncs[idx] = SetLocalizerService(mocks[idx])
			}(i)
		}

		wg.Wait()

		// The final service should be one of the mocks
		finalService := GetLocalizerService()
		found := false
		for _, mock := range mocks {
			if finalService == mock {
				found = true
				break
			}
		}
		assert.True(t, found, "Final service should be one of the set mocks")

		// Clean up (call restore functions)
		for _, restore := range restoreFuncs {
			if restore != nil {
				restore()
			}
		}
	})

	t.Run("Mixed concurrent reads and writes are safe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := NewMockLocalizerService(ctrl)
		defer ctrl.Finish()

		const numOperations = 50
		var wg sync.WaitGroup

		originalRestore := SetLocalizerService(mockService)
		defer originalRestore()

		wg.Add(numOperations * 2) // readers + writers

		// Start readers
		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				service := GetLocalizerService()
				assert.NotNil(t, service) // Should always get some service
			}()
		}

		// Start writers
		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				newMock := NewMockLocalizerService(ctrl)
				restore := SetLocalizerService(newMock)
				// Immediately restore to avoid leaving test in inconsistent state
				restore()
			}()
		}

		wg.Wait()
		// Test should complete without data races or panics
	})
}

// TestGlobalFunctionsWithNilService tests behavior when no service is set
func TestGlobalFunctionsWithNilService(t *testing.T) {
	// Temporarily clear the global service
	restore := SetLocalizerService(nil)
	defer restore()

	t.Run("GetLocalizer panics with nil service", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _, _ = GetLocalizer(language.English)
		})
	})

	t.Run("Translate panics with nil service", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _, _ = Translate("localizer", NewMessage("test"))
		})
	})

	t.Run("MustTranslate panics with nil service", func(t *testing.T) {
		assert.Panics(t, func() {
			MustTranslate("localizer", NewMessage("test"))
		})
	})
}

// TestComplexMessageChaining tests complex message building scenarios
func TestComplexMessageChaining(t *testing.T) {
	t.Run("Empty message ID", func(t *testing.T) {
		msg := NewMessage("")
		assert.Equal(t, "", msg.ID)
	})

	t.Run("Message with complex data structure", func(t *testing.T) {
		complexData := map[string]interface{}{
			"user": map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			"items": []map[string]interface{}{
				{"name": "Item 1", "price": 10.50},
				{"name": "Item 2", "price": 25.00},
			},
			"total":    35.50,
			"currency": "USD",
		}

		msg := NewMessage("order.summary").
			WithData(complexData).
			WithPluralCount(2)

		assert.Equal(t, "order.summary", msg.ID)
		assert.Equal(t, complexData, msg.Data)
		assert.Equal(t, 2, msg.PluralCount)
	})

	t.Run("Overwriting data and plural count", func(t *testing.T) {
		msg := NewMessage("test").
			WithData("first").
			WithPluralCount(1).
			WithData("second").
			WithPluralCount(2)

		assert.Equal(t, "second", msg.Data)
		assert.Equal(t, 2, msg.PluralCount)
	})
}
