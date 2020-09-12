package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/yunomu/kif"
	"github.com/yunomu/kif/ptypes"
)

var (
	inFile  = flag.String("f", "", "Input file")
	outFile = flag.String("o", "", "Output file")
	format  = flag.String("fmt", "", `Input/Output format
	s/S: kif (ShiftJIS) (default)
	u/U: kif (UTF8)
	j/J: Protocol Buffer (JSON)
	b/B: Protocol Buffer (byte strings)
`)
)

func init() {
	flag.Parse()
	log.SetOutput(os.Stderr)
}

func sjisRead(in io.Reader) (*ptypes.Kif, error) {
	return kif.Parse(in)
}

func sjisWrite(out io.Writer, k *ptypes.Kif) error {
	kifWriter := kif.NewWriter()
	return kifWriter.Write(out, k)
}

func jsonRead(in io.Reader) (*ptypes.Kif, error) {
	unmarshaler := &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	kif := &ptypes.Kif{}
	if err := unmarshaler.Unmarshal(in, kif); err != nil {
		return nil, err
	}
	return kif, nil
}

func jsonWrite(out io.Writer, kif *ptypes.Kif) error {
	marshaler := &jsonpb.Marshaler{
		Indent:       "  ",
		EmitDefaults: true,
	}
	return marshaler.Marshal(out, kif)
}

func binRead(in io.Reader) (*ptypes.Kif, error) {
	bs, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	kif := &ptypes.Kif{}
	if err := proto.Unmarshal(bs, kif); err != nil {
		return nil, err
	}
	return kif, nil
}

func binWrite(out io.Writer, kif *ptypes.Kif) error {
	bs, err := proto.Marshal(kif)
	if err != nil {
		return err
	}
	_, err = out.Write(bs)
	return err
}

func parseFormat(fmt string) (
	read func(io.Reader) (*ptypes.Kif, error),
	write func(io.Writer, *ptypes.Kif) error,
) {
	read = sjisRead
	write = sjisWrite
	for _, r := range []rune(fmt) {
		switch r {
		case 's':
			read = sjisRead
		case 'S':
			write = sjisWrite
		case 'u':
			read = func(in io.Reader) (*ptypes.Kif, error) {
				return kif.Parse(in, kif.ParseEncodingUTF8())
			}
		case 'U':
			write = func(out io.Writer, k *ptypes.Kif) error {
				kifWriter := kif.NewWriter(kif.WriteEncodingUTF8())
				return kifWriter.Write(out, k)
			}
		case 'j':
			read = jsonRead
		case 'J':
			write = jsonWrite
		case 'b':
			read = binRead
		case 'B':
			write = binWrite
		}
	}
	return
}

func main() {
	var in io.Reader = os.Stdin
	if *inFile != "" {
		f, err := os.Open(*inFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		in = f
	}

	var out io.Writer = os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		out = f
	}

	read, write := parseFormat(*format)

	kif, err := read(in)
	if err != nil {
		log.Fatalln(err)
	}

	if err := write(out, kif); err != nil {
		log.Fatalln(err)
	}
}
