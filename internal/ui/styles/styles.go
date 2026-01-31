// Package styles provides consistent styling for UI components across terminal themes.
// This ensures good visibility and contrast regardless of the user's terminal theme.
package styles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InputFieldColors defines colors for input field components
type InputFieldColors struct {
	Label           tcell.Color
	Field           tcell.Color
	FieldBackground tcell.Color
	Placeholder     tcell.Color
	Autocomplete    tcell.Color
}

// DefaultInputColors returns colors that work well across most terminal themes.
// Uses high-contrast combinations that are visible on both light and dark themes.
func DefaultInputColors() InputFieldColors {
	return InputFieldColors{
		Label:           tcell.ColorYellow,
		Field:           tcell.ColorWhite,
		FieldBackground: tcell.ColorDarkBlue,
		Placeholder:     tcell.ColorGray,
		Autocomplete:    tcell.ColorGreen,
	}
}

// StyleInputField applies consistent styling to a tview InputField.
// This ensures the input field is visible across different terminal themes.
func StyleInputField(input *tview.InputField) *tview.InputField {
	colors := DefaultInputColors()
	return input.
		SetLabelColor(colors.Label).
		SetFieldTextColor(colors.Field).
		SetFieldBackgroundColor(colors.FieldBackground).
		SetPlaceholderTextColor(colors.Placeholder).
		SetAutocompleteStyles(
			colors.FieldBackground,
			tcell.StyleDefault.Foreground(colors.Field).Background(colors.FieldBackground),
			tcell.StyleDefault.Foreground(colors.FieldBackground).Background(colors.Autocomplete),
		)
}

// StyleInputFieldWithColors applies custom colors to an input field.
func StyleInputFieldWithColors(input *tview.InputField, colors InputFieldColors) *tview.InputField {
	return input.
		SetLabelColor(colors.Label).
		SetFieldTextColor(colors.Field).
		SetFieldBackgroundColor(colors.FieldBackground).
		SetPlaceholderTextColor(colors.Placeholder)
}

// NewStyledInputField creates a new input field with consistent styling applied.
func NewStyledInputField() *tview.InputField {
	return StyleInputField(tview.NewInputField())
}

// FormColors defines colors for form components
type FormColors struct {
	Label           tcell.Color
	Field           tcell.Color
	FieldBackground tcell.Color
	Button          tcell.Color
	ButtonActive    tcell.Color
}

// DefaultFormColors returns colors that work well for forms.
func DefaultFormColors() FormColors {
	return FormColors{
		Label:           tcell.ColorYellow,
		Field:           tcell.ColorWhite,
		FieldBackground: tcell.ColorDarkBlue,
		Button:          tcell.ColorWhite,
		ButtonActive:    tcell.ColorGreen,
	}
}

// StyleForm applies consistent styling to a tview Form.
func StyleForm(form *tview.Form) *tview.Form {
	colors := DefaultFormColors()
	return form.
		SetLabelColor(colors.Label).
		SetFieldTextColor(colors.Field).
		SetFieldBackgroundColor(colors.FieldBackground).
		SetButtonTextColor(colors.Button).
		SetButtonBackgroundColor(tcell.ColorDarkBlue).
		SetButtonActivatedStyle(tcell.StyleDefault.
			Foreground(tcell.ColorBlack).
			Background(colors.ButtonActive))
}

// DropDownColors defines colors for dropdown components
type DropDownColors struct {
	Label           tcell.Color
	Field           tcell.Color
	FieldBackground tcell.Color
	Options         tcell.Color
	OptionsSelected tcell.Color
}

// DefaultDropDownColors returns colors that work well for dropdowns.
func DefaultDropDownColors() DropDownColors {
	return DropDownColors{
		Label:           tcell.ColorYellow,
		Field:           tcell.ColorWhite,
		FieldBackground: tcell.ColorDarkBlue,
		Options:         tcell.ColorWhite,
		OptionsSelected: tcell.ColorGreen,
	}
}

// StyleDropDown applies consistent styling to a tview DropDown.
func StyleDropDown(dropdown *tview.DropDown) *tview.DropDown {
	colors := DefaultDropDownColors()
	return dropdown.
		SetLabelColor(colors.Label).
		SetFieldTextColor(colors.Field).
		SetFieldBackgroundColor(colors.FieldBackground).
		SetListStyles(
			tcell.StyleDefault.Foreground(colors.Options).Background(colors.FieldBackground),
			tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(colors.OptionsSelected),
		)
}

// TextAreaColors defines colors for text area components
type TextAreaColors struct {
	Text            tcell.Color
	Background      tcell.Color
	Placeholder     tcell.Color
}

// DefaultTextAreaColors returns colors that work well for text areas.
func DefaultTextAreaColors() TextAreaColors {
	return TextAreaColors{
		Text:        tcell.ColorWhite,
		Background:  tcell.ColorDarkBlue,
		Placeholder: tcell.ColorGray,
	}
}

// StyleTextArea applies consistent styling to a tview TextArea.
func StyleTextArea(textarea *tview.TextArea) *tview.TextArea {
	colors := DefaultTextAreaColors()
	return textarea.
		SetTextStyle(tcell.StyleDefault.Foreground(colors.Text).Background(colors.Background)).
		SetPlaceholderStyle(tcell.StyleDefault.Foreground(colors.Placeholder).Background(colors.Background))
}
