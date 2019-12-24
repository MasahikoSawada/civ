# civ

A simple CSV interactive viewer.

# Demo

![civ demo](https://github.com/MasahikoSawada/masahikosawada.github.io/raw/master/images/civ.gif)

# Build

```
go get -u github.com/MasahikoSawada/civ
```

# Usage

```
$ civ [options] [FILE]
```

| Option      | Description                                       |
|:------------|:--------------------------------------------------|
| `-d string` | Use `string` as a delimiter instead of comma(`,`) |
|`-H`         | Set dummy header (col_1, col_2 ...)|

* civ reads data from stdin if no file is specified.
* civ processes the first line as a header line by default. If the first line of the file is not header line please use `-H` option to set dummy headers.
* `-d` option allows a speciial argument `\t` to parse TSV.

# Query Buffer

civ has a buffer for user-input query at top of the window. The first character indicates the current mode as described below.

# Modes

civ has 4 modes: view mode, command mode, search mode and filter mode.

You can swtich modes by special character when the query buffer is empty.

* '`:`' : View Mode
* '`@`' : Command Mode
* '`/`' : Search Mode
* '`^`' : Filter Mode

Press `Ctrl-g` always clear all query buffer and switch to view mode.

Press `Ctrl-c` exits but executing `@exit` also exits while output the table data to `stdout`.

Press `Enter` saves the result of the current command (at most one result for each searching and filtering).

## View Mode(`:`)

Viewing the table data with the following ''less-like'' key binds:

|Key|Description|
|:---|:-----------|
|e|Forward one line|
|y|Backward one line|
|f, SPACE|Forward one window|
|b|Forward one window|
|d|Forward one half-window|
|u|Backward one half-window|
|g|Go to first line in file|
|G|Go to last line in file|

## Command Mode(`@`)

Executing the following commands modify the table:

|Command|Description|
|:------|:----------|
|`@hide column-name [...]`|Hide the specified column(s)|
|`@show column-name [...]`|Show the specified hidden column(s)|
|`@show_only column-name [...]`|Show only the specified column(s)|
|`@reset`|Reset all configurations(row filtering, column visibility etc)|
|`@exit`|Output the table data to stdout and exit normally|

Pressing `Enter` key executes the input command.

Note that specifying column name is case-insensive.

## Search Mode(`/`)

civ supports the incremental search hight-lighting the matched words.

## Filter Mode(`^`)

civ supports the incremental filtering rows.

## Limitation

Since civ is continually being improved it has some limitations. These limitation might be resolved in the future.

* Not supports multi-byte characters.
