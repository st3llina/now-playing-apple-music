# 🎵 Apple Music Now Playing

A lightweight **Now Playing tracker** built with Go that retrieves real-time music information from Apple Music on Windows using Windows Runtime (WinRT).

## ✨ Features

- 🎵 Track metadata — Title, Artist, and Album
- 🖼️ Album Artwork extraction and local caching
- ⏱️ Playback progress and track duration
- ▶️ Playback status detection (Playing / Paused)
- ⚡ Windows Media Session integration
- 🚀 Lightweight and written in Go

## 🛠️ Tech Stack

- **Language:** Go
- **Platform:** Windows 10 / 11
- **API:** Windows Runtime (WinRT)
- **Media API:** Windows.Media.Control
- **Bindings:** [go-bindings-winrt](https://github.com/deploymenttheory/go-bindings-winrt)

## 📦 Installation

Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/apple-now-playing.git
cd apple-now-playing
```

Install dependencies:

```bash
go mod tidy
```

## 🚀 Usage

Make sure Apple Music is running and playing a track.

Run the application:

```bash
go run .
```

The application retrieves the current media session and reads its metadata.

## 🖼️ Album Artwork

Album artwork is extracted from the Windows Media Session thumbnail stream and saved locally.

Example:

```text
artworks/
├── Artist - Track One.jpg
├── Artist - Track Two.jpg
└── Artist - Track Three.jpg
```

Existing artwork files can be reused to avoid unnecessary writes.

## 📋 Requirements

- Windows 10 / Windows 11
- Go 1.25+
- Apple Music for Windows
- A media session exposed through Windows Media Control

## ⚙️ How It Works

```text
Apple Music
     │
     ▼
Windows Media Session API
     │
     ▼
Go + WinRT Bindings
     │
     ├── Track Metadata
     ├── Playback Information
     ├── Timeline Properties
     └── Album Artwork
             │
             ▼
       Local Artwork Cache
```

## ⚠️ Limitations

- Windows-specific implementation.
- Requires an active media session.
- Available metadata depends on what the media application exposes through Windows Media Control.
- Other media applications may also be detected.

## 📌 Project Status

**Work in Progress**

This project is actively being developed. Features and APIs may change.

## 📄 License

MIT License.
