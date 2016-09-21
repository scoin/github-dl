Github-DL
==========

Github-DL is a command line search and clone tool of Github repositories, using the open Github API.

It is written in Go and has been compiled and tested to work on Linux x64 an Linux ARMv6. See releases for binaries.

###Usage

```shell
github-dl [-search <repository>] [-in <name,description,readme>] [-user <user>] [-lang <language>] [-stars <min..max>] [-size <min..max>] [-showforks <true/only>] [-sort <field>] [-order <asc/desc>]
```

To run the code locally:

```shell
git clone
go get github.com/nsf/termbox-go
go run githubl.go [parameters]
```