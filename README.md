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
`-u` option updates update local database (tldr repository).  
`-p` option selects platform from `linux`,`osx`,`sunos`,`windows`.  

![alfred-tldr](./alfred-tldr.png)

## License
MIT License.
