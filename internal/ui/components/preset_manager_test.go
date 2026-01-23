package components

import (
	"testing"

	"github.com/jontk/s9s/internal/ui/filters"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPresetManagerUI(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	viewType := "jobs"

	pm := NewPresetManagerUI(app, presetManager, viewType)

	require.NotNil(t, pm)
	assert.Equal(t, app, pm.app)
	assert.Equal(t, presetManager, pm.presetManager)
	assert.Equal(t, viewType, pm.viewType)
	assert.Nil(t, pm.pages)
	assert.Nil(t, pm.onDone)
}

func TestPresetManagerUIShow(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	onDoneCalled := false

	pm.Show(pages, func() {
		onDoneCalled = true
	})

	assert.Equal(t, pages, pm.pages)
	assert.NotNil(t, pm.onDone)

	// Trigger the callback
	pm.onDone()
	assert.True(t, onDoneCalled)
}

func TestPresetManagerUIClose(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	onDoneCalled := false

	pm.Show(pages, func() {
		onDoneCalled = true
	})

	pm.close()

	assert.True(t, onDoneCalled)
}

func TestPresetManagerUICloseWithoutCallback(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	// Should not panic when onDone is nil
	pm.close()
}

func TestPresetManagerUICreateCenteredModal(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	content := tview.NewTextView()
	modal := pm.createCenteredModal(content, 60, 20)

	require.NotNil(t, modal)
	// Modal should be a flex container
	_, ok := modal.(*tview.Flex)
	assert.True(t, ok, "createCenteredModal should return a Flex container")
}

func TestPresetManagerUIShowError(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	// Should create an error modal without panicking
	pm.showError("Test error message")

	// Verify page count increased
	assert.Equal(t, 1, pages.GetPageCount())
}

func TestPresetManagerUIShowSuccess(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	// Should create a success modal without panicking
	pm.showSuccess("Test success message")

	// Verify page count increased
	assert.Equal(t, 1, pages.GetPageCount())
}

func TestPresetManagerUIShowInfo(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	// Should create an info modal without panicking
	pm.showInfo("Test info message")

	// Verify page count increased
	assert.Equal(t, 1, pages.GetPageCount())
}

func TestPresetManagerUITestPreset(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	preset := filters.FilterPreset{
		Name:      "Test Preset",
		FilterStr: "state=running",
	}

	// Should not panic when testing preset
	pm.testPreset(preset)
}

func TestPresetManagerUIShowImportDialog(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	// Should not panic
	pm.showImportDialog()
}

func TestPresetManagerUIShowExportDialog(t *testing.T) {
	app := tview.NewApplication()
	presetManager := filters.NewPresetManager()
	pm := NewPresetManagerUI(app, presetManager, "jobs")

	pages := tview.NewPages()
	pm.pages = pages

	// Should not panic
	pm.showExportDialog()
}
