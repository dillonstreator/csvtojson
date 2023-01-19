package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"github.com/dillonstreator/csvtojson"
)

func main() {
	separatorPtr := flag.String("separator", ",", "csv separator")
	mappingPtr := flag.String("mapping", "", "field mapping json object where keys are incoming headers and values are the output keys")
	injectPtr := flag.String("inject", "", "static key/value(s) to inject into resulting json objects")
	filepathPtr := flag.String("filepath", "", "path to csv file. must use either this or pipe csv via standard input.")
	timeoutDurationPtr := flag.Duration("timeout", 0, "exit program after specified timeout duration")

	flag.Parse()

	separator, s := utf8.DecodeRuneInString(*separatorPtr)
	if separator == utf8.RuneError || s == 0 {
		log.Fatal("invalid separator")
	}

	var mapping map[string]string
	if len(*mappingPtr) > 0 {
		err := json.Unmarshal([]byte(*mappingPtr), &mapping)
		if err != nil {
			log.Fatal(err)
		}
	}

	var inject map[string]string
	if len(*injectPtr) > 0 {
		err := json.Unmarshal([]byte(*injectPtr), &inject)
		if err != nil {
			log.Fatal(err)
		}
	}

	var hasPipedInput bool

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		hasPipedInput = true
	}

	var r io.Reader

	if hasPipedInput {
		r = os.Stdin
	} else {
		if len(*filepathPtr) == 0 {
			log.Fatal("filepath required if no csv content is piped into standard input")
		}

		file, err := os.Open(*filepathPtr)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		r = file
	}

	ctx := context.Background()

	if *timeoutDurationPtr != 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(context.Background(), *timeoutDurationPtr)
		defer cancel()
	}

	err := csvtojson.Copy(
		ctx,
		r,
		os.Stdout,
		csvtojson.WithFields(mapping),
		csvtojson.WithInject(inject),
		csvtojson.WithSeparator(separator),
	)
	if err != nil {
		log.Fatal(err)
	}
}
