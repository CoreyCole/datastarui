# DatePicker Component

A comprehensive date picker component that combines an input field with an auto-opening calendar popover for intuitive date selection.

## Features

- 🗓️ **Single & Range Selection**: Support for both single date and date range selection
- ⌨️ **Smart Input**: Users can type dates with automatic slash insertion (20241225 → 2024/12/25)
- 🚀 **Auto-Opening Popover**: Calendar opens automatically on focus and typing
- 🔄 **Real-Time Sync**: Calendar automatically updates as user types
- 🎨 **Customizable Styling**: Multiple styling options for different use cases
- 📝 **Form Integration**: Seamless integration with forms using data-bind
- ✅ **Validation**: Built-in validation with min/max dates and required fields
- ♿ **Accessible**: Full keyboard navigation and screen reader support
- 📱 **Responsive**: Works on desktop and mobile devices

## Basic Usage

### Single Date Selection

```go
@datepicker.DatePicker(datepicker.DatePickerArgs{
    ID:          "single-date",
    Name:        "selectedDate",
    Mode:        "single",
    Placeholder: "YYYY/MM/DD",
})
```

### Date Range Selection

```go
@datepicker.DatePicker(datepicker.DatePickerArgs{
    ID:             "date-range",
    Name:           "dateRange",
    Mode:           "range",
    Placeholder:    "Select date range...",
    NumberOfMonths: 2,
})
```

## Auto-Opening Behavior

The DatePicker automatically opens the calendar popover when:

- ✅ **Input Focus**: User clicks or tabs into the input field
- ✅ **Start Typing**: User begins typing a date
- ✅ **Icon Click**: User clicks the calendar icon

### Real-Time Calendar Sync

As the user types, the calendar automatically:

- 📅 **Updates Month**: Navigates to the month being typed
- 🎯 **Highlights Date**: Shows the date being entered (if valid)
- ✨ **Live Preview**: Provides immediate visual feedback

## Args Reference

### Core Args

| Prop | Type | Default | Description |
| ------------- | -------- | -------------- | --------------------------------------- |
| `ID` | `string` | auto-generated | Unique identifier for the datepicker |
| `Name` | `string` | - | Name attribute for form submission |
| `Mode` | `string` | `"single"` | Selection mode: `"single"` or `"range"` |
| `Placeholder` | `string` | - | Placeholder text for the input field |

### Date Args

| Prop | Type | Default | Description |
| --------------- | -------- | ------------- | ------------------------------------- |
| `DefaultDate` | `string` | current month | Initial month to display (YYYY-MM-DD) |
| `SelectedDate` | `string` | - | Pre-selected date for single mode |
| `RangeStart` | `string` | - | Start date for range mode |
| `RangeEnd` | `string` | - | End date for range mode |
| `MinDate` | `string` | - | Minimum selectable date |
| `MaxDate` | `string` | - | Maximum selectable date |
| `DisabledDates` | `string` | - | Comma-separated disabled dates |

### Display Args

| Prop | Type | Default | Description |
| ----------------- | ------ | ------- | --------------------------------------- |
| `NumberOfMonths` | `int` | `1` | Number of months to display in calendar |
| `HideOutsideDays` | `bool` | `false` | Hide dates from adjacent months |
| `ShowIcon` | `bool` | `true` | Show calendar icon in input |

### Behavior Args

| Prop | Type | Default | Description |
| --------------- | ------ | ------- | ------------------------------------------- |
| `AutoSlash` | `bool` | `true` | Auto-insert slashes (20241225 → 2024/12/25) |
| `OpenOnFocus` | `bool` | `true` | Open popover when input receives focus |
| `OpenOnType` | `bool` | `true` | Open popover when user starts typing |
| `CloseOnSelect` | `bool` | varies | Close popover after selection |

### State Args

| Prop | Type | Default | Description |
| ---------- | ------ | ------- | ---------------------------------- |
| `Required` | `bool` | `false` | Required field for form validation |
| `Disabled` | `bool` | `false` | Disable the entire component |

### Styling Args

