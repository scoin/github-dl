Github-DL
==========
![Demo](http://i.imgur.com/MAd6RE1.gif)


Github-DL is a command line search and clone tool of Github repositories, using the open [Github API Search Repositories Endpoint](https://developer.github.com/v3/search/#search-repositories).

It is written in Go and has been compiled and tested to work on Linux x64 an Linux ARMv6. See releases for binaries.

While this is not a replacement for github.com or one of many other clients, this is intended for headless and IoT devices, like the [PocketChip](https://getchip.com/pages/pocketchip) and to save you from typing long URLs / SCP commands etc.

###Usage

```shell
github-dl [-search <terms>] [-in <name,description,readme>] [-user <user>] [-lang <language>] [-stars <min..max>] [-size <min..max>] [-showforks <true/only>] [-sort <stars/forks/updated>] [-order <asc/desc>]
```

To run the code locally:

```shell
git clone https://github.com/scoin/github-dl.git
cd github-dl/
export GOPATH=$PWD
go get github.com/nsf/termbox-go
go run githubl.go [parameters]
```
