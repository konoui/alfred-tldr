![Generic badge](https://github.com/konoui/tldr/workflows/test/badge.svg)

## alfred tldr
tldr alfred workflow written in go.

## Install
- Download the workflow form [latest release](https://github.com/konoui/tldr/releases).
- Build the workflow on your computer.
```
$ make package
$ ls
tldr.alfredworkflow (snip)
```

## Usage
`tldr <query>`

Options   
`-u` option updates command list (tldr repository).   
`-p` option chooses platform from `linux`,`osx`,`sunos`,`windows`.   

![alfred-tldr](./alfred-tldr.png)

## License
MIT License.
