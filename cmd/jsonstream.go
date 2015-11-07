package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/brettlangdon/jsonstream"
)

var input io.Reader
var inputFile string

func init() {
	flag.StringVar(&inputFile, "file", nil, "")
}

func main() {
	reader := jsonstream.NewReader(os.Stdin)
	for {
		data, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		var output []byte
		output, err = json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\r\n", output)
	}
}
