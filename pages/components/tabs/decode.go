package tabs

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FormDecoder handles automatic parsing of HTTP form data into protobuf messages
type FormDecoder struct{}

// NewFormDecoder creates a new form decoder instance
func NewFormDecoder() *FormDecoder {
	return &FormDecoder{}
}

// DecodeForm parses HTTP form data into a protobuf message using reflection
func (d *FormDecoder) DecodeForm(r *http.Request, msg proto.Message) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	// Get the protobuf reflection descriptor
	reflectMsg := msg.ProtoReflect()
	descriptor := reflectMsg.Descriptor()

	// Iterate through all fields in the protobuf message
	fields := descriptor.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := string(field.Name())

		// Get form value for this field
		formValue := r.FormValue(fieldName)
		if formValue == "" {
			continue // Skip empty values
		}

		// Set the field value based on its type
		if err := d.setFieldValue(reflectMsg, field, formValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return nil
}

// setFieldValue sets a protobuf field value based on its type
func (d *FormDecoder) setFieldValue(msg protoreflect.Message, field protoreflect.FieldDescriptor, value string) error {
	switch field.Kind() {
	case protoreflect.StringKind:
		msg.Set(field, protoreflect.ValueOfString(value))

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		intVal, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid int32 value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfInt32(int32(intVal)))

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int64 value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfInt64(intVal))

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		uintVal, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid uint32 value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfUint32(uint32(uintVal)))

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint64 value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfUint64(uintVal))

	case protoreflect.FloatKind:
		floatVal, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfFloat32(float32(floatVal)))

	case protoreflect.DoubleKind:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid double value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfFloat64(floatVal))

	case protoreflect.BoolKind:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool value: %s", value)
		}
		msg.Set(field, protoreflect.ValueOfBool(boolVal))

	case protoreflect.BytesKind:
		msg.Set(field, protoreflect.ValueOfBytes([]byte(value)))

	case protoreflect.EnumKind:
		// Handle enum by name or number
		enumDesc := field.Enum()
		var enumVal protoreflect.EnumNumber

		// Try to parse as number first
		if num, err := strconv.ParseInt(value, 10, 32); err == nil {
			enumVal = protoreflect.EnumNumber(num)
		} else {
			// Try to find by name
			enumValue := enumDesc.Values().ByName(protoreflect.Name(value))
			if enumValue == nil {
				return fmt.Errorf("invalid enum value: %s", value)
			}
			enumVal = enumValue.Number()
		}
		msg.Set(field, protoreflect.ValueOfEnum(enumVal))

	case protoreflect.MessageKind:
		// For nested messages, we'd need more complex handling
		return fmt.Errorf("nested message fields not supported yet")

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// DecodeFormWithValidation parses form data and returns validation errors
func (d *FormDecoder) DecodeFormWithValidation(r *http.Request, msg proto.Message) (*FormResponse, error) {
	if err := d.DecodeForm(r, msg); err != nil {
		return &FormResponse{
			Success: false,
			Message: "Form validation failed",
			Errors:  []string{err.Error()},
		}, nil
	}

	// Additional validation can be added here
	if err := d.validateMessage(msg); err != nil {
		return &FormResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &FormResponse{
		Success: true,
		Message: "Form processed successfully",
		Errors:  []string{},
	}, nil
}

// validateMessage performs custom validation on the protobuf message
func (d *FormDecoder) validateMessage(msg proto.Message) error {
	// Get the protobuf reflection descriptor
	reflectMsg := msg.ProtoReflect()
	descriptor := reflectMsg.Descriptor()

	// Iterate through fields and perform validation
	fields := descriptor.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := string(field.Name())

		// Check if required fields are set
		if field.HasPresence() && !reflectMsg.Has(field) {
			return fmt.Errorf("required field %s is missing", fieldName)
		}

		// Validate string fields
		if field.Kind() == protoreflect.StringKind && reflectMsg.Has(field) {
			value := reflectMsg.Get(field).String()
			if err := d.validateStringField(fieldName, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateStringField performs validation on string fields
func (d *FormDecoder) validateStringField(fieldName, value string) error {
	// Trim whitespace
	value = strings.TrimSpace(value)

	// Basic validation rules
	switch fieldName {
	case "name":
		if len(value) < 2 {
			return fmt.Errorf("name must be at least 2 characters long")
		}
		if len(value) > 50 {
			return fmt.Errorf("name must be less than 50 characters")
		}

	case "username":
		if len(value) < 3 {
			return fmt.Errorf("username must be at least 3 characters long")
		}
		if len(value) > 20 {
			return fmt.Errorf("username must be less than 20 characters")
		}
		// Check for valid username characters
		for _, char := range value {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '_' || char == '-') {
				return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
			}
		}

	case "current_password", "new_password":
		if len(value) < 8 {
			return fmt.Errorf("password must be at least 8 characters long")
		}
		if len(value) > 128 {
			return fmt.Errorf("password must be less than 128 characters")
		}
	}

	return nil
}

// Helper functions for easy usage

// DecodeAccountForm decodes form data into AccountForm
func DecodeAccountForm(r *http.Request) (*AccountForm, *FormResponse, error) {
	decoder := NewFormDecoder()
	form := &AccountForm{}

	response, err := decoder.DecodeFormWithValidation(r, form)
	if err != nil {
		return nil, nil, err
	}

	return form, response, nil
}

// DecodePasswordForm decodes form data into PasswordForm
func DecodePasswordForm(r *http.Request) (*PasswordForm, *FormResponse, error) {
	decoder := NewFormDecoder()
	form := &PasswordForm{}

	response, err := decoder.DecodeFormWithValidation(r, form)
	if err != nil {
		return nil, nil, err
	}

	return form, response, nil
}
