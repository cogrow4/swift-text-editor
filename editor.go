package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Editor modes
type EditorMode int

const (
	ViewMode EditorMode = iota
	EditMode
)

// Enhanced text editor with Vim-like modes
type TextEditor struct {
	app           *tview.Application
	textView      *tview.TextView
	statusBar     *tview.TextView
	helpModal     *tview.Modal
	fileModal     *tview.Modal
	saveForm      *tview.Form
	openForm      *tview.Form
	filePath      string
	content       []string
	lineNum       int
	colNum        int
	showHelp      bool
	showWelcome   bool
	modified      bool
	mode          EditorMode
	commandBuffer string
	showingDialog bool
}

func NewTextEditor(filePath string) *TextEditor {
	app := tview.NewApplication()

	editor := &TextEditor{
		app:           app,
		filePath:      filePath,
		content:       []string{""},
		lineNum:       0,
		colNum:        0,
		showWelcome:   filePath == "",
		modified:      false,
		mode:          ViewMode,
		commandBuffer: "",
		showingDialog: false,
	}

	editor.setupUI()
	return editor
}

func (e *TextEditor) setupUI() {
	// Create main text view with better configuration
	e.textView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(false).
		SetScrollable(true).
		SetChangedFunc(func() {
			e.app.Draw()
		})

	// Create status bar with more information
	e.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	// Create help modal
	e.helpModal = tview.NewModal().
		SetText(e.getHelpText()).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.showHelp = false
			e.app.SetRoot(e.getMainLayout(), true)
		})

	// Create save form
	e.saveForm = tview.NewForm()
	e.saveForm.AddInputField("Save as:", "", 50, nil, nil)
	e.saveForm.AddButton("Save", func() {
		filename := e.saveForm.GetFormItem(0).(*tview.InputField).GetText()
		if filename != "" {
			e.filePath = filename
			e.saveFile()
		}
		e.showingDialog = false
		e.app.SetRoot(e.getMainLayout(), true)
	})
	e.saveForm.AddButton("Cancel", func() {
		e.showingDialog = false
		e.app.SetRoot(e.getMainLayout(), true)
	})
	e.saveForm.SetBorder(true)
	e.saveForm.SetTitle("Save File")

	// Create open form
	e.openForm = tview.NewForm()
	e.openForm.AddInputField("Open file:", "", 50, nil, nil)
	e.openForm.AddButton("Open", func() {
		filename := e.openForm.GetFormItem(0).(*tview.InputField).GetText()
		if filename != "" {
			e.filePath = filename
			e.loadFile()
		}
		e.showingDialog = false
		e.app.SetRoot(e.getMainLayout(), true)
	})
	e.openForm.AddButton("Cancel", func() {
		e.showingDialog = false
		e.app.SetRoot(e.getMainLayout(), true)
	})
	e.openForm.SetBorder(true)
	e.openForm.SetTitle("Open File")

	// Set up enhanced key bindings
	e.setupKeyBindings()

	// Load file if specified
	if e.filePath != "" {
		e.loadFile()
	} else {
		e.showWelcomeScreen()
	}

	// Set initial layout
	e.app.SetRoot(e.getMainLayout(), true)
}

func (e *TextEditor) getMainLayout() tview.Primitive {
	if e.showHelp {
		return e.helpModal
	}

	if e.showingDialog {
		// This will be set by the specific dialog functions
		return nil
	}

	// Create flex layout with line numbers
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(e.textView, 0, 1, true).
		AddItem(e.statusBar, 1, 0, false)

	return flex
}

