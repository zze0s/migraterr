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

Required flags:
- `--glob`  Glob to file(s) like `~/.sessions/*.torrent.rtorrent`
- `--replacements`  Array of `oldvalue|newvalue` similar to `sed`.

Optional flags:
- `--dry-run`  Do not edit any data
- `--verbose`  Verbose output
- `--export /new/dir`   Export edited files to new directory:
- `--export /new/dir`   Export edited files to new directory:
