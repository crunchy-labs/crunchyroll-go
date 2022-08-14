# crunchyroll-go

A [Go](https://golang.org) library & cli for the undocumented [crunchyroll](https://www.crunchyroll.com) api. To use it, you need a crunchyroll premium account to for full (api) access.

<p align="center">
  <a href="https://github.com/ByteDream/crunchyroll-go">
    <img src="https://img.shields.io/github/languages/code-size/ByteDream/crunchyroll-go?style=flat-square" alt="Code size">
  </a>
  <a href="https://github.com/ByteDream/crunchyroll-go/releases/latest">
    <img src="https://img.shields.io/github/downloads/ByteDream/crunchyroll-go/total?style=flat-square" alt="Download Badge">
  </a>
  <a href="https://github.com/ByteDream/crunchyroll-go/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/ByteDream/crunchyroll-go?style=flat-square" alt="License">
  </a>
  <a href="https://golang.org">
    <img src="https://img.shields.io/github/go-mod/go-version/ByteDream/crunchyroll-go?style=flat-square" alt="Go version">
  </a>
  <a href="https://github.com/ByteDream/crunchyroll-go/releases/latest">
    <img src="https://img.shields.io/github/v/release/ByteDream/crunchyroll-go?style=flat-square" alt="Release">
  </a>
  <a href="https://discord.gg/gUWwekeNNg">
    <img src="https://img.shields.io/discord/915659846836162561?label=discord&style=flat-square" alt="Discord">
  </a>
  <a href="https://github.com/ByteDream/crunchyroll-go/actions/workflows/ci.yml">
    <img src="https://github.com/ByteDream/crunchyroll-go/workflows/CI/badge.svg?style=flat" alt="CI">
  </a>
</p>

**If you're searching for the command line client, head over to _[crunchy-cli](https://github.com/crunchy-labs/crunchy-cli)_. This repo only contains the golang crunchyroll library!**

> Beginning with version v3, this project is set to maintenance mode (only small fixes, no new features) in favor of rewriting it completely in Rust.
> Go bindings for the finished rust rewrite will be provided.

## 📥 Download

Download the library via `go get`

```shell
$ go get github.com/crunchy-labs/crunchyroll-go/v3
```

## 📚 Documentation

The documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/crunchy-labs/crunchyroll-go/v3).

## ☂️ Coverage

Around _90% - 95%_ of the crunchyroll beta api (at state of writing) are implemented.
Only some endpoints on the home / index screen are missing.
Since the library will be rewritten in Rust, I don't see any further use cases of implementing the missing endpoints in this project.
They would be useless for 99% of the library usage anyway, unless you want to rebuild crunchyroll on top of it (or for whatever reason you want to implement home screen endpoints).

## ⚖ License

This project is licensed under the GNU Lesser General Public License v3.0 (LGPL-3.0) - see the [LICENSE](LICENSE) file for more details.