func (e *TextEditor) showWelcomeScreen() {
	welcomeText := `
╔══════════════════════════════════════════════════════════════╗
║                    SWIFT Text Editor                         ║
║         Streamlined Workflow, Increased Focus Typography     ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  Welcome to SWIFT! A modern, Vim-like text editor.           ║
║                                                              ║
║  🚀 Get Started: Press 'g' for help and commands             ║
║  📝 Open File: swift filename.txt                            ║
║  ⌨️ Vim-like: View Mode and Edit Mode                       ║
║                                                              ║
║  Key Features:                                               ║
║  • View Mode: Navigate and execute commands                 ║
║  • Edit Mode: Insert and edit text (press 'i')             ║
║  • Syntax highlighting                                       ║
║  • Mac keyboard compatible                                   ║
║  • Better than Vim!                                          ║
║                                                              ║
║  Press 'g' to see all commands and get started!              ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	e.textView.SetText(welcomeText)
	e.updateStatusBar("Welcome to SWIFT! Press 'g' for help")
}

func (e *TextEditor) setupKeyBindings() {
	e.textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle mode-specific behavior
		if e.mode == EditMode {
			return e.handleEditMode(event)
		} else {
			return e.handleViewMode(event)
		}
	})
}

func (e *TextEditor) handleEditMode(event *tcell.EventKey) *tcell.EventKey {
	// In Edit Mode, only handle text input and ESC to exit
	switch event.Key() {
	case tcell.KeyEscape:
		// Exit Edit Mode
		e.mode = ViewMode
		e.updateStatusBar("View Mode")
		return nil
	case tcell.KeyUp:
		e.moveUp()
		return nil
	case tcell.KeyDown:
		e.moveDown()
		return nil
	case tcell.KeyLeft:
		e.moveLeft()
		return nil
	case tcell.KeyRight:
		e.moveRight()
		return nil
	case tcell.KeyHome:
		e.moveToLineStart()
		return nil
	case tcell.KeyEnd:
		e.moveToLineEnd()
		return nil
	case tcell.KeyEnter:
		e.insertNewline()
		return nil
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.backspace()
		return nil
	case tcell.KeyDelete:
		e.delete()
		return nil
	case tcell.KeyTab:
		e.insertTab()
		return nil
	}

	// Handle regular characters for text input
	if event.Rune() != 0 {
		e.insertChar(event.Rune())
		return nil
	}

	return event
}

func (e *TextEditor) handleViewMode(event *tcell.EventKey) *tcell.EventKey {
	// Handle special keys
	switch event.Key() {
	case tcell.KeyCtrlQ:
		if e.modified {
			// TODO: Add unsaved changes warning
		}
		e.app.Stop()
		return nil
	case tcell.KeyUp:
		e.moveUp()
		return nil
	case tcell.KeyDown:
		e.moveDown()
		return nil
	case tcell.KeyLeft:
		e.moveLeft()
		return nil
	case tcell.KeyRight:
		e.moveRight()
		return nil
	case tcell.KeyHome:
		e.moveToLineStart()
		return nil
	case tcell.KeyEnd:
		e.moveToLineEnd()
		return nil
	case tcell.KeyEnter:
		// Execute command buffer
		e.executeCommand()
		return nil
	}

	// Handle regular characters for commands
	if event.Rune() != 0 {
		switch event.Rune() {
		case 'i':
			// Enter Edit Mode
			e.mode = EditMode
			e.updateStatusBar("Edit Mode - Press ESC to exit")
			return nil
		case 'g':
			if e.showWelcome {
				e.showHelp = true
				e.app.SetRoot(e.helpModal, true)
				return nil
			}
		case 'h':
			e.showHelp = true
			e.app.SetRoot(e.helpModal, true)
			return nil
		default:
			// Add to command buffer
			e.commandBuffer += string(event.Rune())
			e.updateStatusBar(fmt.Sprintf("View Mode - Command: %s", e.commandBuffer))
			return nil
		}
	}

	return event
}

func (e *TextEditor) executeCommand() {
	command := e.commandBuffer
	e.commandBuffer = ""

	switch command {
	case "q":
		if e.showWelcome {
			e.app.Stop()
		} else {
			e.app.Stop()
		}
	case "w":
		e.saveFile()
	case "wq":
		e.saveFile()
		e.app.Stop()
	case "o":
		e.openFile()
	case "n":
		e.newFile()
	case "h":
		e.showHelp = true
		e.app.SetRoot(e.helpModal, true)
	default:
		e.updateStatusBar(fmt.Sprintf("Unknown command: %s", command))
	}
}

func (e *TextEditor) getHelpText() string {
	return `
╔══════════════════════════════════════════════════════════════╗
║                    SWIFT Help & Commands                    ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  🎯 MODES (Vim-like):                                       ║
║  • View Mode: Navigate and execute commands                 ║
║  • Edit Mode: Insert and edit text (press 'i' to enter)    ║
║  • ESC: Exit Edit Mode and return to View Mode             ║
║                                                              ║
║  🎯 NAVIGATION (Mac Compatible):                            ║
║  • Arrow Keys: Move cursor in all directions                ║
║  • Home/End: Jump to beginning/end of line                  ║
║  • Ctrl+Home/End: Jump to start/end of file                 ║
║                                                              ║
║  📝 EDITING (Edit Mode Only):                               ║
║  • Type normally to insert text                             ║
║  • Backspace/Delete: Remove characters                      ║
║  • Enter: New line                                          ║
║  • Tab: Smart indentation (4 spaces)                        ║
║  • ESC: Exit Edit Mode                                      ║
║                                                              ║
║  💾 COMMANDS (View Mode + Enter):                           ║
║  • 'w' + Enter: Save file (prompts for filename)            ║
║  • 'q' + Enter: Quit                                        ║
║  • 'wq' + Enter: Save and quit                              ║
║  • 'o' + Enter: Open file (prompts for path)                ║
║  • 'n' + Enter: New file                                    ║
║  • 'h' + Enter: Show this help                              ║
║                                                              ║
║  🎨 FEATURES:                                               ║
║  • Syntax highlighting for many languages                   ║
║  • Vim-like modes (View/Edit)                               ║
║  • Mac keyboard compatible                                  ║
║  • Mouse-free operation                                     ║
║  • Smart indentation                                        ║
║                                                              ║
║  💡 TIP: Press 'i' to edit, ESC to return to view mode!    ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
}

