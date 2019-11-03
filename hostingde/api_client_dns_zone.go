package hostingde

// Adapted from hostingde provider in https://github.com/go-acme/lego

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v3"
	"time"
)

// https://www.hosting.de/api/?json#list-zoneconfigs
func (d *Client) listZoneConfigs(findRequest ZoneConfigsFindRequest) (*ZoneConfigsFindResponse, error) {
	uri := defaultBaseURL + "/zoneConfigsFind"

	findResponse := &ZoneConfigsFindResponse{}

	rawResp, err := d.post(uri, findRequest, findResponse)
	if err != nil {
		return nil, err
	}

	if len(findResponse.Response.Data) == 0 {
		return nil, fmt.Errorf("%v: %s", err, toUnreadableBodyMessage(uri, rawResp))
	}

	if findResponse.Status != "success" && findResponse.Status != "pending" {
		return findResponse, errors.New(toUnreadableBodyMessage(uri, rawResp))
	}

	return findResponse, nil
}

// https://www.hosting.de/api/?json#updating-zones
func (d *Client) updateZone(updateRequest ZoneUpdateRequest) (*ZoneUpdateResponse, error) {
	uri := defaultBaseURL + "/zoneUpdate"

	updateResponse := &ZoneUpdateResponse{}

	rawResp, err := d.post(uri, updateRequest, updateResponse)
	if err != nil {
		return nil, err
	}

	if updateResponse.Status != "success" && updateResponse.Status != "pending" {
		return nil, errors.New(toUnreadableBodyMessage(uri, rawResp))
	}

	return updateResponse, nil
}

// https://www.hosting.de/api/?json#creating-new-zones
func (d *Client) createZone(createRequest ZoneCreateRequest) (*ZoneCreateResponse, error) {
	uri := defaultBaseURL + "/zoneCreate"

	createResponse := &ZoneCreateResponse{}

	rawResp, err := d.post(uri, createRequest, createResponse)
	if err != nil {
		return nil, err
	}

	if createResponse.Status != "success" && createResponse.Status != "pending" {
		return nil, errors.New(toUnreadableBodyMessage(uri, rawResp))
	}

	return createResponse, nil
}

// https://www.hosting.de/api/?json#deleting-zones
func (d *Client) deleteZone(deleteRequest ZoneDeleteRequest) (*ZoneDeleteResponse, error) {
	uri := defaultBaseURL + "/zoneDelete"

	deleteResponse := &ZoneDeleteResponse{}

	rawResp, err := d.post(uri, deleteRequest, deleteResponse)
	if err != nil {
		return nil, err
	}

	if deleteResponse.Status != "success" && deleteResponse.Status != "pending" {
		return nil, errors.New(toUnreadableBodyMessage(uri, rawResp))
	}

	return deleteResponse, nil
}

func (d *Client) getZone(findRequest ZoneConfigsFindRequest) (*ZoneConfig, error) {
	var zoneConfig *ZoneConfig
	ctx, cancel := context.WithCancel(context.Background())

	operation := func() error {

		findResponse, err := d.listZoneConfigs(findRequest)
		if err != nil {
			cancel()
			return err
		}

		if findResponse.Response.Data[0].Status != "active" {
			return fmt.Errorf("unexpected status: %q", findResponse.Response.Data[0].Status)
		}

		zoneConfig = &findResponse.Response.Data[0]
		return nil
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 3 * time.Second
	bo.MaxInterval = 10 * bo.InitialInterval
	bo.MaxElapsedTime = 100 * bo.InitialInterval

	// retry in case the zone was edited recently and is not yet active
	err := backoff.Retry(operation, backoff.WithContext(bo, ctx))
	if err != nil {
		return nil, err
	}

	return zoneConfig, nil
}
