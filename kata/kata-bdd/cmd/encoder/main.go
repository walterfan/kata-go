package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/waltfy/kata-bdd/encoder/pkg/encoder"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: encoder --mode [encode|decode] --type [base64|hex|url] --text <string>\n")
	flag.PrintDefaults()
}

func main() {
	var mode string
	var typ string
	var text string

	flag.StringVar(&mode, "mode", "encode", "mode: encode or decode")
	flag.StringVar(&typ, "type", "base64", "type: base64, hex, or url")
	flag.StringVar(&text, "text", "", "text to process")
	flag.Usage = usage
	flag.Parse()

	if text == "" {
		usage()
		os.Exit(2)
	}

	switch typ {
	case "base64":
		if mode == "encode" {
			fmt.Println(encoder.EncodeBase64(text))
			return
		}
		out, err := encoder.DecodeBase64(text)
		if err != nil { log.Fatal(err) }
		fmt.Println(out)
	case "hex":
		if mode == "encode" {
			fmt.Println(encoder.EncodeHex(text))
			return
		}
		out, err := encoder.DecodeHex(text)
		if err != nil { log.Fatal(err) }
		fmt.Println(out)
	case "url":
		if mode == "encode" {
			fmt.Println(encoder.EncodeURL(text))
			return
		}
		out, err := encoder.DecodeURL(text)
		if err != nil { log.Fatal(err) }
		fmt.Println(out)
	default:
		log.Fatalf("unknown type: %s", typ)
	}
}
