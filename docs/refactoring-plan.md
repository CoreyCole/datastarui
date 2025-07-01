# DatastarUI Component Refactoring Plan

This document outlines a systematic approach to refactor all DatastarUI components to use the new utility libraries for cleaner, more maintainable Datastar integration.

## 🎯 Objectives

1. **Eliminate string concatenation** in Datastar expressions
2. **Standardize patterns** using utility libraries
3. **Improve maintainability** and reduce errors
4. **Ensure full functionality** with comprehensive testing
5. **Document patterns** for future development

## 📊 Component Assessment

### 🔴 High Priority - Complex Refactoring Needed

#### 1. Calendar Component
- **Complexity**: Extremely high
- **Issues**: Complex date calculations, range selection, multi-month support
- **Expected Impact**: Major improvement in maintainability
- **Estimated Effort**: 3-4 sessions

#### 2. DatePicker Component  
- **Complexity**: High
- **Issues**: Signal coordination between DateInput and Calendar
- **Expected Impact**: Simplified signal orchestration
- **Estimated Effort**: 2-3 sessions

#### 3. DateInput Component
- **Complexity**: High
- **Issues**: Auto-formatting, calendar coordination, date parsing
- **Expected Impact**: Cleaner input handling logic
- **Estimated Effort**: 2-3 sessions

### 🟡 Medium Priority - Moderate Improvements

#### 4. Dropdown Component
- **Complexity**: Medium
- **Issues**: Basic keyboard navigation, positioning
- **Expected Impact**: Standardized keyboard patterns
- **Estimated Effort**: 1-2 sessions

#### 5. Tabs Component
- **Complexity**: Medium  
- **Issues**: Manual data-class objects for active styling
- **Expected Impact**: Cleaner conditional classes
- **Estimated Effort**: 1 session

### 🟢 Low Priority - Complete Existing Work

#### 6. Select Component
- **Status**: Partially refactored (80% complete)
- **Issues**: Minor cleanup needed
- **Expected Impact**: Template completion
- **Estimated Effort**: 1 session

#### 7. Dialog Component
- **Status**: Mostly good (already using some utilities)
- **Issues**: Minor keyboard handler improvements
- **Expected Impact**: Consistency improvements
- **Estimated Effort**: 0.5 sessions

## 🛠️ Refactoring Methodology

### Phase 1: Preparation
1. **Create test baseline** using Playwright MCP
2. **Document current behavior** with screenshots
3. **Identify expression patterns** to refactor
4. **Plan utility usage** for each component

### Phase 2: Implementation  
1. **Refactor expressions** using utility builders
2. **Update signal management** with structured patterns
3. **Implement data-class helpers** for conditional styling
4. **Test incrementally** after each change

### Phase 3: Validation
1. **Comprehensive testing** with Playwright MCP
2. **Visual regression testing** with screenshots
3. **Keyboard navigation testing** for accessibility
4. **Cross-browser validation** (Chrome, Firefox, Safari)

### Phase 4: Documentation
1. **Update component examples** in guide.md
2. **Document new patterns** discovered
3. **Add testing scenarios** to debugging.md

## 📋 Standard Testing Protocol

For each component refactoring:

### Pre-Refactoring Baseline
```bash
# 1. Start development environment
just up
just docker-tail app  # Verify compilation

# 2. Navigate to component page
mcp__playwright__playwright_navigate http://localhost:4242/components/[component]

# 3. Capture baseline behavior
mcp__playwright__playwright_console_logs type="all" clear=true
mcp__playwright__playwright_screenshot name="[component]-baseline"

# 4. Test all interactive features
# - Click all buttons/triggers
# - Test keyboard navigation  
# - Test different states/variants
# - Capture screenshots for each state

# 5. Document current expressions
# Review .templ file for string concatenation patterns
```

### Post-Refactoring Validation
```bash
# 1. Verify compilation
just docker-tail app

# 2. Test identical functionality  
mcp__playwright__playwright_navigate http://localhost:4242/components/[component]
mcp__playwright__playwright_console_logs type="all" clear=true

# 3. Compare behavior
mcp__playwright__playwright_screenshot name="[component]-refactored"

# 4. Test edge cases
# - Error conditions
# - Boundary values  
# - Keyboard accessibility
# - Multiple instances

# 5. Validate no console errors
mcp__playwright__playwright_console_logs type="error"
```

## 🗓️ Implementation Schedule

### Week 1: High Priority Components
- **Day 1-2**: Calendar Component
  - Analyze complex date expressions
  - Refactor using expression builders
  - Test range selection and multi-month
  
- **Day 3**: DatePicker Component  
  - Signal coordination refactoring
  - Test single/range modes
  
- **Day 4**: DateInput Component
  - Auto-formatting logic
  - Calendar integration testing

### Week 2: Medium Priority & Completion
- **Day 1**: Dropdown Component
  - Keyboard navigation standardization
  - Test positioning and states
  
- **Day 2**: Tabs Component
  - data-class helper implementation
  - Active state testing
  
- **Day 3**: Select Component (completion)
  - Finish partial refactoring
  - Complete testing coverage
  
- **Day 4**: Dialog Component  
  - Minor improvements
  - Consistency updates

### Week 3: Testing & Documentation
- **Day 1-2**: Comprehensive cross-component testing
- **Day 3**: Documentation updates
- **Day 4**: Final validation and cleanup

## 📖 Refactoring Templates

### Expression Builder Pattern
```go
// Before: String concatenation
clickHandler := fmt.Sprintf("$%s.open = !$%s.open; $%s.highlighted = 0", 
    id, id, id)

// After: Expression builder
clickHandler := utils.NewExpression().
    Statement(signals.Toggle("open")).
    SetSignal("highlighted", "0").
    Build()
```

### Data-Class Pattern
```go
// Before: Manual object syntax
dataClass := fmt.Sprintf("{'active': $%s.tab === '%s'}", id, value)

// After: Data-class helper
dataClass := utils.ActiveTab(fmt.Sprintf("%s.tab", id), value)
```

### Keyboard Navigation Pattern
```go
// Before: Complex conditional
keyHandler := fmt.Sprintf("evt.key === 'Enter' ? (%s) : evt.key === 'Escape' ? (%s) : null", 
    enterAction, escapeAction)

// After: Keyboard builder
keyHandler := utils.NewKeyboardHandler("Enter", "Escape").
    OnKeys(func(expr *DatastarExpression) *DatastarExpression {
        return expr.
            PreventDefault().
            Conditional("evt.key === 'Enter'", enterAction, escapeAction)
    }).
    Build()
```

## 🎯 Success Criteria

For each component refactoring:

### Functional Requirements
- [ ] All existing functionality preserved
- [ ] No JavaScript console errors
- [ ] Keyboard navigation works identically
- [ ] Visual appearance unchanged
- [ ] All variants/states functional

### Code Quality Requirements  
- [ ] No string concatenation for expressions
- [ ] All complex expressions use utility builders
- [ ] Conditional classes use data-class helpers
- [ ] Signal management uses utils.Signals()
- [ ] Code is readable and maintainable

### Testing Requirements
- [ ] Playwright MCP tests pass
- [ ] Screenshot comparison shows no regressions
- [ ] Cross-browser compatibility maintained
- [ ] Accessibility features preserved
- [ ] Performance characteristics unchanged

## 🚀 Getting Started

To begin the refactoring process:

1. **Review this plan** and select first component
2. **Follow testing protocol** to establish baseline
3. **Implement refactoring** using utility patterns
4. **Validate thoroughly** with Playwright MCP
5. **Document learnings** for future components

The Calendar component is recommended as the first target due to its complexity and potential for significant improvement.