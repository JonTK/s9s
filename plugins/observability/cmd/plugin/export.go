package main

import (
	"context"
	"fmt"
	
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	
	"github.com/jontk/s9s/internal/dao"
	"github.com/jontk/s9s/internal/plugin"
	"github.com/jontk/s9s/internal/plugins"
	observability "github.com/jontk/s9s/plugins/observability"
)

// PluginAdapter adapts the observability plugin to the s9s plugin interface
type PluginAdapter struct {
	plugin *observability.ObservabilityPlugin
}

// NewPlugin is the exported function that s9s looks for when loading plugins
func NewPlugin() plugins.Plugin {
	return &PluginAdapter{
		plugin: observability.New(),
	}
}

// GetInfo returns plugin information
func (a *PluginAdapter) GetInfo() plugins.PluginInfo {
	info := a.plugin.GetInfo()
	return plugins.PluginInfo{
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Website:     info.License, // Using License field for website since it's available
	}
}

// Initialize initializes the plugin
func (a *PluginAdapter) Initialize(ctx context.Context, client dao.SlurmClient) error {
	// Convert the config from context if available
	configMap := make(map[string]interface{})
	
	// The s9s plugin manager might pass config through context
	if config, ok := ctx.Value("plugin_config").(map[string]interface{}); ok {
		configMap = config
	}
	
	// Log debug info about config
	fmt.Printf("[DEBUG] Plugin config map: %+v\n", configMap)
	
	// Initialize the observability plugin
	if err := a.plugin.Init(ctx, configMap); err != nil {
		fmt.Printf("[ERROR] Plugin Init failed: %v\n", err)
		return err
	}
	
	// Start the plugin
	if err := a.plugin.Start(ctx); err != nil {
		fmt.Printf("[ERROR] Plugin Start failed: %v\n", err)
		return err
	}
	
	return nil
}

// GetCommands returns the commands this plugin provides
func (a *PluginAdapter) GetCommands() []plugins.Command {
	// The observability plugin doesn't provide commands currently
	return []plugins.Command{}
}

// GetViews returns the views this plugin provides
func (a *PluginAdapter) GetViews() []plugins.View {
	views := a.plugin.GetViews()
	result := make([]plugins.View, 0, len(views))
	
	for _, v := range views {
		// Create view adapter
		viewAdapter := &ViewAdapter{
			plugin: a.plugin,
			info:   v,
		}
		result = append(result, viewAdapter)
	}
	
	return result
}

// GetKeyBindings returns custom key bindings
func (a *PluginAdapter) GetKeyBindings() []plugins.KeyBinding {
	// The observability plugin doesn't provide custom key bindings currently
	return []plugins.KeyBinding{}
}

// OnEvent handles events
func (a *PluginAdapter) OnEvent(event plugins.Event) error {
	// The observability plugin doesn't handle events currently
	return nil
}

// Cleanup cleans up the plugin
func (a *PluginAdapter) Cleanup() error {
	ctx := context.Background()
	return a.plugin.Stop(ctx)
}

// ViewAdapter adapts the observability view to the s9s view interface
type ViewAdapter struct {
	plugin *observability.ObservabilityPlugin
	info   plugin.ViewInfo
	view   plugin.View
}

// GetName returns the view name
func (v *ViewAdapter) GetName() string {
	return v.info.ID
}

// GetTitle returns the view title
func (v *ViewAdapter) GetTitle() string {
	return v.info.Name
}

// Render returns the tview primitive
func (v *ViewAdapter) Render() tview.Primitive {
	if v.view != nil {
		return v.view.GetPrimitive()
	}
	return nil
}

// OnKey handles key events
func (v *ViewAdapter) OnKey(event *tcell.EventKey) *tcell.EventKey {
	if v.view != nil {
		handled := v.view.HandleKey(event)
		if handled {
			return nil // Event was handled
		}
	}
	return event
}

// Refresh updates the view data
func (v *ViewAdapter) Refresh() error {
	if v.view != nil {
		ctx := context.Background()
		return v.view.Update(ctx)
	}
	return nil
}

// Init initializes the view
func (v *ViewAdapter) Init(ctx context.Context) error {
	// Create the actual view
	view, err := v.plugin.CreateView(ctx, v.info.ID)
	if err != nil {
		return err
	}
	v.view = view
	return nil
}