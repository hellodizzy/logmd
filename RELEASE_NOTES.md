# logmd v1.0.0 - Complete Journal CLI 📖

A minimal, local-first journal CLI written in Go. Complete implementation with all core features.

## 🎯 Features
- ✅ **Daily entries**: `logmd today` - create/edit journal entries
- ✅ **Interactive timeline**: `logmd timeline` - browse entries with beautiful T
- ✅ **Entry viewing**: `logmd view YYYY-MM-DD` - render markdown beautifully
- ✅ **Configuration**: `logmd config` - manage settings
- ✅ **Cross-platform**: macOS (Intel/ARM) and Linux (Intel/ARM)

## 📦 Installation

### Homebrew (recommended)
```bash
brew install hellodizzy/tap/logmd
```

### Direct Download
Download the appropriate binary for your platform below.

## 🧪 Usage
```bash
logmd today           # Create/edit today's entry
logmd timeline        # Browse all entries
logmd view 2024-01-15 # View specific entry
logmd config          # Check configuration
```

## 🔧 Technical Details
- **Language**: Go 1.22+
- **Dependencies**: MIT/Apache 2.0 licensed only
- **Architecture**: Cross-platform static binaries
- **Size**: ~18MB per binary
- **Tests**: 57 passing tests across all packages

## 📋 What's Included
This release includes all 5 development phases:
- **Phase 0**: Bootstrap (CLI setup, CI pipeline)  
- **Phase 1**: Vault package (file system operations)
- **Phase 2**: `today` command (entry creation/editing)
- **Phase 3**: `timeline` command (interactive TUI browsing)
- **Phase 4**: `view` command (markdown rendering)
- **Phase 5**: `config` command (configuration display)

## 🏗️ Architecture
Built with professional Go libraries:
- **CLI**: Cobra for command-line interfaces
- **TUI**: Bubble Tea for interactive terminals
- **Config**: Viper for flexible configuration
- **Markdown**: Goldmark + Glamour for parsing/rendering
- **Styling**: Lipgloss for beautiful terminal output

## 📄 License
MIT License - fully open source and commercial-friendly. p