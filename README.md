# hxd

hxd prints a file as a hex dump. Each line shows 16 bytes: the offset, the bytes in hex, and the printable ASCII characters.

```
$ hxd -f game.exe -s 0x40 -n 48
0x00000040: 0e 1f ba 0e | 00 b4 09 cd | 21 b8 01 4c | cd 21 54 68 | ........!..L.!Th
0x00000050: 69 73 20 70 | 72 6f 67 72 | 61 6d 20 63 | 61 6e 6e 6f | is program canno
0x00000060: 74 20 62 65 | 20 72 75 6e | 20 69 6e 20 | 44 4f 53 20 | t be run in DOS
```

## Install

Requires Go 1.27 or later.

```
go install github.com/KalebHawkins/hxd@latest
```

Or build from source. `make` runs `go fmt` and `go vet`, then writes the binary to `bin/`.

```
git clone https://github.com/KalebHawkins/hxd.git
cd hxd
make
```

## Usage

```
hxd -f <file> [-s offset] [-n length]
```

| Flag | Description |
|---|---|
| `-f`, `--file` | File to dump. Required. |
| `-s`, `--offset` | Byte offset to start at. Accepts decimal or prefixed values: `64`, `0x40`, `0o100`, `0b1000000`. Defaults to 0. |
| `-n`, `--length` | Maximum number of bytes to dump. `0` dumps to the end of the file. |

The offset must not be negative. An offset past the end of the file prints nothing. A length larger than the rest of the file stops at the end of the file.

## Examples

Dump the first 64 bytes of an executable (the DOS and PE headers start here):

```
hxd -f game.exe -n 64
```

Dump 64 bytes starting at offset 0x40:

```
hxd -f game.exe -s 0x40 -n 64
```

Save a dump of a save file, change something in the game, dump it again, and compare the two to find which bytes changed:

```
hxd -f save.dat > before.txt
hxd -f save.dat > after.txt
diff before.txt after.txt
```

Every line has the same width and starts with its offset, so line-based tools like `diff` and `grep` work on the output.

## Output format

```
0x000000c0: 69 6e 64 69 | 72 65 63 74 | 0a 29 0a    |             | indirect.).
```

- The offset is the position of the line's first byte in the file, in hex. With `-s`, offsets show the real file position, not a count from zero.
- Bytes are shown in lowercase hex, in groups of four separated by `|`.
- The ASCII column shows bytes `0x20` through `0x7e` as characters. Everything else, including newlines, is shown as `.`.
- A short last line is padded so the ASCII column stays aligned.

The dump is the only thing written to stdout. Errors go to stderr, and the exit code is 1 when something fails.

## License

Apache 2.0. See [LICENSE](LICENSE).