// Movement functions
func (e *TextEditor) moveUp() {
	if e.lineNum > 0 {
		e.lineNum--
		if e.colNum >= len(e.content[e.lineNum]) {
			e.colNum = len(e.content[e.lineNum])
		}
		e.updateDisplay()
	}
}

func (e *TextEditor) moveDown() {
	if e.lineNum < len(e.content)-1 {
		e.lineNum++
		if e.colNum >= len(e.content[e.lineNum]) {
			e.colNum = len(e.content[e.lineNum])
		}
		e.updateDisplay()
	}
}

func (e *TextEditor) moveLeft() {
	if e.colNum > 0 {
		e.colNum--
	} else if e.lineNum > 0 {
		e.lineNum--
		e.colNum = len(e.content[e.lineNum])
	}
	e.updateDisplay()
}

func (e *TextEditor) moveRight() {
	if e.colNum < len(e.content[e.lineNum]) {
		e.colNum++
	} else if e.lineNum < len(e.content)-1 {
		e.lineNum++
		e.colNum = 0
	}
	e.updateDisplay()
}

func (e *TextEditor) moveToLineStart() {
	e.colNum = 0
	e.updateDisplay()
}

func (e *TextEditor) moveToLineEnd() {
	e.colNum = len(e.content[e.lineNum])
	e.updateDisplay()
}

// Page Up/Down removed for Mac keyboard compatibility

// Editing functions
func (e *TextEditor) insertChar(char rune) {
	if e.showWelcome {
		return
	}

	line := e.content[e.lineNum]
	before := line[:e.colNum]
	after := line[e.colNum:]
	e.content[e.lineNum] = before + string(char) + after
	e.colNum++
	e.modified = true
	e.updateDisplay()
}

func (e *TextEditor) insertNewline() {
	if e.showWelcome {
		return
	}

	line := e.content[e.lineNum]
	before := line[:e.colNum]
	after := line[e.colNum:]

	e.content[e.lineNum] = before
	newLine := after

	// Insert new line
	e.content = append(e.content[:e.lineNum+1], append([]string{newLine}, e.content[e.lineNum+1:]...)...)
	e.lineNum++
	e.colNum = 0
	e.modified = true
	e.updateDisplay()
}

func (e *TextEditor) backspace() {
	if e.showWelcome {
		return
	}

	if e.colNum > 0 {
		line := e.content[e.lineNum]
		before := line[:e.colNum-1]
		after := line[e.colNum:]
		e.content[e.lineNum] = before + after
		e.colNum--
		e.modified = true
	} else if e.lineNum > 0 {
		// Join with previous line
		prevLine := e.content[e.lineNum-1]
		currentLine := e.content[e.lineNum]
		e.content[e.lineNum-1] = prevLine + currentLine
		e.content = append(e.content[:e.lineNum], e.content[e.lineNum+1:]...)
		e.lineNum--
		e.colNum = len(prevLine)
		e.modified = true
	}
	e.updateDisplay()
}

