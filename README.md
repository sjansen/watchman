# Go Watchman Client

[![GoDoc](https://godoc.org/github.com/sjansen/watchman?status.svg)](https://godoc.org/github.com/sjansen/watchman)
[![Go Report Card](https://goreportcard.com/badge/github.com/sjansen/watchman)](https://goreportcard.com/report/github.com/sjansen/watchman)
[![Latest Release](https://img.shields.io/github/release/sjansen/watchman.svg)](https://github.com/sjansen/watchman/releases)
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![License](https://img.shields.io/github/license/sjansen/watchman.svg)](https://github.com/sjansen/watchman/blob/master/LICENSE)

This Go module provides a client for Facebook's [Watchman](https://facebook.github.io/watchman/), a file watching service with support for recording and querying filesystem state. Two alternatives are provided. The [high-level API](https://godoc.org/github.com/sjansen/watchman) is designed for common use cases. The [low-level API](https://godoc.org/github.com/sjansen/watchman/protocol) is designed to enable advanced use cases.


## Frequently Asked Questions

**Why use Watchman instead of fsnotify/inotify/kevent/etc?**

- Watchman conserves system resources by enabling applications to share watches instead of duplicating them.
- Watchman waits for changes to settle down before sending notifications.
- Watchman can also be queried for changes since a previous check, or queried for current filesystem state.
- Watchman has been battle hardened by use in projects such as Buck, Mercurial, and React Native.

**Are all Watchman features supported?**

All primitives necessary to access the full Watchman protocol are implemented, however this project is still a work in progress. Most Watchman commands still need to be mapped to more friendly  data structures and methods. In addition, the eventual replacement of JSON ecoding with more efficient [BSER](https://facebook.github.io/watchman/docs/bser.html) encoding is planned but not yet implemented.

For details, see [docs/status.md](docs/status.md).


## Roadmap

This is a personal project. I work on it when I feel like it.
The following is a list of things I might work on. It is not
a list of commitments.

Contributions are welcome.

### 0.2
 - improved test coverage
 - start testing on Windows
 - expose existing command options

### 0.3
 - support additional commands


## Additional Resources

https://facebook.github.io/watchman/docs/install.html

https://www.facebook.com/notes/facebook-engineering/watchman-faster-builds-with-large-source-trees/10151457195103920/
