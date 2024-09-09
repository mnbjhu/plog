# Pretty Logger (plog)

A TUI to watch application logs in the terminal

## Usage

```
plog <command>
```

E.g.

```bash
plog java -jar app.jar
```

## Install

### Using Go Package Manager

```bash
go install github.com/mnbjhu/plog@latest
```

### Download The Binary

Binaries for each target can be found on the releases page [here](https://github.com/mnbjhu/plog/releases/latest).

## Custom Logging Format

You can use the`init` command to initialze the project.
You'll be prompted to enter details about your specific log format.

```bash
plog init
```

The cofiguration will be saved to the working directory under the name `.plog.json`.
You can edit this configuration manually for more contol over the log matching.
