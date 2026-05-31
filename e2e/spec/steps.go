package spec

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/coreycole/datastarui/e2e/runtime"
)

func OpenPage(keyOrPath string) Step {
	return Custom("open page "+keyOrPath, func(t testing.TB, ctx *runtime.Context) {
		path, err := ctx.Config.PagePath(keyOrPath)
		if err != nil {
			t.Fatal(err)
		}
		Visit(path).Run(t, ctx)
	})
}

func Visit(path string) Step {
	return Custom("visit "+path, func(t testing.TB, ctx *runtime.Context) {
		_, err := ctx.Page.Goto(ctx.Config.BaseURL+path, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func Click(locator Locator) Step {
	return Custom("click "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := resolved.First().Click(); err != nil {
			t.Fatal(err)
		}
	})
}

func Fill(locator Locator, value string) Step {
	return Custom("fill "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := resolved.First().Fill(value); err != nil {
			t.Fatal(err)
		}
	})
}

func SelectOption(locator Locator, value string) Step {
	return Custom("select "+value+" in "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		values := []string{value}
		if _, err := resolved.First().SelectOption(playwright.SelectOptionValues{Values: &values}); err != nil {
			t.Fatal(err)
		}
	})
}

func Visible(locator Locator) Step {
	return Custom("visible "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := resolved.First().WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
			t.Fatal(err)
		}
	})
}

func Hidden(locator Locator) Step {
	return Custom("hidden "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := resolved.First().WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateHidden}); err != nil {
			t.Fatal(err)
		}
	})
}

func WaitFor(locator Locator) Step { return Visible(locator) }

func TextVisible(text string) Step { return Visible(Text(text)) }

func TextAbsent(text string) Step {
	return Custom("text absent "+text, func(t testing.TB, ctx *runtime.Context) {
		body, err := ctx.Page.Locator("body").InnerText()
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(body, text) {
			t.Fatalf("expected text %q to be absent", text)
		}
	})
}

func URLContains(text string) Step {
	return Custom("url contains "+text, func(t testing.TB, ctx *runtime.Context) {
		if !strings.Contains(ctx.Page.URL(), text) {
			t.Fatalf("expected URL %q to contain %q", ctx.Page.URL(), text)
		}
	})
}

func ConsoleClean() Step {
	return Custom("console clean", func(t testing.TB, ctx *runtime.Context) {
		time.Sleep(250 * time.Millisecond)
		if problems := ctx.Console.Problems(); len(problems) > 0 {
			t.Fatalf("console problems:\n%s", runtime.FormatConsoleProblems(problems))
		}
	})
}

func AuthenticatedAs(email string) Step {
	return Custom("authenticated as "+email, func(t testing.TB, ctx *runtime.Context) {
		if err := ctx.Config.App.Authenticate(context.Background(), ctx.Page, ctx.Config, email); err != nil {
			t.Fatal(err)
		}
	})
}

func locatorName(locator Locator) string {
	if stringer, ok := locator.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%T", locator)
}
