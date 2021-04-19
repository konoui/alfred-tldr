![Generic badge](https://github.com/konoui/alfred-tldr/workflows/test/badge.svg)
![Code Grade](https://www.code-inspector.com/project/20715/status/svg)

## alfred tldr

[tldr-pages](https://github.com/tldr-pages/tldr) alfred workflow written in go.

![alfred-tldr](./alfred-tldr.gif)

## Usage

`tldr <query>`

Options  
`--version`/`-v` option shows the current version of the client.  
`--update`/`-u` option updates local database (tldr repository).  
`--platform`/`-p` option selects platform from `linux`,`osx`,`sunos`,`windows`.  
`--language`/`-L` option selects preferred language for the page.

## Install

- Download and open the workflow with terminal for macOS.

```
$ curl -O -L https://github.com/konoui/alfred-tldr/releases/latest/download/tldr.alfredworkflow && open tldr.alfredworkflow
```

- Build the workflow on your computer.

```
$ make package
$ ls
tldr.alfredworkflow (snip)
```

## Configurations

The workflow shows update recommendations when tldr database is older or newer workflow is available.
update recommendations are configurable via [workflow environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/).

#### Tldr Database Recommendation

`TLDR_DB_UPDATE_RECOMMENDATION` is variable that enable or disable showing update recommendation of tldr database .
The value is `true` by default.
If tldr database is older than `two weeks`, the workflow shows it.

#### Newer Alfred Workflow Recommendation

`TLDR_WORKFLOW_UPDATE_RECOMMENDATION` is variable that enables or disables showing update recommendation if newer version alfred workflow is released.
The value is `true` by default.

`TLDR_WORKFLOW_UPDATE_INTERVAL_DAYS` is variable that defines how frequency the workflow checks remote git repository.
The value is `7` days by default.

The workflow checks new workflow version by accessing remote git repository per `7` days.
It shows update recommendation if newer version is available.

## License

MIT License.

## Special Thanks

Icons are provided by [takafumi-max](https://github.com/takafumi-max)