func (e *TextEditor) delete() {
	if e.showWelcome {
		return
	}

	line := e.content[e.lineNum]
	if e.colNum < len(line) {
		before := line[:e.colNum]
		after := line[e.colNum+1:]
		e.content[e.lineNum] = before + after
		e.modified = true
	} else if e.lineNum < len(e.content)-1 {
		// Join with next line
		nextLine := e.content[e.lineNum+1]
		e.content[e.lineNum] = line + nextLine
		e.content = append(e.content[:e.lineNum+1], e.content[e.lineNum+2:]...)
		e.modified = true
	}
	e.updateDisplay()
}

func (e *TextEditor) insertTab() {
	if e.showWelcome {
		return
	}

	// Smart indentation - insert 4 spaces
	for i := 0; i < 4; i++ {
		e.insertChar(' ')
	}
}

func (e *TextEditor) updateDisplay() {
	if e.showWelcome {
		return
	}

	// Create display with line numbers, syntax highlighting, and cursor indicator
	var display strings.Builder

	for i, line := range e.content {
		// Add line number with highlighting for current line
		if i == e.lineNum {
			display.WriteString(fmt.Sprintf("[yellow:blue]%3d[white] | ", i+1))
		} else {
			display.WriteString(fmt.Sprintf("%3d | ", i+1))
		}

		// Add syntax highlighted line with cursor indicator
		if i == e.lineNum {
			// Current line - show cursor position and highlight line background only
			display.WriteString("[white:blue]" + e.highlightLineWithCursor(line) + "[white]")
		} else {
			// Other lines - normal highlighting
			display.WriteString(e.highlightLine(line))
		}
		display.WriteString("\n")
	}

	e.textView.SetText(display.String())

	// Update status bar with mode information
	modeText := "View Mode"
	if e.mode == EditMode {
		modeText = "Edit Mode"
	}

	status := fmt.Sprintf("SWIFT | %s | %s | Line %d, Col %d",
		e.getStatusText(), modeText, e.lineNum+1, e.colNum+1)
	if e.modified {
		status += " | MODIFIED"
	}
	e.statusBar.SetText(status)
}

func (e *TextEditor) getStatusText() string {
	if e.filePath == "" {
		return "Untitled"
	}
	return filepath.Base(e.filePath)
}

func (e *TextEditor) highlightLine(line string) string {
	// Simple syntax highlighting based on file extension
	if e.filePath == "" {
		return line
	}

	ext := strings.ToLower(filepath.Ext(e.filePath))

	switch ext {
	case ".go":
		return e.highlightGo(line)
	case ".py":
		return e.highlightPython(line)
	case ".js", ".ts":
		return e.highlightJavaScript(line)
	case ".html", ".htm":
		return e.highlightHTML(line)
	case ".css":
		return e.highlightCSS(line)
	case ".json":
		return e.highlightJSON(line)
	default:
		return line
	}
}

func (e *TextEditor) highlightLineWithCursor(line string) string {
	// Add cursor indicator at the current column position
	if e.colNum >= len(line) {
		// Cursor at end of line - add cursor indicator after highlighting
		highlighted := e.highlightLine(line)
		return highlighted + "[black:white]▌[white]"
	} else {
		// Cursor in middle of line - split line and highlight each part
		before := line[:e.colNum]
		after := line[e.colNum:]

		// Highlight the parts separately to avoid color code issues
		highlightedBefore := e.highlightLine(before)
		highlightedAfter := e.highlightLine(after)

		return highlightedBefore + "[black:white]▌[white]" + highlightedAfter
	}
}

func (e *TextEditor) highlightGo(line string) string {
	// Simple Go syntax highlighting
	keywords := []string{"package", "import", "func", "var", "const", "type", "struct", "interface", "if", "else", "for", "range", "return", "go", "defer", "select", "case", "default", "switch", "break", "continue", "fallthrough"}

	result := line
	for _, keyword := range keywords {
		result = strings.ReplaceAll(result, keyword, fmt.Sprintf("[blue]%s[white]", keyword))
	}

	// Highlight strings
	result = strings.ReplaceAll(result, "\"", "[green]\"[white]")

	return result
}

