jsonstream
==========
[![GoDoc](https://godoc.org/github.com/brettlangdon/jsonstream?status.svg)](https://godoc.org/github.com/brettlangdon/jsonstream)

`jsonstream` is a utility for interacting with a newline delimited JSON stream (e.g. log files which log as JSON).

The goal of `jsonstream` is to help convert a stream of JSON data into a format more friendly to typical UNIX tools.

## Quick start
To install `jsonstream` run `go get github.com/brettlangdon/jsonstream/cmd/...`

Given an example log file `example.log`:
```
{"line": 1, "data": {"type": "string", "value": "Hello World"}, "extra": {"user_id": 12345} }
{"line": 2, "data": {"type": "bool", "value": true}, "extra": {"user_id": 12346} }
{"line": 3, "data": {"type": "integer", "value": 56}, "extra": {} }
```

We can then pipe this file into `jsonstream`:

```
$ cat example.log | jsonstream
{"data":{"type":"string","value":"Hello World"},"extra":{"user_id":12345},"line":1}
{"data":{"type":"bool","value":true},"extra":{"user_id":12346},"line":2}
{"data":{"type":"integer","value":56},"extra":{},"line":3}
```

By default `jsonstream` uses a JSON formatter for output which just sorts the keys of the JSON input.

However, we can use the alternate TSV format to produce a more friendly output:

```
$ cat example.log | jsonstream --tsv
{"type":"string","value":"Hello World"}	{"user_id":12345}	1
{"type":"bool","value":true}	{"user_id":12346}	2
{"type":"integer","value":56}	{}	3
```

This TSV format allows us to more easily pipe this data into other UNIX tools, like `awk`.

For example, if we wanted to use `awk` to print just the `extra` field we could do:
```
$ cat example.log| jsonstream --tsv |  awk -F '\t' '{print $2}'
{"user_id":12345}
{"user_id":12346}
{}
```

We can also pipe `jsonstream` calls together to extract nested data:

```
$ cat example.log | jsonstream --tsv data | jsonstream --tsv --key
type="string"	value="Hello World"
type="bool"	value=true
type="integer"	value=56
```


## Installation
`jsonstream` requires `GO15VENDOREXPERIMENT="1"`.

You can install via `go get`:

```
go get github.com/brettlangdon/jsonstream/cmd/...
```

Or install from source:

```
git clone git://github.com/brettlangdon/jsonstream
cd ./jsonstream
make
```

## Usage
```
$ jsonstream --help
usage: jsonstream [--file FILE] [--tsv] [--key] [KEYS [KEYS ...]]

positional arguments:
  keys                   Limit the output format to only include the listed keys

options:
  --file FILE, -f FILE   JSON stream file to read from
  --tsv, -t              Reformat the JSON stream to TSV '<value>	<value>'
  --key, -k              Whether or not to include the key in --tsv. '<key>=<value>	<key>=<value>'
```

### Supplying a file
`jsonstream` supports either piping data in via stdin, or else you can supply a single file with the `--file FILE` command switch.

```
$ cat example.log | jsonstream
```

```
$ jsonstream --file example.log
$ jsonstream -f example.log
```

### Output formats
`jsonstream` currently supports 3 output formats: `JSON`, `TSV`, and `TSV` with keys.

#### JSON (default)
This format will still parse the input JSON to ensure it is valid, sort the keys in ascending alphabetical order (to ensure consistency) and then output as a JSON object on a single line.

**example.log**
```
{"line": 1, "data": {"type": "string", "value": "Hello World"}, "extra": {"user_id": 12345} }
```

```
$ cat example.log | jsonstream
{"data":{"type":"string","value":"Hello World"},"extra":{"user_id":12345},"line":1}
```

#### TSV
This format will parse the input JSON to ensure it is valid, sort the keys in ascending alphabetical order (to ensure consistency) and then output the values of the top level keys as tab separated JSON encoded values.

**example.log**
```
{"line": 1, "data": {"type": "string", "value": "Hello World"}, "extra": {"user_id": 12345} }
```

```
$ cat example.log | jsonstream --tsv
{"type":"string","value":"Hello World"}	{"user_id":12345}	1
```

#### TSV with keys
This format is the same as `TSV` except that the output tab separated `<key>=<value>` pairs.


**example.log**
```
{"line": 1, "data": {"type": "string", "value": "Hello World"}, "extra": {"user_id": 12345} }
```

```
$ cat example.log | jsonstream --tsv --key
data={"type":"string","value":"Hello World"}	extra={"user_id":12345}	line=1
```

### Limiting keys
We can supply optional positional arguments to `jsonstream` specifying which top levels keys we want included in the output.

For example, given the following `example.log`:
```
{"line": 1, "data": {"type": "string", "value": "Hello World"}, "extra": {"user_id": 12345} }
```

If we only cared about the `data` property we can run:

```
$ cat example.log | jsonstream data
{"data":{"type":"string","value":"Hello World"}}
```

And if we wanted to extract just the `value` property from `data` we can do:

```
$ cat example.log | jsonstream --tsv data | jsonstream --tsv value
"Hello World"
```

## Alternatives
`jsonstream` is meant to be a very simple utility for transforming a stream of JSON into something that can be piped into another command.

If you are interested in something more feature rich, check out these alternatives:
* `jq` - https://stedolan.github.io/jq/
    * Supports processing stream of newline delimited JSON
* `underscore-cli` - https://github.com/ddopson/underscore-cli
    * Does not support newline delimited JSON as input (as of right now)