| Prop | Type | Default | Description |
| --------------- | -------- | ------- | ------------------------------------ |
| `Class` | `string` | - | Additional CSS classes for container |
| `InputClass` | `string` | - | Additional CSS classes for input |
| `CalendarClass` | `string` | - | Additional CSS classes for calendar |

### Form Args

| Prop | Type | Default | Description |
| ------------ | ------------------ | ------- | --------------------------------- |
| `FormID` | `string` | - | Form ID for data-bind integration |
| `Attributes` | `templ.Attributes` | - | Additional HTML attributes |

## Usage Examples

### Form Integration

```go
@form.Form(form.FormArgs{ID: "booking-form"}) {
    @datepicker.DatePicker(datepicker.DatePickerArgs{
        ID:          "checkin-date",
        Name:        "checkinDate",
        FormID:      "booking-form",
        Mode:        "single",
        Required:    true,
        MinDate:     "2024-01-01",
    })
}
```

### Custom Styling

```go
@datepicker.DatePicker(datepicker.DatePickerArgs{
    ID:         "styled-date",
    InputClass: "h-12 text-lg border-2 border-blue-500",
    CalendarClass: "shadow-2xl",
})
```

### Range with Validation

```go
@datepicker.DatePicker(datepicker.DatePickerArgs{
    Mode:           "range",
    NumberOfMonths: 2,
    MinDate:       "2024-01-01",
    MaxDate:       "2024-12-31",
    Required:      true,
})
```

### Customized Popover Behavior

```go
@datepicker.DatePicker(datepicker.DatePickerArgs{
    ID:            "custom-behavior",
    OpenOnFocus:   false,  // Only open on typing or icon click
    OpenOnType:    true,   // Open when user starts typing
    CloseOnSelect: false,  // Keep open after selection
})
```

### Auto-Slash Input Demo

```go
@datepicker.DatePicker(datepicker.DatePickerArgs{
    ID:          "auto-slash-demo",
    Placeholder: "Type: 20241225 → 2024/12/25",
    AutoSlash:   true, // Default behavior
})
```

## User Experience Flow

### Typical User Interaction

1. **Focus Input** → Calendar popover opens automatically
1. **Start Typing** → Calendar navigates to relevant month
1. **Continue Typing** → Calendar highlights the date being entered
1. **Complete Date** → Calendar shows selected date
1. **Click Calendar** → Alternative selection method
1. **Select Date** → Input updates, popover closes (if `CloseOnSelect: true`)

### Real-Time Feedback

| User Input | Input Display | Calendar Action |
| ------------ | ------------- | ---------------------- |
| Focus | - | Opens to current month |
| `2` | `2` | Opens to current month |
| `20` | `20` | Opens to current month |
| `2024` | `2024` | Navigates to 2024 |
| `2024/1` | `2024/1` | Navigates to Jan 2024 |
| `2024/12` | `2024/12` | Navigates to Dec 2024 |
| `2024/12/2` | `2024/12/2` | Highlights Dec 2nd |
| `2024/12/25` | `2024/12/25` | Highlights Dec 25th |

## Accessibility

The DatePicker component is fully accessible:

- ✅ Keyboard navigation (Tab, Enter, Escape, Arrow keys)
- ✅ Screen reader support with proper ARIA labels
- ✅ Focus management and visual indicators
- ✅ High contrast mode support
- ✅ Auto-announcements for date changes

## Browser Support

- ✅ Chrome 90+
- ✅ Firefox 88+
- ✅ Safari 14+
- ✅ Edge 90+

## Notes

- **Format Consistency**: Display uses `YYYY/MM/DD`, internal storage uses `YYYY-MM-DD`
- **Auto-Slash Magic**: Typing `20241225` automatically becomes `2024/12/25`
- **Smart Opening**: Popover opens on focus and typing for seamless UX
- **Real-Time Sync**: Calendar updates live as user types valid dates
- **Range Selection**: Requires two clicks (start and end dates) or typing both dates
- **Form Integration**: Uses Datastar's data-bind for reactive state management

## Related Components

- **Calendar**: The underlying calendar component
- **Input**: Base input field component
- **Popover**: Popover positioning system
- **Form**: Form state management
