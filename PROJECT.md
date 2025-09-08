# SWIFT Text Editor - Project Overview

## 🚀 What is SWIFT?

**SWIFT** (Streamlined Workflow, Increased Focus Typography) is a modern, intuitive CLI text editor designed to be better than Vim with smooth, mouse-free operation and easy-to-remember commands.

## ✨ Key Features

### 🎯 Vim-like Modes
- **View Mode**: Navigate and execute commands
- **Edit Mode**: Insert and edit text (press 'i' to enter, ESC to exit)

### 🎨 Visual Excellence
- **Visible Cursor**: Clear cursor indicator with `▌` symbol
- **Syntax Highlighting**: Support for Go, Python, JavaScript, HTML, CSS, JSON, RTF
- **Line Highlighting**: Current line highlighted in blue
- **Line Numbers**: Clear numbering with current line emphasis

### ⌨️ Intuitive Commands
- **Navigation**: Arrow keys, Home/End work as expected
- **File Operations**: `w` (save), `o` (open), `q` (quit), `wq` (save & quit)
- **Mac Compatible**: No Page Up/Down, optimized for Mac keyboards

### 💾 File Management
- **Save Dialog**: Prompts for filename when saving new files
- **Open Dialog**: Prompts for file path when opening files
- **Smart Dialogs**: Clean forms with input fields and buttons

## 🛠️ Technical Stack

- **Language**: Go 1.21+
- **UI Framework**: tview (terminal UI library)
- **Dependencies**: tcell for terminal control
- **Build System**: Simple shell script

## 📁 Project Structure

```
swift-text-editor/
├── main.go          # Entry point and CLI parsing
├── editor.go        # Core editor functionality
├── go.mod           # Go module dependencies
├── build.sh         # Build script
├── README.md        # User documentation
├── PROJECT.md       # This project overview
├── .gitignore       # Git ignore rules
├── sample.txt       # Sample text file
└── test.rtf         # RTF test file
```

## 🎯 Why SWIFT?

Unlike Vim, SWIFT is designed with modern usability in mind:

- ✅ **Intuitive**: Commands make sense (no cryptic key combinations)
- ✅ **Smooth**: Fluid cursor movement and text manipulation
- ✅ **Fast**: Built with Go for excellent performance
- ✅ **Modern**: Clean interface with syntax highlighting
- ✅ **Accessible**: Easy to learn and remember

## 🚀 Getting Started

### Quick Build
```bash
./build.sh
```

### Usage
```bash
swft filename.txt    # Edit a file
swft                 # Start with welcome screen
```

### Installation
```bash
sudo cp swft /usr/local/bin/
```

## 🎨 Screenshots

The editor features:
- Welcome screen with help
- Syntax highlighted text
- Visible cursor indicator
- File dialogs for save/open
- Status bar with mode and position info

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📝 License

This project is open source. Feel free to contribute and improve SWIFT!

---

**SWIFT** - Because text editing should be intuitive, not cryptic! 🚀
