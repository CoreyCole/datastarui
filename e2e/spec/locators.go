package spec

import (
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"

	"github.com/coreycole/datastarui/e2e/runtime"
)

type Locator interface {
	Resolve(ctx *runtime.Context) (playwright.Locator, error)
}

type locatorFunc struct {
	label string
	fn    func(*runtime.Context) (playwright.Locator, error)
}

func (l locatorFunc) Resolve(ctx *runtime.Context) (playwright.Locator, error) { return l.fn(ctx) }
func (l locatorFunc) String() string                                           { return l.label }

func CSS(selector string) Locator {
	return locatorFunc{label: "css:" + selector, fn: func(ctx *runtime.Context) (playwright.Locator, error) {
		if strings.TrimSpace(selector) == "" {
			return nil, fmt.Errorf("empty CSS selector")
		}
		return ctx.Page.Locator(selector), nil
	}}
}

func TestID(id string) Locator {
	selector := fmt.Sprintf("[data-testid=%q], [data-test-id=%q], [data-e2e=%q], [data-slot=%q]", id, id, id, id)
	return CSS(selector)
}

func SelectorAlias(key string) Locator {
	return locatorFunc{label: "selector:" + key, fn: func(ctx *runtime.Context) (playwright.Locator, error) {
		entry, err := ctx.Config.Selectors.Resolve(key)
		if err != nil {
			return nil, err
		}
		return ctx.Page.Locator(entry.CSS), nil
	}}
}

func Text(text string) Locator {
	return locatorFunc{label: "text:" + text, fn: func(ctx *runtime.Context) (playwright.Locator, error) {
		return ctx.Page.GetByText(text).First(), nil
	}}
}

func Label(text string) Locator {
	return locatorFunc{label: "label:" + text, fn: func(ctx *runtime.Context) (playwright.Locator, error) {
		return ctx.Page.GetByLabel(text).First(), nil
	}}
}

func Role(role, name string) Locator {
	return locatorFunc{label: "role:" + role + ":" + name, fn: func(ctx *runtime.Context) (playwright.Locator, error) {
		return ctx.Page.GetByRole(playwright.AriaRole(role), playwright.PageGetByRoleOptions{Name: name}).First(), nil
	}}
}
