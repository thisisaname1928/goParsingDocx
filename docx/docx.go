package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"strings"
)

type Property struct {
	XMLName xml.Name
	Val     *string `xml:"val,attr"`
	Fill    *string `xml:"fill,attr"`
}

type Text struct {
	XMLName xml.Name `xml:"t"`
	Text    *string  `xml:",chardata"`
}

type RunProperties struct {
	XMLName    xml.Name    `xml:"rPr"`
	Properties *[]Property `xml:",any"`
}

type Run struct {
	XMLName       xml.Name      `xml:"r"`
	Text          Text          `xml:"t"`
	RunProperties RunProperties `xml:"rPr"`
	Drawing       *[]Drawing    `xml:"drawing"`
}

type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Runs    []Run    `xml:"r"`
}

type Body struct {
	XMLName    xml.Name    `xml:"body"`
	Paragraphs []Paragraph `xml:"p"`
}

type Document struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

func parseSubPropery(p *Property) bool {
	if p != nil {

		if p.Val != nil { // when property tag have w:val
			// some exception
			if p.Fill != nil {
				switch *p.Fill {
				case "auto":
					return false
				}
			}

			if *p.Val != "false" {
				return true
			} else {
				return false
			}
		} else { // when it doesn't have
			return true
		}

	}
	// if property isn't exist
	return false
}

func ParseProperties(p *[]Property, text *string) {
	if p != nil {
		for _, prop := range *p {
			if parseSubPropery(&prop) {
				switch prop.XMLName.Local {
				case "b":
					*text = "<b>" + *text + "</b>"
				case "i":
					*text = "<i>" + *text + "</i>"
				case "u":
					*text = "<u>" + *text + "</u>"
				case "shd":
					*text = "<mark>" + *text + "</mark>"
				}
			}
		}
	}
}

func Parse(path string) (Document, error) {
	var doc Document
	invalidDocxFile := errors.New("invalid docx file")

	documentXML, e := DecompressFile(path, "word/document.xml")
	if e != nil || len(documentXML) == 0 {
		return doc, invalidDocxFile
	}

	decoder := xml.NewDecoder(bytes.NewReader(documentXML))
	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "p" {
				var p Paragraph
				if err := decoder.DecodeElement(&p, &se); err == nil {
					doc.Body.Paragraphs = append(doc.Body.Paragraphs, p)
				}
			}
		}
	}

	return doc, nil
}

func Parse2Html() string {

	var doc Document
	doc, e := Parse("./test.docx")
	rIDTable := GetRID("./test.docx")

	if e != nil {
		panic(e)
	}

	p := doc.Body.Paragraphs

	output := ""
	for _, v := range p {
		r := v.Runs
		hastext := false
		for _, i := range r {
			if i.Drawing != nil {
				hastext = true
				for _, drawing := range *i.Drawing {
					output += "<img src=\"./" + rIDTable[ParseDrawing(&drawing)] + "\">"
				}
			} else if i.Text.Text != nil {
				hastext = true
				text := i.Text.Text
				ParseProperties(i.RunProperties.Properties, text)
				output += "<label>" + *text + "</label>"
			}
		}
		if hastext {
			output += "\n<br>\n"
		}
	}

	return "<body>\n" + output + "</body>"
}

func CheckIsMediaFile(f string) bool {
	res := strings.Split(f, "/")
	if res[0] == "word" && res[1] == "media" {
		return true
	}

	return false
}

func DecompressDocxMedia(path string, outpath string) error {

	f, e := os.ReadFile(path)

	if e != nil {
		return e
	}

	a, e := zip.NewReader(bytes.NewReader(f), int64(len(f)))

	if e != nil {
		return e
	}

	for _, f := range a.File {
		if CheckIsMediaFile(f.Name) {
			dat, e := f.Open()

			if e != nil {
				return e
			}
			defer dat.Close()

			tmp := strings.Split(f.Name, "/")
			fname := tmp[len(tmp)-1]

			f, e := os.Create(outpath + fname)
			if e != nil {
				return e
			}
			defer f.Close()

			if _, e := io.Copy(f, dat); e != nil {
				return e
			}
		}
	}

	return nil
}

func DecompressFile(path string, fname string) ([]byte, error) {
	f, e := os.ReadFile(path)

	if e != nil {
		return []byte(""), e
	}

	arch, e := zip.NewReader(bytes.NewReader(f), int64(len(f)))

	if e != nil {
		return []byte(""), e
	}

	for _, f := range arch.File {
		if f.Name == fname {
			dat, e := f.Open()
			if e != nil {
				return []byte(""), e
			}

			defer dat.Close()

			buf := new(bytes.Buffer)

			if _, e := io.Copy(buf, dat); e != nil {
				return []byte(""), e
			}

			return buf.Bytes(), nil

		}
	}

	return []byte(""), nil
}
