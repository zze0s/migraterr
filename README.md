# Migraterr

Utilities to help with migrating files between bt clients and machines.

## Commands

- bencode

### bencode

Edit and show bencoded files such as fast resume files from bt clients.

`migraterr bencode <subcommand>`

Available sub commands:
- `info`  Display info about file(s)
- `edit`  Edit data

#### edit

Usage:
```shell
migraterr bencode edit "~/.sessions/*.torrent.rtorrent" --replace "/old/path/|/new/path/" --dry-run --verbose
migraterr bencode edit "~/.sessions/*.torrent.rtorrent" --replace "/old/path/|/new/path/" --export "./edited-files" -v
```

Required flags:
- `--replace`  Array of `oldvalue|newvalue` similar to `sed`.

Optional flags:
- `--dry-run`  Do not edit any data
- `--verbose`  Verbose output
- `--export /new/dir`   Export edited files to new directory:

#### info

Usage:
```shell
migraterr bencode info "~/.sessions/*.torrent.rtorrent"
migraterr bencode edit "~/.sessions/123.torrent.rtorrent" 
```
