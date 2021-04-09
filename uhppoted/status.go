package uhppoted

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
)

type Status struct {
	EventIndex     uint32          `json:"event-index"`
	EventType      byte            `json:"event-type"`
	Granted        bool            `json:"access-granted"`
	Door           byte            `json:"door"`
	Direction      uint8           `json:"direction"`
	CardNumber     uint32          `json:"card-number"`
	Timestamp      *types.DateTime `json:"event-timestamp,omitempty"`
	Reason         uint8           `json:"event-reason"`
	DoorState      map[uint8]bool  `json:"door-states"`
	DoorButton     map[uint8]bool  `json:"door-buttons"`
	SystemError    uint8           `json:"system-error"`
	SystemDateTime types.DateTime  `json:"system-datetime"`
	SequenceId     uint32          `json:"sequence-id"`
	SpecialInfo    uint8           `json:"special-info"`
	RelayState     uint8           `json:"relay-state"`
	InputState     uint8           `json:"input-state"`
}

type GetStatusRequest struct {
	DeviceID DeviceID
}

type GetStatusResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Status   Status   `json:"status"`
}

func (u *UHPPOTED) GetStatus(request GetStatusRequest) (*GetStatusResponse, error) {
	u.debug("get-status", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	status, err := u.Uhppote.GetStatus(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving status for %v (%w)", device, err))
	}

	response := GetStatusResponse{
		DeviceID: DeviceID(status.SerialNumber),
		Status: Status{
			EventIndex:     status.EventIndex,
			EventType:      status.EventType,
			Granted:        status.Granted,
			Door:           status.Door,
			Direction:      status.Direction,
			CardNumber:     status.CardNumber,
			Timestamp:      status.Timestamp,
			Reason:         status.Reason,
			DoorState:      status.DoorState,
			DoorButton:     status.DoorButton,
			SystemError:    status.SystemError,
			SystemDateTime: status.SystemDateTime,
			SequenceId:     status.SequenceId,
			SpecialInfo:    status.SpecialInfo,
			RelayState:     status.RelayState,
			InputState:     status.InputState,
		},
	}

	u.debug("get-status", fmt.Sprintf("response %+v", response))

	return &response, nil
}
