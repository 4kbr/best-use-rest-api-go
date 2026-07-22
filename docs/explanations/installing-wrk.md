# Installing wrk

## Linux

### Debian/Ubuntu
```bash
sudo apt update && sudo apt install wrk
```

### Fedora
```bash
sudo dnf install wrk
```

### Arch Linux
```bash
sudo pacman -S wrk
```

### From source (jika tidak tersedia di package manager)
```bash
sudo apt install build-essential libssl-dev git -y
git clone https://github.com/wg/wrk.git
cd wrk && make && sudo cp wrk /usr/local/bin/
```

## macOS

### Homebrew
```bash
brew install wrk
```

## Windows

wrk tidak support native Windows. Opsi:

### Opsi 1 — WSL (recommended)
Install WSL + Ubuntu, lalu:
```bash
sudo apt update && sudo apt install wrk
```

### Opsi 2 — Scoop (port wrk untuk Windows)
```bash
scoop install wrk
```
