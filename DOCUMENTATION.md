# Documentation Setup

This project uses [mdBook](https://rust-lang.github.io/mdBook/) to generate beautiful documentation sites from Markdown files.

## How it Works

### 1. **mdBook Structure**
```
book/
├── book.toml          # Configuration file
└── src/
    ├── SUMMARY.md     # Table of contents
    ├── README.md      # Main documentation
    ├── RELEASE.md     # Release guide
    └── commands/      # Command documentation
        ├── git-status.md
        ├── git-sync.md
        ├── k8s-kctx.md
        └── ...
```

### 2. **GitHub Pages Workflow**
When you push to the `main` branch, GitHub Actions automatically:

1. **Sets up mdBook** using `peaceiris/actions-mdbook@v1`
2. **Creates book structure** with all documentation files
3. **Builds the site** using `mdbook build book`
4. **Deploys to GitHub Pages** at `https://nghiadaulau.github.io/opsbrew`

### 3. **Features**
- ✅ **Fast**: mdBook is written in Rust, very fast
- ✅ **Simple**: Just Markdown files, no complex setup
- ✅ **Beautiful**: Clean, responsive design
- ✅ **Search**: Built-in search functionality
- ✅ **Navigation**: Automatic table of contents
- ✅ **Dark mode**: Toggle between light/dark themes
- ✅ **Mobile friendly**: Responsive design

## Local Development

### Install mdBook
```bash
# Using cargo (if you have Rust installed)
cargo install mdbook

# Or download from releases
# https://github.com/rust-lang/mdBook/releases
```

### Build Locally
```bash
# Create book structure (first time only)
mkdir -p book/src
cp README.md book/src/
cp RELEASE.md book/src/

# Create book.toml and SUMMARY.md (see workflow for content)

# Build the book
mdbook build book

# Serve locally for preview
mdbook serve book --open
```

### Add New Documentation

1. **Add new Markdown file** to `book/src/`
2. **Update SUMMARY.md** to include the new page
3. **Push to main** - GitHub Actions will rebuild automatically

Example:
```bash
# Add new command documentation
echo "# New Command" > book/src/commands/new-command.md

# Update SUMMARY.md
echo "- [new command](commands/new-command.md)" >> book/src/SUMMARY.md

# Push changes
git add .
git commit -m "Add new command documentation"
git push origin main
```

## Configuration

The mdBook configuration is in `book/book.toml`:

```toml
[book]
authors = ["nghiadaulau"]
language = "en"
multilingual = false
src = "src"
title = "opsbrew Documentation"

[output.html]
git-repository-url = "https://github.com/nghiadaulau/opsbrew"
git-repository-icon = "fa-github"
default-theme = "light"
preferred-dark-theme = "navy"
mathjax-support = true

[output.html.fold]
enable = true
level = 1
```

## Customization

### Themes
mdBook supports multiple themes:
- `light` (default)
- `navy` (dark)
- `coal` (dark)
- `ayu` (dark)

### Additional Features
- **MathJax**: For mathematical equations
- **Git repository**: Links back to GitHub
- **Folding**: Collapsible sections
- **Search**: Built-in search functionality

## Troubleshooting

### Common Issues

1. **Build fails**: Check that all files referenced in `SUMMARY.md` exist
2. **Missing pages**: Ensure files are in the correct directory structure
3. **Broken links**: Verify all internal links in Markdown files

### Debug Locally
```bash
# Check for errors
mdbook build book --log-level debug

# Validate links
mdbook test book
```

## Benefits over Jekyll

| Feature | Jekyll | mdBook |
|---------|--------|--------|
| Setup complexity | High | Low |
| Build speed | Slow | Fast |
| Dependencies | Ruby, gems | Single binary |
| Configuration | Complex | Simple |
| Markdown support | Good | Excellent |
| Search | Plugin needed | Built-in |
| Mobile | Depends on theme | Always responsive |

mdBook is perfect for technical documentation - simple, fast, and beautiful! 🚀
