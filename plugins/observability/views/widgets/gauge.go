package widgets

import (
	"fmt"
	"math"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GaugeWidget displays a value as a horizontal gauge/progress bar
type GaugeWidget struct {
	*tview.Box

	title string
	value float64
	min   float64
	max   float64
	unit  string

	// Display options
	showValue      bool
	showPercentage bool
	colorFunc      func(float64) tcell.Color
}

// NewGaugeWidget creates a new gauge widget
func NewGaugeWidget(title string, minVal, maxVal float64, unit string) *GaugeWidget {
	g := &GaugeWidget{
		Box:            tview.NewBox(),
		title:          title,
		min:            minVal,
		max:            maxVal,
		unit:           unit,
		showValue:      true,
		showPercentage: true,
		colorFunc:      defaultColorFunc,
	}

	g.SetBorder(true).SetTitle(title)
	return g
}

// SetValue updates the gauge value
func (g *GaugeWidget) SetValue(value float64) {
	g.value = math.Max(g.min, math.Min(g.max, value))
}

// SetColorFunc sets a custom color function
func (g *GaugeWidget) SetColorFunc(f func(float64) tcell.Color) {
	g.colorFunc = f
}

// GetPrimitive returns the primitive for this widget
func (g *GaugeWidget) GetPrimitive() tview.Primitive {
	return g
}

// Draw draws the gauge
func (g *GaugeWidget) Draw(screen tcell.Screen) {
	g.DrawForSubclass(screen, g)

	x, y, width, height := g.GetInnerRect()
	if width <= 0 || height <= 0 {
		return
	}

	percentage := g.calculatePercentage()
	gaugeY := g.calculateGaugeY(y, height)
	filledWidth := int(float64(width) * percentage)
	color := g.colorFunc(g.value)

	g.drawGaugeBar(screen, x, gaugeY, width, filledWidth, color)
	g.drawGaugeText(screen, x, gaugeY, width, filledWidth, color)
}

func (g *GaugeWidget) calculatePercentage() float64 {
	percentage := (g.value - g.min) / (g.max - g.min)
	if math.IsNaN(percentage) || math.IsInf(percentage, 0) {
		percentage = 0
	}
	return math.Max(0, math.Min(1, percentage))
}

func (g *GaugeWidget) calculateGaugeY(y, height int) int {
	if height > 1 {
		return y + (height-1)/2
	}
	return y
}

func (g *GaugeWidget) drawGaugeBar(screen tcell.Screen, x, gaugeY, width, filledWidth int, color tcell.Color) {
	for i := 0; i < width; i++ {
		if i < filledWidth {
			style := g.getFilledStyle(color)
			screen.SetContent(x+i, gaugeY, '█', nil, style)
		} else {
			screen.SetContent(x+i, gaugeY, '░', nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
		}
	}
}

func (g *GaugeWidget) getFilledStyle(color tcell.Color) tcell.Style {
	style := tcell.StyleDefault.Background(color).Foreground(tcell.ColorBlack)
	if color == tcell.ColorBlack {
		style = style.Foreground(tcell.ColorWhite)
	}
	return style
}

func (g *GaugeWidget) drawGaugeText(screen tcell.Screen, x, gaugeY, width, filledWidth int, color tcell.Color) {
	valueText := g.formatValue()
	if len(valueText) >= width {
		return
	}

	textX := x + (width-len(valueText))/2
	for i, ch := range valueText {
		style := g.getTextStyle(textX+i, x+filledWidth, color)
		screen.SetContent(textX+i, gaugeY, ch, nil, style)
	}
}

func (g *GaugeWidget) getTextStyle(textPos, filledEnd int, color tcell.Color) tcell.Style {
	style := tcell.StyleDefault.Bold(true)
	if textPos < filledEnd {
		style = style.Background(color).Foreground(tcell.ColorBlack)
		if color == tcell.ColorBlack {
			style = style.Foreground(tcell.ColorWhite)
		}
	}
	return style
}

// formatValue formats the current value for display
func (g *GaugeWidget) formatValue() string {
	parts := []string{}

	if g.showValue {
		parts = append(parts, fmt.Sprintf("%.1f%s", g.value, g.unit))
	}

	if g.showPercentage {
		percentage := (g.value - g.min) / (g.max - g.min) * 100
		parts = append(parts, fmt.Sprintf("%.1f%%", percentage))
	}

	return strings.Join(parts, " ")
}

// defaultColorFunc returns colors based on percentage thresholds
func defaultColorFunc(value float64) tcell.Color {
	switch {
	case value >= 90:
		return tcell.ColorRed
	case value >= 75:
		return tcell.ColorYellow
	case value >= 50:
		return tcell.ColorOrange
	default:
		return tcell.ColorGreen
	}
}

// GaugeGroup manages multiple gauges with consistent sizing
type GaugeGroup struct {
	*tview.Flex
	gauges []*GaugeWidget
}

// NewGaugeGroup creates a new gauge group
func NewGaugeGroup(direction int) *GaugeGroup {
	return &GaugeGroup{
		Flex:   tview.NewFlex().SetDirection(direction),
		gauges: []*GaugeWidget{},
	}
}

// AddGauge adds a gauge to the group
func (gg *GaugeGroup) AddGauge(title string, minVal, maxVal float64, unit string) *GaugeWidget {
	gauge := NewGaugeWidget(title, minVal, maxVal, unit)
	gg.gauges = append(gg.gauges, gauge)
	gg.AddItem(gauge, 3, 1, false) // Height of 3 for border + gauge + padding
	return gauge
}

// Update updates all gauges in the group
func (gg *GaugeGroup) Update(values map[string]float64) {
	for _, gauge := range gg.gauges {
		if value, ok := values[gauge.title]; ok {
			gauge.SetValue(value)
		}
	}
}
