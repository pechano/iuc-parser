package main

import (
	"encoding/xml"
)
type transferInstructions struct {
	from string
	to string
}
type fileInfo struct{
	UUID string
	subtype string
	MD5 string
	RealName string
	BPRFolder string
	Matched bool 
}

type Info struct{
	XMLName xml.Name `xml:"manifest"`
	Date string `xml:"general-information>created"`
}

type dossier struct{
	XMLName xml.Name `xml:"manifest"`
	Attachment []attachment `xml:"contained-documents>attachment"`
	Document []document `xml:"contained-documents>document"`

}

type attachment struct {
	MD5Filename attachLink `xml:"linked-attachments>linked-doc"`
	RealFilename string `xml:"name"`
	Container string `xml:"container-uuid"`
}

type attachLink struct {
	LinkedDoc string `xml:"href,attr"`
}
type document struct{
	DocType string `xml:"type"`
	Container string   `xml:"uuid"`
	Category string   `xml:"subtype"`
}

type legislationKey struct{
	XMLkey string
	section string
}

