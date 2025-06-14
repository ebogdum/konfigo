# Installation

There are a couple of ways to install Konfigo:

## From Source

Alternatively, you can clone the repository and build from source:

```bash
git clone https://github.com/ebogdum/konfigo.git
cd konfigo
cd cmd/konfigo
go build
```

This will create a `konfigo` binary in the `cmd/konfigo` directory. You can then move this binary to a directory in your `PATH`, like `/usr/local/bin` or `~/bin`.

## Pre-compiled Binaries (If Available)

If pre-compiled binaries are provided for your operating system on the project's GitHub Releases page, you can download the appropriate binary, make it executable, and move it to a directory in your `PATH`.

```bash
# Example for Linux/macOS
wget https://github.com/ebogdum/konfigo/releases/download/v1.0.0/konfigo_linux_amd64 # Adjust URL
chmod +x konfigo_linux_amd64
sudo mv konfigo_linux_amd64 /usr/local/bin/konfigo
```

Verify the installation by running:
```bash
konfigo -h
```
This should display the help message for Konfigo.
