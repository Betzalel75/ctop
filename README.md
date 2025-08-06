<p align="center"><img width="200px" src="/_docs/img/logo.png" alt="ctop"/></p>

#

![License: MIT](https://img.shields.io/badge/license-MIT-blue)
![release](https://img.shields.io/github/v/release/Betzalel75/ctop)

Enhanced fork of ctop with additional features

`ctop` provides a concise and condensed overview of real-time metrics for multiple containers:
<p align="center"><img src="_docs/img/grid.gif" alt="ctop"/></p>

This fork maintains all original functionality while adding:
- Improved container management interface
- Additional viewing options
- Enhanced navigation controls
- Extended filtering capabilities

## Installation

### From Source (Recommended)

#### Prerequisites
- Go 1.24+
- Git

```bash
git clone https://github.com/Betzalel75/ctop.git
cd ctop
make build
sudo mv ctop /usr/local/bin/
```

### Pre-built Binaries (Linux/macOS)


**curl :**

```bash
curl -fsSL https://raw.githubusercontent.com/Betzalel75/ctop/master/scripts/install.sh | sh
````

**wget :**

```bash
wget -qO- https://raw.githubusercontent.com/Betzalel75/ctop/master/scripts/install.sh | sh
```

---

## 🧹 Uninstallation

**curl :**

```bash
curl -fsSL https://raw.githubusercontent.com/Betzalel75/ctop/master/scripts/uninstall.sh | sh
```

**wget :**

```bash
wget -qO- https://raw.githubusercontent.com/Betzalel75/ctop/master/scripts/uninstall.sh | sh
```


Download the latest binary for Linux amd64:

```bash
sudo wget https://github.com/Betzalel75/ctop/releases/download/vX.X.X/ctop-X.X.X-linux-amd64 -O /usr/local/bin/ctop
sudo chmod +x /usr/local/bin/ctop
```

### Docker

```bash
docker run --rm -ti \
--name=ctop \
--volume /var/run/docker.sock:/var/run/docker.sock:ro \
ghcr.io/Betzalel75/ctop:latest
```

## Building

Build steps:

1. Clone the repository:
```bash
git clone https://github.com/Betzalel75/ctop.git
cd ctop
```

2. Build for your current platform:
```bash
make build
```

3. For cross-platform builds:
```bash
make build-all
```

Build artifacts will be placed in the `_build` directory.

## New Features

This fork includes several enhancements over the original ctop:

- **Extended Container Management**:
  - Bulk operations for containers/images/volumes
  - Advanced filtering options

- **Improved UI**:
  - Dual-pane interface (Running/All views)
  - Better pagination support
  - Enhanced help system

- **Additional Functionality**:
  - Image publishing tools
  - Volume management
  - Extended container actions

## Usage

Basic usage remains the same as original ctop:

```bash
ctop
```

### New Keybindings

| Key | Action |
|-----|--------|
| <kbd>tab</kbd> | Toggle between Running/All views |
| <kbd>1</kbd> | Switch to Running view |
| <kbd>2</kbd> | Switch to All view |
| <kbd>d</kbd> (All view) | Delete selected items |
| <kbd>space</kbd> (All view) | Toggle item selection |

## Contributing

Contributions are welcome! Please open an issue or pull request on GitHub.

## License

MIT - See [LICENSE](LICENSE) file

## Acknowledgments

This project is a fork of the original [ctop](https://github.com/bcicen/ctop) by bcicen.

