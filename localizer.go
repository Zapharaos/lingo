package lingo

import (
	"sync"

	"golang.org/x/text/language"
)

// LocalizerService defines the interface for handling localizers
type LocalizerService interface {
	GetLocalizer(language language.Tag) (interface{}, bool, error)
	Translate(localizer interface{}, message *Message) (string, bool, error)
	MustTranslate(localizer interface{}, message *Message) string
}

// Message represents a translatable message item
type Message struct {
	ID          string
	Data        interface{}
	PluralCount interface{}
}

// NewMessage creates a new Message instance with the given ID
func NewMessage(id string) *Message {
	return &Message{
		ID: id,
	}
}

// WithData sets the ID for the message
func (m *Message) WithData(data interface{}) *Message {
	m.Data = data
	return m
}

// WithPluralCount sets the plural count for the message
func (m *Message) WithPluralCount(count interface{}) *Message {
	m.PluralCount = count
	return m
}

var (
	_globalServiceMu sync.RWMutex
	_globalService   LocalizerService
)

// SetLocalizerService affect a new repository to the global service singleton
func SetLocalizerService(service LocalizerService) func() {
	_globalServiceMu.Lock()
	defer _globalServiceMu.Unlock()

	prev := _globalService
	_globalService = service
	return func() { SetLocalizerService(prev) }
}

// GetLocalizerService is used to access the global service singleton
func GetLocalizerService() LocalizerService {
	_globalServiceMu.RLock()
	defer _globalServiceMu.RUnlock()
	return _globalService
}

// Directly exposes global localizer implementation

// GetLocalizer Directly exposes the current service GetLocalizer function.
func GetLocalizer(language language.Tag) (interface{}, bool, error) {
	return GetLocalizerService().GetLocalizer(language)
}

// Translate Directly exposes the current service Translate function.
func Translate(localizer interface{}, message *Message) (string, bool, error) {
	return GetLocalizerService().Translate(localizer, message)
}

// MustTranslate Directly exposes the current service MustTranslate function.
func MustTranslate(localizer interface{}, message *Message) string {
	return GetLocalizerService().MustTranslate(localizer, message)
}
