# togram

CLI tool to pipe stdin or send files to Telegram — auto-detects text, photos, videos and documents.

```bash
echo "deploy done" | togram
cat screenshot.png | togram
togram report.pdf
togram -c @mychat video.mp4
```

## Install

### apt (Debian / Ubuntu)

```bash
curl -s https://packagecloud.io/geqo/togram/gpgkey | sudo apt-key add -
echo "deb https://packagecloud.io/geqo/togram/any/ any main" \
  | sudo tee /etc/apt/sources.list.d/togram.list
sudo apt update && sudo apt install togram
```

### rpm (Fedora / RHEL)

```bash
curl -s https://packagecloud.io/geqo/togram/gpgkey | sudo rpm --import /dev/stdin
echo "[togram]
name=togram
baseurl=https://packagecloud.io/geqo/togram/el/9/x86_64
enabled=1
gpgcheck=1" | sudo tee /etc/yum.repos.d/togram.repo
sudo dnf install togram
```

### Binary

Download the latest release from [Releases](https://github.com/geqo/togram/releases) and place the binary in your `$PATH`.

### Go

```bash
go install github.com/geqo/togram@latest
```

## Configuration

Create `/etc/togram/config`:

```ini
token = 123456:ABC-DEF...
chat  = @mychannel
```

Get a bot token from [@BotFather](https://t.me/BotFather). Find your chat ID via [@userinfobot](https://t.me/userinfobot).

## Usage

```
togram [flags] [file]

Flags:
  -c, --chat string    chat ID or @username (overrides config)
      --token string   bot token (overrides config)
  -h, --help           help for togram
```

### Examples

```bash
# pipe text
echo "hello" | togram

# pipe a file — type detected automatically
cat photo.jpg | togram
cat video.mp4 | togram

# send a file by path
togram document.pdf

# override chat and token inline
togram -c @otherchat --token 123:ABC archive.zip
```

### Content type detection

| Input | Telegram method |
|---|---|
| plain text ≤ 4096 chars | `sendMessage` |
| plain text > 4096 chars | `sendDocument` (as `.txt`) |
| image (jpg, png, gif, webp…) | `sendPhoto` |
| video (mp4, mkv, mov…) | `sendVideo` |
| audio (mp3, ogg, flac…) | `sendAudio` |
| anything else | `sendDocument` |

For stdin, type is detected from the first 512 bytes (magic bytes). For named files, the extension takes precedence.

## License

MIT
