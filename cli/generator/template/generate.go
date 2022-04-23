// DO NOT CHANGE PACKAGE

// Package template provides a template used by copygen to generate custom code.
package template

import (
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// Generate generates code.
// GENERATOR FUNCTION.
// EDITABLE.
// DO NOT REMOVE.
func Generate(gen *models.Generator) (string, error) {
	var content strings.Builder

	content.WriteString(string(gen.Keep) + "\n")
	for i := range gen.Functions {
		content.WriteString(Function(&gen.Functions[i]) + "\n")
	}

	return content.String(), nil
}

// Function provides generated code for a function.
func Function(function *models.Function) string {
	var fn strings.Builder
	fn.WriteString(generateComment(function) + "\n")
	fn.WriteString(generateSignature(function) + "\n")
	fn.WriteString(generateBody(function))
	fn.WriteString(generateReturn(function))
	return fn.String()
}

// generateComment generates a function comment.
func generateComment(function *models.Function) string {
	var toComment strings.Builder
	for i, toType := range function.To {
		if i+1 == len(function.To) {
			toComment.WriteString(toType.Name())
			break
		}

		toComment.WriteString(toType.Name() + ", ")
	}

	var fromComment strings.Builder
	for i, fromType := range function.From {
		if i+1 == len(function.From) {
			fromComment.WriteString(fromType.Name())
			break
		}

		fromComment.WriteString(fromType.Name() + ", ")
	}

	return "// " + function.Name + " copies a " + fromComment.String() + " to a " + toComment.String() + "."
}

// generateSignature generates a function's signature.
func generateSignature(function *models.Function) string {
	return "func " + function.Name + "(" + generateParameters(function) + ") {"
}

// generateParameters generates the parameters of a function.
func generateParameters(function *models.Function) string {
	var parameters strings.Builder
	for _, toType := range function.To {
		parameters.WriteString(toType.Field.VariableName + " " + toType.Name() + ", ")
	}

	for i, fromType := range function.From {
		if i+1 == len(function.From) {
			parameters.WriteString(fromType.Field.VariableName + " " + fromType.Name())
			break
		}

		parameters.WriteString(fromType.Field.VariableName + " " + fromType.Name() + ", ")
	}

	return parameters.String()
}

// generateBody generates the body of a function.
func generateBody(function *models.Function) string {
	var body strings.Builder

	// Assign fields to ToType(s).
	for i, toType := range function.To {
		body.WriteString("// " + toType.Name() + " fields\n")

		for _, toField := range toType.Field.Fields {
			body.WriteString(toField.FullVariableName("") + " = ")

			fromField := toField.From
			switch {
			case fromField.Options.Convert != "":
				body.WriteString(fromField.Options.Convert + "(" + fromField.FullVariableName("") + ")\n")
			case toField.Definition != fromField.Definition:
				// match alias types to respective basic types.
				if toField.Package != "" {
					body.WriteString(toField.Package + ".")
				}
				body.WriteString(toField.Definition + "(" + fromField.FullVariableName("") + ")" + "\n")
			default:
				body.WriteString(fromField.FullVariableName("") + "\n")
			}
		}

		if i+1 != len(function.To) {
			body.WriteString("\n")
		}
	}

	return body.String()
}

// generateReturn generates a return statement for the function.
func generateReturn(function *models.Function) string {
	return "}"
}
