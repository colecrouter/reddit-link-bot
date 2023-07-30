package media

import (
	"encoding/xml"
	"fmt"
)

type MPD struct {
	XMLName                   xml.Name `xml:"MPD"`
	Text                      string   `xml:",chardata"`
	Xmlns                     string   `xml:"xmlns,attr"`
	Xsi                       string   `xml:"xsi,attr"`
	MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr"`
	MinBufferTime             string   `xml:"minBufferTime,attr"`
	Profiles                  string   `xml:"profiles,attr"`
	Type                      string   `xml:"type,attr"`
	SchemaLocation            string   `xml:"schemaLocation,attr"`
	Period                    struct {
		Text          string `xml:",chardata"`
		Duration      string `xml:"duration,attr"`
		ID            string `xml:"id,attr"`
		AdaptationSet []struct {
			XMLName                 xml.Name           `xml:"AdaptationSet"`
			Text                    string             `xml:",chardata"`
			ContentType             string             `xml:"contentType,attr"`
			ID                      string             `xml:"id,attr"`
			SegmentAlignment        string             `xml:"segmentAlignment,attr"`
			StartWithSAP            string             `xml:"startWithSAP,attr"`
			SubsegmentAlignment     string             `xml:"subsegmentAlignment,attr"`
			SubsegmentStartsWithSAP string             `xml:"subsegmentStartsWithSAP,attr"`
			Representation          RepresentationList `xml:"Representation"`
		} `xml:"AdaptationSet"`
	} `xml:"Period"`
}

type RepresentationList []interface{}

type AudioRepresentation struct {
	XMLName                   xml.Name `xml:"Representation"`
	Text                      string   `xml:",chardata"`
	AudioSamplingRate         string   `xml:"audioSamplingRate,attr"`
	Bandwidth                 string   `xml:"bandwidth,attr"`
	Codecs                    string   `xml:"codecs,attr"`
	ID                        string   `xml:"id,attr"`
	MimeType                  string   `xml:"mimeType,attr"`
	AudioChannelConfiguration struct {
		Text        string `xml:",chardata"`
		SchemeIdUri string `xml:"schemeIdUri,attr"`
		Value       string `xml:"value,attr"`
	} `xml:"AudioChannelConfiguration"`
	BaseURL     string `xml:"BaseURL"`
	SegmentBase struct {
		Text           string `xml:",chardata"`
		IndexRange     string `xml:"indexRange,attr"`
		Timescale      string `xml:"timescale,attr"`
		Initialization struct {
			Text  string `xml:",chardata"`
			Range string `xml:"range,attr"`
		} `xml:"Initialization"`
	} `xml:"SegmentBase"`
}

type VideoRepresentation struct {
	XMLName     xml.Name `xml:"Representation"`
	Text        string   `xml:",chardata"`
	Bandwidth   string   `xml:"bandwidth,attr"`
	Codecs      string   `xml:"codecs,attr"`
	FrameRate   string   `xml:"frameRate,attr"`
	Height      string   `xml:"height,attr"`
	ID          string   `xml:"id,attr"`
	MimeType    string   `xml:"mimeType,attr"`
	Width       string   `xml:"width,attr"`
	BaseURL     string   `xml:"BaseURL"`
	SegmentBase struct {
		Text           string `xml:",chardata"`
		IndexRange     string `xml:"indexRange,attr"`
		Timescale      string `xml:"timescale,attr"`
		Initialization struct {
			Text  string `xml:",chardata"`
			Range string `xml:"range,attr"`
		} `xml:"Initialization"`
	} `xml:"SegmentBase"`
}

func (r *RepresentationList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var setType string
	for _, attr := range start.Attr {
		if attr.Name.Local == "mimeType" {
			setType = attr.Value
			break
		}
	}

	switch setType {
	case "audio/mp4":
		var v AudioRepresentation
		err = d.DecodeElement(&v, &start)
		*r = append(*r, v)
	case "video/mp4":
		var v VideoRepresentation
		err = d.DecodeElement(&v, &start)
		*r = append(*r, v)
	default:
		return fmt.Errorf("unknown contentType: %s", setType)
	}

	return err
}

func (m *MPD) GetMediaLinks() (audio string, video string, err error) {
	for _, set := range m.Period.AdaptationSet {
		for _, rep := range set.Representation {
			switch v := rep.(type) {
			case AudioRepresentation:
				audio = v.BaseURL
			case VideoRepresentation:
				video = v.BaseURL
			default:
				err = fmt.Errorf("unknown representation type: %T", v)
			}
		}
	}

	return
}
