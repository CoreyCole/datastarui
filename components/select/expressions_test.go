package selectcomponent

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSStringLiteralEncodesJavaScriptString(t *testing.T) {
	got := jsStringLiteral(`project 'alpha' \\ beta`)
	if got != `"project 'alpha' \\\\ beta"` {
		t.Fatalf("literal = %s", got)
	}
}

func TestSelectItemClickHandlerUsesEncodedValueLiteral(t *testing.T) {
	got := NewSelectItemHandler("project_filter", `project 'alpha' \\ beta`).BuildClickHandler()
	if strings.Contains(got, `$project_filter.value = 'project 'alpha'`) {
		t.Fatalf("handler used unsafe raw single-quoted value: %s", got)
	}
	if !strings.Contains(got, `$project_filter.value = "project 'alpha' \\\\ beta"`) {
		t.Fatalf("handler missing encoded value literal: %s", got)
	}
}

func TestSelectOnChangeEffectAllowsEmptyValue(t *testing.T) {
	got := selectOnChangeEffect("project_filter", "this.closest('form').requestSubmit()")
	if strings.Contains(got, "$project_filter.value &&") {
		t.Fatalf("effect still blocks empty value changes: %s", got)
	}
	for _, want := range []string{
		"$project_filter.value !== $project_filter._lastValue",
		"$project_filter._lastValue = $project_filter.value",
		"this.closest('form').requestSubmit()",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("effect missing %q: %s", want, got)
		}
	}
}

func TestSelectValueEqualsExpressionUsesEncodedValueLiteral(t *testing.T) {
	got := selectValueEqualsExpression("$project_filter.value", `project 'alpha' \\ beta`)
	if strings.Contains(got, `=== 'project 'alpha'`) {
		t.Fatalf("comparison used unsafe raw single-quoted value: %s", got)
	}
	if !strings.Contains(got, `$project_filter.value === "project 'alpha' \\\\ beta"`) {
		t.Fatalf("comparison missing encoded value literal: %s", got)
	}
}

func TestSelectRendersEmptyOptionAndOnChangeWithoutTruthyGuard(t *testing.T) {
	var body bytes.Buffer
	err := Select(SelectArgs{
		ID:          "project_filter",
		Name:        "project",
		Value:       "",
		Options:     []SelectOptionArgs{{Value: "", Label: "All projects"}, {Value: `project 'alpha' \\ beta`, Label: "Project Alpha"}},
		Placeholder: "All projects",
		OnChange:    "this.closest('form').requestSubmit()",
	}).Render(t.Context(), &body)
	if err != nil {
		t.Fatal(err)
	}

	html := body.String()
	for _, want := range []string{
		`name="project"`,
		`All projects`,
		`$project_filter.value !== $project_filter._lastValue`,
		`this.closest(&#39;form&#39;).requestSubmit()`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("rendered select missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, "$project_filter.value &amp;&amp;") || strings.Contains(html, "$project_filter.value &&") {
		t.Fatalf("rendered select still has truthy guard: %s", html)
	}
	if strings.Contains(html, `$project_filter.value === &#39;project &#39;alpha&#39;`) || strings.Contains(html, `$project_filter.value = &#39;project &#39;alpha&#39;`) {
		t.Fatalf("rendered select used unsafe single-quoted option value: %s", html)
	}
}

func TestSelectWithNonNilEmptyOptionsRendersDisabledTrigger(t *testing.T) {
	var body bytes.Buffer
	err := Select(SelectArgs{
		ID:      "empty_select",
		Options: []SelectOptionArgs{},
	}).Render(t.Context(), &body)
	if err != nil {
		t.Fatal(err)
	}

	html := body.String()
	for _, want := range []string{`data-slot="select-trigger"`, `disabled`} {
		if !strings.Contains(html, want) {
			t.Fatalf("empty-options select missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, `data-slot="select-content"`) {
		t.Fatalf("empty-options select rendered content: %s", html)
	}
}
