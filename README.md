# SWIFT Text Editor

**Streamlined Workflow, Increased Focus Typography**

A modern, intuitive CLI text editor designed to be better than Vim with smooth, mouse-free operation and easy-to-remember commands.

## 🚀 Features

- **Intuitive Interface**: No cryptic key combinations like Vim
- **Syntax Highlighting**: Support for Go, Python, JavaScript, HTML, CSS, JSON, and more
- **Smooth Navigation**: Arrow keys, Home/End, Page Up/Down work as expected
- **Mouse-Free Operation**: Complete keyboard control
- **Smart Indentation**: Automatic indentation with Tab key
- **Line Numbers**: Clear line numbering for easy navigation
- **Status Bar**: Shows current position and file status
- **Welcome Screen**: Get started quickly with built-in help

## 📦 Installation

### Quick Build
```bash
./build.sh
```

### Manual Build
```bash
go mod tidy
go build -o swft .
```

### Global Installation
```bash
# After building
sudo cp swft /usr/local/bin/
```

## 🎯 Usage

```bash
# Edit a file
swft filename.txt

# Start with welcome screen
swft

# Get help
swft --help
```

## ⌨️ Commands

### Navigation (Intuitive & Smooth)
- **Arrow Keys**: Move cursor in all directions
- **Home/End**: Jump to beginning/end of line
- **Page Up/Down**: Scroll by page
- **Ctrl+Home/End**: Jump to start/end of file

### Editing (Better than Vim)
- **Type normally**: Insert text at cursor
- **Backspace/Delete**: Remove characters
- **Enter**: New line
- **Tab**: Smart indentation (4 spaces)

### File Operations
- **Ctrl+S** or **'w'**: Save file
- **Ctrl+O** or **'o'**: Open file
- **Ctrl+N** or **'n'**: New file

### Help & Exit
- **'h'**: Show help
- **'g'**: Get started (from welcome screen)
- **Ctrl+Q** or **'q'**: Quit (from welcome screen)

## 🎨 Syntax Highlighting

SWIFT automatically detects file types and provides syntax highlighting for:

- **Go** (.go)
- **Python** (.py)
- **JavaScript/TypeScript** (.js, .ts)
- **HTML** (.html, .htm)
- **CSS** (.css)
- **JSON** (.json)

## 💡 Why SWIFT?

Unlike Vim, SWIFT is designed with modern usability in mind:

- ✅ **Intuitive**: Commands make sense (no cryptic key combinations)
- ✅ **Smooth**: Fluid cursor movement and text manipulation
- ✅ **Fast**: Built with Go for excellent performance
- ✅ **Modern**: Clean interface with syntax highlighting
- ✅ **Accessible**: Easy to learn and remember

## 🔧 Development

### Requirements
- Go 1.21 or later
- Terminal with color support

### Dependencies
- `github.com/gdamore/tcell/v2` - Terminal UI framework
- `github.com/rivo/tview` - Rich text widgets

### Building
```bash
go mod tidy
go build -o swift .
```

## 📝 License

This project is open source. Feel free to contribute and improve SWIFT!

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

**SWIFT** - Because text editing should be intuitive, not cryptic! 🚀
