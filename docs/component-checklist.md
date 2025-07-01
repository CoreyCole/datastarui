# Component Refactoring Checklist

Use this checklist to track progress when refactoring each component. Copy and customize for each component.

## 📋 Component: [COMPONENT_NAME]

### Pre-Refactoring Assessment

#### Current State Analysis
- [ ] Reviewed component .templ file for string concatenation patterns
- [ ] Identified complex Datastar expressions that need refactoring
- [ ] Documented current signal management approach
- [ ] Listed all interactive features (click, keyboard, hover, etc.)
- [ ] Noted any conditional CSS class patterns

#### Baseline Testing
- [ ] Development server running (`just up` + `just docker-tail app`)
- [ ] Component page accessible at `http://localhost:4242/components/[component]`
- [ ] Baseline screenshot captured: `[component]-baseline.png`
- [ ] All variants/states tested and documented
- [ ] Console logs captured (no errors): `type="all" clear=true`
- [ ] Interactive features tested manually

#### Expression Pattern Analysis
- [ ] Located string concatenation patterns (search for `fmt.Sprintf`, `+` concatenation)
- [ ] Identified complex conditional logic in Datastar expressions  
- [ ] Found manual data-class object syntax (`{'class': condition}`)
- [ ] Listed keyboard navigation handlers
- [ ] Documented signal coordination between components

### Refactoring Implementation

#### Signal Management
- [ ] Updated to use `utils.Signals(id, StructType{...})`
- [ ] Replaced manual JSON with struct definitions
- [ ] Updated signal references to use `signals.Signal("prop")`
- [ ] Replaced toggle logic with `signals.Toggle("prop")`
- [ ] Updated conditional expressions with `signals.Conditional()`

#### Expression Builders
- [ ] Replaced string concatenation with `utils.NewExpression()`
- [ ] Updated click handlers using expression builders
- [ ] Refactored keyboard navigation with `utils.NewKeyboardHandler()`
- [ ] Converted complex conditionals to builder pattern
- [ ] Added proper `PreventDefault()` where needed

#### Data-Class Helpers
- [ ] Replaced manual class objects with `utils.NewDataClass()`
- [ ] Used specialized helpers (`HighlightedItem`, `ActiveTab`, etc.)
- [ ] Updated conditional styling to use `data-class` attribute
- [ ] Verified class conditions use proper comparison operators

#### Component-Specific Patterns
- [ ] [Add component-specific refactoring tasks here]
- [ ] [e.g., Calendar: Date calculation expressions]
- [ ] [e.g., Select: Keyboard navigation with highlighting]
- [ ] [e.g., Tabs: Active state management]

### Post-Refactoring Validation

#### Compilation & Build
- [ ] Component compiles without errors (`just docker-tail app`)
- [ ] No Go compilation errors
- [ ] templ files generate successfully
- [ ] Tailwind CSS builds correctly

#### Functional Testing
- [ ] All interactive features work identically
- [ ] Click handlers function correctly
- [ ] Keyboard navigation preserves all functionality
- [ ] State management works as expected
- [ ] Multiple instances work without conflicts

#### Visual Regression Testing
- [ ] Component renders identically to baseline
- [ ] All variants display correctly
- [ ] Interactive states (hover, focus, active) unchanged
- [ ] Responsive behavior maintained
- [ ] Dark mode support preserved

#### Browser Console Validation
- [ ] No JavaScript errors in console
- [ ] No Datastar expression parsing errors
- [ ] Signal updates work correctly
- [ ] Console logs captured: `type="error"` returns empty

#### Cross-Browser Testing
- [ ] Chrome/Chromium functionality verified
- [ ] Firefox compatibility confirmed
- [ ] Safari/WebKit behavior tested (if possible)
- [ ] Mobile viewport testing completed

#### Accessibility Testing
- [ ] Keyboard navigation fully functional
- [ ] ARIA attributes preserved
- [ ] Screen reader compatibility maintained
- [ ] Focus management works correctly
- [ ] Tab order is logical

### Testing Protocol

#### Playwright MCP Commands Used
```bash
# Navigation
mcp__playwright__playwright_navigate http://localhost:4242/components/[component]

# Console monitoring
mcp__playwright__playwright_console_logs type="all" clear=true
mcp__playwright__playwright_console_logs type="error"

# Interaction testing
mcp__playwright__playwright_click "[data-selector]"
mcp__playwright__playwright_press_key "ArrowDown"
mcp__playwright__playwright_hover "[element]"

# Visual verification
mcp__playwright__playwright_screenshot name="[component]-[state]"
```

#### Specific Test Scenarios
- [ ] [List component-specific test scenarios]
- [ ] [e.g., Calendar: Month navigation, date selection, range mode]
- [ ] [e.g., Select: Dropdown open/close, option selection, keyboard nav]
- [ ] [e.g., Tabs: Tab switching, active state, keyboard navigation]

### Performance Verification

#### Expression Complexity
- [ ] New expressions are simpler than original
- [ ] No unnecessary computation in expressions
- [ ] Signal updates are efficient
- [ ] No memory leaks or excessive DOM manipulation

#### Bundle Impact
- [ ] No significant increase in generated code size
- [ ] Utility functions reused effectively
- [ ] No duplicate expression patterns

### Documentation & Cleanup

#### Code Quality
- [ ] Code is more readable than original
- [ ] Expressions are self-documenting
- [ ] Consistent patterns used throughout
- [ ] No TODO comments or debugging code left

#### Pattern Documentation
- [ ] New patterns documented in guide.md (if novel)
- [ ] Component added as example (if exemplary)
- [ ] Common pitfalls noted in debugging.md (if encountered)
- [ ] Testing scenarios added to checklist template (if reusable)

### Sign-off

#### Completion Criteria
- [ ] All checklist items completed
- [ ] No regressions identified
- [ ] Code quality improved over original
- [ ] Full test coverage achieved
- [ ] Documentation updated as needed

#### Final Validation
- [ ] Component ready for production use
- [ ] Patterns can be applied to other components
- [ ] Knowledge captured for team use
- [ ] Ready to move to next component

---

**Refactoring Completed**: [DATE]
**Time Spent**: [HOURS]
**Issues Encountered**: [NOTES]
**Patterns Learned**: [INSIGHTS]
**Next Component**: [COMPONENT_NAME]