func (e *TextEditor) highlightPython(line string) string {
	keywords := []string{"def", "class", "if", "else", "elif", "for", "while", "import", "from", "return", "yield", "try", "except", "finally", "with", "as", "pass", "break", "continue", "and", "or", "not", "in", "is", "lambda", "True", "False", "None"}

	result := line
	for _, keyword := range keywords {
		result = strings.ReplaceAll(result, keyword, fmt.Sprintf("[blue]%s[white]", keyword))
	}

	return result
}

func (e *TextEditor) highlightJavaScript(line string) string {
	keywords := []string{"function", "var", "let", "const", "if", "else", "for", "while", "return", "class", "extends", "import", "export", "async", "await", "try", "catch", "finally", "throw", "new", "this", "true", "false", "null", "undefined"}

	result := line
	for _, keyword := range keywords {
		result = strings.ReplaceAll(result, keyword, fmt.Sprintf("[blue]%s[white]", keyword))
	}

	return result
}

func (e *TextEditor) highlightHTML(line string) string {
	// Simple HTML highlighting
	result := line
	result = strings.ReplaceAll(result, "<", "[red]<[white]")
	result = strings.ReplaceAll(result, ">", "[red]>[white]")
	return result
}

func (e *TextEditor) highlightCSS(line string) string {
	// Simple CSS highlighting
	keywords := []string{"color", "background", "margin", "padding", "border", "width", "height", "display", "position", "float", "clear", "font", "text", "line", "letter", "word", "white", "space", "overflow", "visibility", "opacity", "z-index"}

	result := line
	for _, keyword := range keywords {
		result = strings.ReplaceAll(result, keyword, fmt.Sprintf("[blue]%s[white]", keyword))
	}

	return result
}

func (e *TextEditor) highlightJSON(line string) string {
	// Simple JSON highlighting
	result := line
	result = strings.ReplaceAll(result, "\"", "[green]\"[white]")
	result = strings.ReplaceAll(result, ":", "[yellow]:[white]")
	result = strings.ReplaceAll(result, ",", "[yellow],[white]")
	return result
}

func (e *TextEditor) loadFile() {
	if e.filePath == "" {
		return
	}

	content, err := os.ReadFile(e.filePath)
	if err != nil {
		e.textView.SetText(fmt.Sprintf("Error loading file: %v\n\nPress 'n' for new file or 'o' to open another file.", err))
		e.updateStatusBar(fmt.Sprintf("Error: %v", err))
		return
	}

	// Split content into lines
	e.content = strings.Split(string(content), "\n")
	if len(e.content) == 0 {
		e.content = []string{""}
	}

	e.lineNum = 0
	e.colNum = 0
	e.showWelcome = false
	e.modified = false
	e.updateDisplay()
}

func (e *TextEditor) saveFile() {
	if e.filePath == "" {
		e.showSaveDialog()
		return
	}

	content := strings.Join(e.content, "\n")
	err := os.WriteFile(e.filePath, []byte(content), 0644)
	if err != nil {
		e.updateStatusBar(fmt.Sprintf("Error saving: %v", err))
	} else {
		e.modified = false
		e.updateStatusBar(fmt.Sprintf("Saved: %s", filepath.Base(e.filePath)))
	}
}

func (e *TextEditor) showSaveDialog() {
	e.showingDialog = true
	e.app.SetRoot(e.saveForm, true)
}

func (e *TextEditor) openFile() {
	e.showOpenDialog()
}

func (e *TextEditor) showOpenDialog() {
	e.showingDialog = true
	e.app.SetRoot(e.openForm, true)
}

func (e *TextEditor) newFile() {
	e.filePath = ""
	e.content = []string{""}
	e.lineNum = 0
	e.colNum = 0
	e.showWelcome = false
	e.modified = false
	e.updateDisplay()
}

func (e *TextEditor) updateStatusBar(message string) {
	e.statusBar.SetText(fmt.Sprintf("SWIFT | %s", message))
}

func (e *TextEditor) Run() error {
	return e.app.Run()
}
