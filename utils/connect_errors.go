package utils

import (
	"errors"
	"log"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

// ParseConnectViolations extracts field-level validation errors from a Connect error
func ParseConnectViolations(err error) map[string]string {
	fieldErrors := make(map[string]string)

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		return fieldErrors
	}

	// Iterate through error details
	for _, detail := range connectErr.Details() {
		log.Printf("DEBUG: Detail type: %s", detail.Type())

		if detail.Type() != "buf.validate.Violations" {
			continue
		}

		// detail.Bytes() returns already-decoded bytes, not base64
		violationsBytes := detail.Bytes()
		log.Printf("DEBUG: Violations bytes length: %d", len(violationsBytes))

		// Unmarshal as Violations protobuf
		violations := &validate.Violations{}
		if err := proto.Unmarshal(violationsBytes, violations); err != nil {
			log.Printf("DEBUG: Failed to unmarshal violations: %v", err)
			continue
		}

		log.Printf("DEBUG: Found %d violations", len(violations.GetViolations()))

		// Extract field name and message from each violation
		for _, violation := range violations.GetViolations() {
			fieldName := extractFieldName(violation)
			message := violation.GetMessage()
			log.Printf("DEBUG: Violation - field: %s, message: %s", fieldName, message)
			if fieldName != "" {
				fieldErrors[fieldName] = message
			}
		}
	}

	return fieldErrors
}

// extractFieldName gets the field name from violation.field.elements
func extractFieldName(violation *validate.Violation) string {
	field := violation.GetField()
	if field == nil {
		return ""
	}

	elements := field.GetElements()
	if len(elements) == 0 {
		return ""
	}

	// For simple fields, take first element's field name
	if elem := elements[0].GetFieldName(); elem != "" {
		return elem
	}

	return ""
}
