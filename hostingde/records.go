package hostingde

import (
	"errors"
	"fmt"
	"net/http"
)

// https://www.hosting.de/api/?json#list-recordconfigs
func (d *Client) listRecords(findRequest RecordsFindRequest) (*RecordsFindResponse, error) {
	uri := defaultBaseURL + "/recordsFind"

	findResponse := &RecordsFindResponse{}

	rawResp, err := d.doRequest(http.MethodPost, uri, findRequest, findResponse)
	if err != nil {
		return nil, err
	}

	if len(findResponse.Response.Data) == 0 {
		return nil, fmt.Errorf("%v: %s", err, toErrorWithNewlines(uri, rawResp))
	}

	if findResponse.Status != "success" {
		return findResponse, errors.New(toErrorWithNewlines(uri, rawResp))
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

// https://www.hosting.de/api/?json#updating-records-in-a-zone
func (c *Client) updateRecords(updateRequest RecordsUpdateRequest) (*RecordsUpdateResponse, error) {
	uri := defaultBaseURL + "/recordsUpdate"

	updateResponse := &RecordsUpdateResponse{}

	rawResp, err := c.doRequest(http.MethodPost, uri, updateRequest, updateResponse)
	if err != nil {
		return nil, err
	}

	if updateResponse.Status != "success" && updateResponse.Status != "pending" {
		return nil, errors.New(toErrorWithNewlines(uri, rawResp))
	}

	return updateResponse, nil
}
