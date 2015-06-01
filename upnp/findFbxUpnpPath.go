package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)


// XML Structures for service answer
type BrowseResponse struct {
	XMLName        xml.Name `xml:"BrowseResponse"`
	Result         string   `xml:"Result"`
	NumberReturned int      `xml:"NumberReturned"`
	TotalMatches   int      `xml:"TotalMatches"`
	UpdateID       int      `xml:"UpdateID"`
}

type Body struct {
	XMLName    xml.Name       `xml:"Body"`
	BrResponse BrowseResponse `xml:"BrowseResponse"`
}

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Bdy     Body     `xml:"Body"`
}

// XML Structures for request answer
type Container struct {
	XMLName xml.Name `xml:"container"`
	Title string `xml:"title"`
	Class string `xml:"class"`
	Id string `xml:"id,attr"`
}

type DIDLLite struct {
	XMLName xml.Name `xml:"DIDL-Lite"`
	Containers []Container `xml:"container"`
}

// Functions
func extractResultFromXmlResponse(response string) string {
	var e Envelope
	xml.Unmarshal([]byte(response), &e)
	return e.Bdy.BrResponse.Result
}

func findResourcePathInServiceResponse(response string) string {
	var d DIDLLite
	var id string
	xml.Unmarshal([]byte(response), &d)
        for _,container := range d.Containers {
		if container.Title == "Freebox TV" {
			id = container.Id
		}
	}

	return id
}

// Main
func main() {
	url := "http://192.168.0.254:52424/service/ContentDirectory/control"

	// Our request
	soap_data := []byte(`
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<s:Body><u:Browse xmlns:u="urn:schemas-upnp-org:service:ContentDirectory:1">
<ObjectID>0/0</ObjectID>
<BrowseFlag>BrowseDirectChildren</BrowseFlag>
<Filter>id,dc:title,res,sec:CaptionInfo,sec:CaptionInfoEx,pv:subtitlefile</Filter>
<StartingIndex>0</StartingIndex>
<RequestedCount>0</RequestedCount>
<SortCriteria></SortCriteria>
</u:Browse>
</s:Body>
</s:Envelope>
`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(soap_data))
	req.Header.Set("SOAPACTION", "urn:schemas-upnp-org:service:ContentDirectory:1#Browse")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	result := extractResultFromXmlResponse(string(body))
        fmt.Println(findResourcePathInServiceResponse(result))
}
