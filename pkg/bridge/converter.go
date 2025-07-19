package bridge

import (
	"context"
	"fmt"
)

// CharacterConverter defines the interface for converting between
// different character format systems and the universal format.
type CharacterConverter interface {
	// Name returns the name of this converter (e.g., "character.ai", "roleplay")
	Name() string

	// CanConvert checks if this converter can handle the given data
	CanConvert(data interface{}) bool

	// ToUniversal converts from the source format to UniversalCharacter
	ToUniversal(ctx context.Context, data interface{}) (*UniversalCharacter, error)

	// FromUniversal converts from UniversalCharacter to the target format
	FromUniversal(ctx context.Context, char *UniversalCharacter) (interface{}, error)
}

// ConverterRegistry manages available character converters.
type ConverterRegistry struct {
	converters map[string]CharacterConverter
}

// NewConverterRegistry creates a new converter registry.
func NewConverterRegistry() *ConverterRegistry {
	return &ConverterRegistry{
		converters: make(map[string]CharacterConverter),
	}
}

// Register adds a converter to the registry.
func (r *ConverterRegistry) Register(converter CharacterConverter) error {
	name := converter.Name()
	if _, exists := r.converters[name]; exists {
		return fmt.Errorf("converter %s already registered", name)
	}
	r.converters[name] = converter
	return nil
}

// Get retrieves a converter by name.
func (r *ConverterRegistry) Get(name string) (CharacterConverter, error) {
	converter, exists := r.converters[name]
	if !exists {
		return nil, fmt.Errorf("converter %s not found", name)
	}
	return converter, nil
}

// FindConverter attempts to find a converter that can handle the given data.
func (r *ConverterRegistry) FindConverter(data interface{}) (CharacterConverter, error) {
	for _, converter := range r.converters {
		if converter.CanConvert(data) {
			return converter, nil
		}
	}
	return nil, fmt.Errorf("no converter found for data type %T", data)
}

// List returns all registered converter names.
func (r *ConverterRegistry) List() []string {
	names := make([]string, 0, len(r.converters))
	for name := range r.converters {
		names = append(names, name)
	}
	return names
}

// ConversionError represents an error during character conversion.
type ConversionError struct {
	Source string
	Target string
	Field  string
	Err    error
}

func (e *ConversionError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("conversion error from %s to %s (field: %s): %v", 
			e.Source, e.Target, e.Field, e.Err)
	}
	return fmt.Sprintf("conversion error from %s to %s: %v", 
		e.Source, e.Target, e.Err)
}

// ConversionWarning represents a non-fatal issue during conversion.
type ConversionWarning struct {
	Field   string
	Message string
}

// ConversionResult contains the result of a conversion operation.
type ConversionResult struct {
	Character *UniversalCharacter
	Warnings  []ConversionWarning
}

// BaseConverter provides common functionality for converters.
type BaseConverter struct {
	name string
}

// NewBaseConverter creates a new base converter.
func NewBaseConverter(name string) *BaseConverter {
	return &BaseConverter{name: name}
}

// Name returns the converter name.
func (b *BaseConverter) Name() string {
	return b.name
}