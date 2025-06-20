// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.887
package input

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "strings"

func Input(props InputProps) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)

		// Generate CSS classes
		classes := inputVariants(props.Class)

		// Set up input attributes
		var inputAttrs templ.Attributes
		inputAttrs = make(templ.Attributes)

		// Copy user-provided attributes first
		for k, v := range props.Attributes {
			inputAttrs[k] = v
		}

		// Standard input attributes
		inputAttrs["data-slot"] = "input"
		inputAttrs["class"] = classes

		if props.Type != "" {
			inputAttrs["type"] = props.Type
		}
		if props.Placeholder != "" {
			inputAttrs["placeholder"] = props.Placeholder
		}
		if props.Value != "" {
			inputAttrs["value"] = props.Value
		}
		if props.Name != "" {
			inputAttrs["name"] = props.Name
		}
		if props.ID != "" {
			inputAttrs["id"] = props.ID
		}
		if props.Disabled {
			inputAttrs["disabled"] = true
		}
		if props.Required {
			inputAttrs["required"] = true
		}

		// Set up data-bind if FormID and Name are provided
		if props.FormID != "" && props.Name != "" {
			// Convert form ID to valid signal name (replace hyphens with underscores)
			signalName := strings.ReplaceAll(props.FormID, "-", "_")
			dataBindValue := signalName + "." + props.Name
			inputAttrs["data-bind"] = dataBindValue

			// Don't try to initialize nested signals here - let the form handle it
			// The data-bind will automatically create the signal if it doesn't exist
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<input")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderAttributes(ctx, templ_7745c5c3_Buffer, inputAttrs)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
