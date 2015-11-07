package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/brettlangdon/jsonstream"
)

var args struct {
	File string   `arg:"-f,help:JSON stream file to read from"`
	TSV  bool     `arg:"-t,help:Reformat the JSON stream to TSV '<value>\t<value>'"`
	Key  bool     `arg:"-k,help:Whether or not to include the key in --tsv. '<key>=<value>\t<key>=<value>'"`
	Keys []string `arg:"positional,help:Which keys from the input JSON stream to include in the output"`
}

func init() {
	arg.MustParse(&args)
}

func getReader() (reader *jsonstream.Reader, err error) {
	var input io.Reader
	input = os.Stdin
	if args.File != "" {
		input, err = os.Open(args.File)
	}

	if err == nil {
		reader = jsonstream.NewReader(input, args.Keys)
	}
	return reader, err
}

func getFormatter() (formatter *jsonstream.Formatter, err error) {
	var format jsonstream.FormatType
	format = jsonstream.FormatJSON

	if args.TSV {
		format = jsonstream.FormatTSV
		if args.Key {
			format = jsonstream.FormatTSVKey
		}
	}

	if err == nil {
		formatter = jsonstream.NewFormatter(format)
	}
	return formatter, err
}

func main() {
	var err error
	var reader *jsonstream.Reader
	var formatter *jsonstream.Formatter

	reader, err = getReader()
	if err != nil {
		panic(err)
	}
	formatter, err = getFormatter()
	if err != nil {
		panic(err)
	}

	for {
		data, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		var output []byte
		output, err = formatter.Format(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\r\n", output)
	}
}
