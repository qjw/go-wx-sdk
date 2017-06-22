package utils

import "encoding/xml"

type CharData string
func (n CharData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct{
		S string `xml:",innerxml"`
	}{
		S: "<![CDATA[" + string(n) + "]]>",
	}, start)
}
