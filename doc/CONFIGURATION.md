## Configurations

This workflow provides several configurations through [workflow environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/).

### Open URL Modifier Key

The `TLDR_MOD_KEY_OPEN_URL` variable configures the keys to use in order to open the URL of a detail command with a web browser.

By default, you can open the URL by pressing cmd(⌘) + enter, if a specified command has a related URL.

The available variables are:

- `cmd`
- `ctrl`
- `alt`
- `fn`
- `shift`

For example, if you specify ctrl for `TLDR_MOD_KEY_OPEN_URL`, you can open the URL by pressing `ctrl(^)` + `enter`.

### Command Format

The `TLDR_COMMAND_FORMAT` variable switches the command output format for user input parameters.

The workflow adopts `single` as the default value.

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

Several commands use uppercase in places where no user input parameters are present.
For instance, `LISTEN` of the `lsof` command is not a user input parameter but is written in uppercase.
If you care, I recommend using `single` or `original`.

```
lsof -iTCP:{{port}} -sTCP:LISTEN
```

### Recommendations

This workflow shows update recommendations when the tldr database is out of date or when a newer version of the workflow is available.

#### Tldr Database Recommendation

The `TLDR_DB_UPDATE_RECOMMENDATION` variable enables or disables showing the update recommendation of the tldr database.
The value is `true` by default.
If the tldr database is older than `two weeks`, the workflow shows the update recommendation.

#### Newer Alfred Workflow Recommendation

The `TLDR_WORKFLOW_UPDATE_RECOMMENDATION` variable enables or disables showing the update recommendation if a newer version of the Alfred workflow is released.
The value is `true` by default.

The `TLDR_WORKFLOW_UPDATE_INTERVAL_DAYS` variable defines how frequently the workflow checks the remote git repository.
The value is `7` days by default.
The workflow checks for a new workflow version by accessing the remote git repository every `7` days.
It shows an update recommendation if a newer version is available.
