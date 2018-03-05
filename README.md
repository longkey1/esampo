# esampo

It is a tool to open a daily report of a specific day at once with a browser in [esa.io](https://esa.io/).

![esampo](https://user-images.githubusercontent.com/58566/36960482-929094ea-208a-11e8-9958-975926d1b34c.gif)

## Usage

Open yesterday's daily reports.

```shell
$ esampo
```

Open daily reports of 3 days ago.

```shell
$ esampo -b 3
```

## Installation 

You can download binary from [release page](https://github.com/longkey1/esampo/releases).

## Configuration

```toml
# $HOME/.esamporc

access_token = "your access token"
team_name = "your team name"
my_screen_name = "your screen name"
path = "日報/2006/01/02"
```

`path` is using [golang's time format](https://golang.org/src/time/format.go).
