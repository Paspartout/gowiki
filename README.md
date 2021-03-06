# Gowiki

[![Build Status](https://travis-ci.org/Paspartout/gowiki.svg?branch=master)](https://travis-ci.org/Paspartout/gowiki)

![Screenshot](https://paspartout.github.io/gowiki/screenshot.png)

Gowiki started out as an extensions to the 
[golang.org web application](https://golang.org/doc/articles/wiki/) tutorial.
The idea behind it is to implement a simple file based wiki web application
that can serve wikis from [vimwiki](https://github.com/vimwiki/vimwiki) in
the markdown format.

The use case of the current version is browsing and editing a wiki that
consists of markdown files in a single directory.

You may take a look at the [todo.md](https://github.com/Paspartout/gowiki/blob/master/todo.md)
for what I implented and what I may implement in future versions.

## Installing

To install gowiki manually from source type the following commands in your terminal:

```sh
$ go get -d github.com/Paspartout/gowiki
$ go get github.com/mjibson/esc # for embedding resource files into executable
$ cd $GOPATH/src/github.com/Paspartout/gowiki
$ go generate
$ go install
```

Alternatively you can download the latest realease from the [Github Releases](https://github.com/Paspartout/gowiki/releases).

## License

Gowiki itself is licensed under the MIT License.
See the LICENSE file.

Gowiki is using the [bulma.css framework](https://bulma.io/) and the [openiconic icons](https://useiconic.com/open).
Their licenses can be found in the static directory.

