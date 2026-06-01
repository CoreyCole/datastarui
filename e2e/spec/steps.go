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

type Key string

const (
	KeyEscape Key = "Escape"
	KeyTab    Key = "Tab"
	KeyEnter  Key = "Enter"
	KeySpace  Key = " "
)

type Attribute string

const (
	AriaExpanded Attribute = "aria-expanded"
	AriaModal    Attribute = "aria-modal"
	RoleAttr     Attribute = "role"
	ClassAttr    Attribute = "class"
	DataShow     Attribute = "data-show"
	DataText     Attribute = "data-text"
)

type pathPage string

func Path(path string) Page        { return pathPage(path) }
func (p pathPage) VisitStep() Step { return Visit(string(p)) }

type stepExpectation struct{ step Step }

func ExpectStep(step Step) Expectation    { return stepExpectation{step: step} }
func (e stepExpectation) CheckStep() Step { return e.step }

func All(steps ...Step) Step {
	return Custom("all", func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		for _, step := range steps {
			step.Run(t, ctx)
		}
	})
}

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

func Press(locator Locator, key Key) Step {
	return Custom("press "+string(key)+" on "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := resolved.First().Press(string(key)); err != nil {
			t.Fatal(err)
		}
	})
}

func PressPage(key Key) Step {
	return Custom("press page "+string(key), func(t testing.TB, ctx *runtime.Context) {
		if err := ctx.Page.Keyboard().Press(string(key)); err != nil {
			t.Fatal(err)
		}
	})
}

func InputValue(locator Locator, want string) Expectation {
	return ExpectStep(Custom("input value "+locatorName(locator)+" = "+want, func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		pollUntil(t, 2*time.Second, func() (bool, string, error) {
			got, err := resolved.First().InputValue()
			return got == want, got, err
		})
	}))
}

func TextEquals(locator Locator, want string) Expectation {
	return ExpectStep(Custom("text equals "+locatorName(locator)+" = "+want, func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		pollUntil(t, 2*time.Second, func() (bool, string, error) {
			got, err := resolved.First().InnerText()
			return strings.TrimSpace(got) == want, got, err
		})
	}))
}

func TextContains(locator Locator, want string) Expectation {
	return ExpectStep(Custom("text contains "+locatorName(locator)+" = "+want, func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		pollUntil(t, 2*time.Second, func() (bool, string, error) {
			got, err := resolved.First().InnerText()
			return strings.Contains(got, want), got, err
		})
	}))
}

func AttributeEquals(locator Locator, attr Attribute, want string) Expectation {
	return ExpectStep(Custom("attribute "+string(attr)+" equals "+locatorName(locator)+" = "+want, func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		pollUntil(t, 2*time.Second, func() (bool, string, error) {
			got, err := resolved.First().GetAttribute(string(attr))
			return got == want, got, err
		})
	}))
}

func AttributeContains(locator Locator, attr Attribute, want string) Expectation {
	return ExpectStep(Custom("attribute "+string(attr)+" contains "+locatorName(locator)+" = "+want, func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		pollUntil(t, 2*time.Second, func() (bool, string, error) {
			got, err := resolved.First().GetAttribute(string(attr))
			return strings.Contains(got, want), got, err
		})
	}))
}

func Focused(locator Locator) Expectation {
	return ExpectStep(Custom("focused "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		pollUntil(t, 2*time.Second, func() (bool, string, error) {
			result, err := resolved.First().Evaluate(`el => el === document.activeElement`, nil)
			if err != nil {
				return false, "", err
			}
			if ok, _ := result.(bool); ok {
				return true, "active element matched", nil
			}
			active, _ := ctx.Page.Locator(":focus").First().Evaluate(`el => el ? (el.id || el.getAttribute('data-slot') || el.tagName) : 'none'`, nil)
			return false, fmt.Sprint(active), nil
		})
	}))
}

