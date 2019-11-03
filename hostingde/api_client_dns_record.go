package hostingde

// Adapted from hostingde provider in https://github.com/go-acme/lego

import (
	"errors"
	"fmt"
)

// https://www.hosting.de/api/?json#list-recordconfigs
func (d *Client) listRecords(findRequest RecordsFindRequest) (*RecordsFindResponse, error) {
	uri := defaultBaseURL + "/recordsFind"

	findResponse := &RecordsFindResponse{}

	rawResp, err := d.post(uri, findRequest, findResponse)
	if err != nil {
		return nil, err
	}

	if len(findResponse.Response.Data) == 0 {
		return nil, fmt.Errorf("%v: %s", err, toUnreadableBodyMessage(uri, rawResp))
	}

	if findResponse.Status != "success" {
		return findResponse, errors.New(toUnreadableBodyMessage(uri, rawResp))
	}

	return findResponse, nil
}

func (d *Client) getRecord(findRequest RecordsFindRequest) (*DNSRecord, error) {
	var record *DNSRecord

	findResponse, err := d.listRecords(findRequest)
	if err != nil {
		return nil, err
	}

	record = &findResponse.Response.Data[0]

	return record, nil
}
