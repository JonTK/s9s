package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jontk/s9s/internal/dao"
	"github.com/jontk/s9s/pkg/slurm"
	"github.com/rivo/tview"
)

// JobOutputView displays job output (stdout/stderr)
type JobOutputView struct {
	app      *tview.Application
	pages    *tview.Pages
	client   dao.SlurmClient
	modal    *tview.Flex
	textView *tview.TextView
	jobID    string
	jobName  string
	outputType string // "stdout" or "stderr"
	autoRefresh bool
	refreshTicker *time.Ticker
}

// NewJobOutputView creates a new job output view
func NewJobOutputView(client dao.SlurmClient, app *tview.Application) *JobOutputView {
	return &JobOutputView{
		client: client,
		app:    app,
	}
}

// SetPages sets the pages manager for modal display
func (v *JobOutputView) SetPages(pages *tview.Pages) {
	v.pages = pages
}

// ShowJobOutput displays job output in a modal
func (v *JobOutputView) ShowJobOutput(jobID, jobName, outputType string) {
	v.jobID = jobID
	v.jobName = jobName
	v.outputType = outputType
	v.autoRefresh = false

	v.buildUI()
	v.loadOutput()
	v.show()
}

// buildUI creates the job output viewer interface
func (v *JobOutputView) buildUI() {
	// Create text view for output
	v.textView = tview.NewTextView()
	v.textView.SetDynamicColors(true)
	v.textView.SetScrollable(true)
	v.textView.SetWrap(true)
	v.textView.SetBorder(true)
	v.textView.SetTitle(fmt.Sprintf(" Job %s - %s (%s) ", v.jobID, v.jobName, strings.ToUpper(v.outputType)))
	v.textView.SetTitleAlign(tview.AlignCenter)

	// Create help text
	helpText := tview.NewTextView()
	helpText.SetDynamicColors(true)
	helpText.SetText("[yellow]Keys:[white] r=Refresh a=Auto-refresh s=Switch stdout/stderr f=Follow e=Export Esc=Close")
	helpText.SetTextAlign(tview.AlignCenter)

	// Create modal container
	v.modal = tview.NewFlex()
	v.modal.SetDirection(tview.FlexRow)
	v.modal.AddItem(nil, 0, 1, false)
	v.modal.AddItem(tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(v.textView, 0, 1, true).
			AddItem(helpText, 1, 0, false), 0, 4, true).
		AddItem(nil, 0, 1, false), 0, 3, true)
	v.modal.AddItem(nil, 0, 1, false)

	v.modal.SetBorder(true)
	v.modal.SetTitle(" Job Output Viewer ")
	v.modal.SetTitleAlign(tview.AlignCenter)

	// Setup event handlers
	v.setupEventHandlers()
}

// setupEventHandlers configures keyboard shortcuts
func (v *JobOutputView) setupEventHandlers() {
	v.modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			v.close()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'r', 'R':
				v.loadOutput()
				return nil
			case 'a', 'A':
				v.toggleAutoRefresh()
				return nil
			case 's', 'S':
				v.switchOutputType()
				return nil
			case 'f', 'F':
				v.followOutput()
				return nil
			case 'e', 'E':
				v.exportOutput()
				return nil
			}
		}
		return event
	})
}

// loadOutput loads job output from SLURM
func (v *JobOutputView) loadOutput() {
	v.textView.SetText("[yellow]Loading output...[white]")

	go func() {
		var content string
		var err error

		// Try to get job output from client
		if jobMgr := v.client.Jobs(); jobMgr != nil {
			// For mock client, generate sample output
			if _, ok := v.client.(*slurm.MockClient); ok {
				content = v.generateMockOutput()
			} else {
				// For real client, try to get actual output
				content, err = v.getJobOutput()
				if err != nil {
					content = fmt.Sprintf("[red]Error loading output: %v[white]\n\nThis could mean:\n• Job is still running\n• Output files not accessible\n• Job completed without output\n• Insufficient permissions", err)
				}
			}
		} else {
			content = "[red]Job manager not available[white]"
		}

		// Update UI on main thread
		v.app.QueueUpdateDraw(func() {
			v.textView.SetText(content)
			v.textView.ScrollToEnd()
		})
	}()
}

// getJobOutput retrieves actual job output (placeholder for real implementation)
func (v *JobOutputView) getJobOutput() (string, error) {
	// This would be implemented to read actual SLURM output files
	// For now, return a placeholder
	return fmt.Sprintf("Output for job %s (%s) would be retrieved from SLURM output files.\n\nThis requires access to the job's stdout/stderr files on the cluster filesystem.", v.jobID, v.outputType), nil
}