func WithinViewport(locator Locator) Expectation {
	return ExpectStep(Custom("within viewport "+locatorName(locator), func(t testing.TB, ctx *runtime.Context) {
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		box, err := resolved.First().BoundingBox()
		if err != nil {
			t.Fatal(err)
		}
		viewport := ctx.Page.ViewportSize()
		if viewport == nil {
			t.Fatal("page has no viewport size")
		}
		if box == nil {
			t.Fatalf("%s has no bounding box", locatorName(locator))
		}
		if box.X < 0 || box.Y < 0 || box.X+box.Width > float64(viewport.Width) || box.Y+box.Height > float64(viewport.Height) {
			t.Fatalf("%s outside viewport: box=%+v viewport=%+v", locatorName(locator), box, viewport)
		}
	}))
}

func NotClippedBy(child Locator, ancestor Locator) Expectation {
	return ExpectStep(Custom("not clipped "+locatorName(child)+" by "+locatorName(ancestor), func(t testing.TB, ctx *runtime.Context) {
		childResolved, err := child.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		ancestorResolved, err := ancestor.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		childBox, err := childResolved.First().BoundingBox()
		if err != nil {
			t.Fatal(err)
		}
		ancestorBox, err := ancestorResolved.First().BoundingBox()
		if err != nil {
			t.Fatal(err)
		}
		if childBox == nil || ancestorBox == nil {
			t.Fatalf("missing bounding box: child=%+v ancestor=%+v", childBox, ancestorBox)
		}

		childBottom := childBox.Y + childBox.Height
		ancestorBottom := ancestorBox.Y + ancestorBox.Height
		childRight := childBox.X + childBox.Width
		ancestorRight := ancestorBox.X + ancestorBox.Width
		points := clippedCheckPoints(childBox, ancestorBox)
		if len(points) == 0 {
			t.Fatalf("%s appears fully constrained inside %s: child=%+v ancestor=%+v", locatorName(child), locatorName(ancestor), childBox, ancestorBox)
		}
		visible, err := childResolved.First().Evaluate(`(el, points) => points.some((point) => {
			const x = point.x ?? point.X;
			const y = point.y ?? point.Y;
			const hit = document.elementFromPoint(x, y);
			return hit && (hit === el || el.contains(hit));
		})`, points)
		if err != nil {
			t.Fatal(err)
		}
		if ok, _ := visible.(bool); ok {
			return
		}
		t.Fatalf("%s extends outside %s but outside points are clipped: child=%+v ancestor=%+v childBottom=%v ancestorBottom=%v childRight=%v ancestorRight=%v", locatorName(child), locatorName(ancestor), childBox, ancestorBox, childBottom, ancestorBottom, childRight, ancestorRight)
	}))
}

func clippedCheckPoints(childBox, ancestorBox *playwright.Rect) []map[string]interface{} {
	childBottom := childBox.Y + childBox.Height
	ancestorBottom := ancestorBox.Y + ancestorBox.Height
	childRight := childBox.X + childBox.Width
	ancestorRight := ancestorBox.X + ancestorBox.Width
	midX := childBox.X + childBox.Width/2
	midY := childBox.Y + childBox.Height/2
	points := []map[string]interface{}{}
	if childBottom > ancestorBottom {
		points = append(points, map[string]interface{}{"x": midX, "y": (ancestorBottom + childBottom) / 2})
	}
	if childRight > ancestorRight {
		points = append(points, map[string]interface{}{"x": (ancestorRight + childRight) / 2, "y": midY})
	}
	if childBox.Y < ancestorBox.Y {
		points = append(points, map[string]interface{}{"x": midX, "y": (childBox.Y + ancestorBox.Y) / 2})
	}
	if childBox.X < ancestorBox.X {
		points = append(points, map[string]interface{}{"x": (childBox.X + ancestorBox.X) / 2, "y": midY})
	}
	return points
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
		if ctx.Config.App.Authenticate == nil {
			return
		}
		if err := ctx.Config.App.Authenticate(context.Background(), ctx.Page, ctx.Config, email); err != nil {
			t.Fatal(err)
		}
	})
}

func pollUntil(t testing.TB, timeout time.Duration, check func() (bool, string, error)) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	last := ""
	for time.Now().Before(deadline) {
		ok, got, err := check()
		if err != nil {
			last = err.Error()
		} else if ok {
			return
		} else {
			last = got
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("condition not met before timeout; last value %q", last)
}

func locatorName(locator Locator) string {
	if stringer, ok := locator.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%T", locator)
}
