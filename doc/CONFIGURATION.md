## Configurations

The workflow provides some configurations via [workflow environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/).

### Open URL Modifier Key

`TLDR_MOD_KEY_OPEN_URL` is variable that configures to open detail command url with default web browser.

By default, you can open url with pressing `cmd` + `enter` if a specified command has related URL.

Available variables are the followings

- `cmd`
- `ctrl`
- `alt`
- `fn`
- `shift`

For example, if you specified `ctrl` for `TLDR_MOD_KEY_OPEN_URL`, you can open url with pressing `ctrl` + `enter`.

### Command Format

`TLDR_COMMAND_FORMAT` is variable that switches command output format for user input parameters.

The workflow adopts `single` for default value.

By default, user input parameters are quoted by `{}`.

```
tar czf {target.tar.gz} --directory={path/to/directory} .
```

When the value is `uppercase`, user input parameters are converted to uppercase.

```
tar czf TARGET.TAR.GZ --directory=PATH/TO/DIRECTORY .
```

When the value is `original`, user input parameters are quoted by `{{}}`.

```
tar czf {{target.tar.gz}} --directory={{path/to/directory}} .
```

When the value is `remove`, user input parameters are not quoted by any words.

```
tar czf target.tar.gz --directory=path/to/directory .
```

Note:

Several commands are using uppercase in a place where no user input parameters.
For instance, `LISTEN` of `lsof` command is not user input parameter, but upppercase.
If you care, I recommend using `single` or `original`.

```
lsof -iTCP:{{port}} -sTCP:LISTEN
```

### Recommendations

The workflow shows update recommendations when tldr database is older or newer workflow is available.

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