// generateMockOutput creates sample output for demonstration
func (v *JobOutputView) generateMockOutput() string {
	if v.outputType == "stderr" {
		return fmt.Sprintf(`[red]STDERR for Job %s (%s)[white]

[yellow]Warning:[white] This is mock stderr output for demonstration
[red]Error:[white] Sample error message
[yellow]Debug:[white] Verbose debugging information

Job started at: %s
Working directory: /home/user/jobs
Environment variables loaded: 15

Some sample error output...
`, v.jobID, v.jobName, time.Now().Add(-10*time.Minute).Format("2006-01-02 15:04:05"))
	}

	return fmt.Sprintf(`[green]STDOUT for Job %s (%s)[white]

Job started at: %s
Loading modules...
  - module load gcc/11.2.0
  - module load openmpi/4.1.1

Setting up environment...
Working directory: /home/user/simulation
Input files: config.dat, mesh.inp

[yellow]Starting simulation...[white]
Iteration 1/1000: Energy = -1234.56
Iteration 100/1000: Energy = -1245.78
Iteration 200/1000: Energy = -1251.23
Iteration 300/1000: Energy = -1255.67
...
Iteration 1000/1000: Energy = -1289.45

[green]Simulation completed successfully![white]
Results written to: output/results.dat
Execution time: 15 minutes 23 seconds
Peak memory usage: 2.3 GB

Job completed at: %s
`, v.jobID, v.jobName, 
	time.Now().Add(-25*time.Minute).Format("2006-01-02 15:04:05"),
	time.Now().Add(-10*time.Minute).Format("2006-01-02 15:04:05"))
}

// toggleAutoRefresh toggles automatic refresh
func (v *JobOutputView) toggleAutoRefresh() {
	v.autoRefresh = !v.autoRefresh

	if v.autoRefresh {
		v.textView.SetTitle(fmt.Sprintf(" Job %s - %s (%s) [AUTO-REFRESH] ", v.jobID, v.jobName, strings.ToUpper(v.outputType)))
		v.startAutoRefresh()
	} else {
		v.textView.SetTitle(fmt.Sprintf(" Job %s - %s (%s) ", v.jobID, v.jobName, strings.ToUpper(v.outputType)))
		v.stopAutoRefresh()
	}
}

// startAutoRefresh starts automatic refresh timer
func (v *JobOutputView) startAutoRefresh() {
	if v.refreshTicker != nil {
		v.refreshTicker.Stop()
	}

	v.refreshTicker = time.NewTicker(3 * time.Second)
	go func() {
		for range v.refreshTicker.C {
			if v.autoRefresh {
				v.loadOutput()
			} else {
				break
			}
		}
	}()
}

// stopAutoRefresh stops automatic refresh
func (v *JobOutputView) stopAutoRefresh() {
	if v.refreshTicker != nil {
		v.refreshTicker.Stop()
		v.refreshTicker = nil
	}
}

// switchOutputType switches between stdout and stderr
func (v *JobOutputView) switchOutputType() {
	if v.outputType == "stdout" {
		v.outputType = "stderr"
	} else {
		v.outputType = "stdout"
	}

	v.textView.SetTitle(fmt.Sprintf(" Job %s - %s (%s) ", v.jobID, v.jobName, strings.ToUpper(v.outputType)))
	v.loadOutput()
}

// followOutput scrolls to end and enables auto-refresh
func (v *JobOutputView) followOutput() {
	v.textView.ScrollToEnd()
	if !v.autoRefresh {
		v.toggleAutoRefresh()
	}
}

// exportOutput exports the current output to a file
func (v *JobOutputView) exportOutput() {
	// For now, just show a notification
	// In a real implementation, this would save to a file
	v.showNotification("Export functionality would save output to file")
}

// showNotification shows a temporary notification
func (v *JobOutputView) showNotification(message string) {
	notification := tview.NewModal()
	notification.SetText(message)
	notification.AddButtons([]string{"OK"})
	notification.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		v.pages.RemovePage("notification")
	})

	v.pages.AddPage("notification", notification, true, true)
}

// show displays the modal
func (v *JobOutputView) show() {
	if v.pages != nil {
		v.pages.AddPage("job-output", v.modal, true, true)
		v.app.SetFocus(v.textView)
	}
}

// close closes the output viewer
func (v *JobOutputView) close() {
	v.stopAutoRefresh()
	if v.pages != nil {
		v.pages.RemovePage("job-output")
	}
}