package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
)

type DoorControl struct {
	Device struct {
		ID           uint32 `json:"id"`
		Door         uint8  `json:"door"`
		ControlState string `json:"control"`
	} `json:"device"`
}

type GetDoorDelayRequest struct {
	DeviceID uint32
	Door     uint8
}

type GetDoorDelayResponse struct {
	DeviceID uint32 `json:"device-id"`
	Door     uint8  `json:"door"`
	Delay    uint8  `json:"delay"`
}

func (u *UHPPOTED) GetDoorDelay(ctx context.Context, request GetDoorDelayRequest) (*GetDoorDelayResponse, int, error) {
	u.debug(ctx, "get-door-delay", fmt.Sprintf("request  %+v", request))

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := GetDoorDelayResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug(ctx, "get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

type SetDoorDelayRequest struct {
	DeviceID uint32
	Door     uint8
	Delay    uint8
}

type SetDoorDelayResponse struct {
	DeviceID uint32 `json:"device-id"`
	Door     uint8  `json:"door"`
	Delay    uint8  `json:"delay"`
}

func (u *UHPPOTED) SetDoorDelay(ctx context.Context, request SetDoorDelayRequest) (*SetDoorDelayResponse, int, error) {
	u.debug(ctx, "set-door-delay", fmt.Sprintf("request  %+v", request))

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(request.DeviceID, request.Door, state.ControlState, request.Delay)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := SetDoorDelayResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug(ctx, "get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

type GetDoorControlRequest struct {
	DeviceID uint32
	Door     uint8
}

type GetDoorControlResponse struct {
	DeviceID uint32 `json:"device-id"`
	Door     uint8  `json:"door"`
	Control  string `json:"control"`
}

func (u *UHPPOTED) GetDoorControl(ctx context.Context, request GetDoorControlRequest) (*GetDoorControlResponse, int, error) {
	u.debug(ctx, "get-door-control", fmt.Sprintf("request  %+v", request))

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	lookup := map[uint8]string{
		1: "normally open",
		2: "normally closed",
		3: "controlled",
	}

	response := GetDoorControlResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Control:  lookup[result.ControlState],
	}

	u.debug(ctx, "get-door-control", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

func (u *UHPPOTED) SetDoorControl(ctx context.Context, rq Request) {
	u.debug(ctx, "set-door-control", rq)

	id, door, control, err := rq.DeviceDoorControl()
	if err != nil {
		u.warn(ctx, 0, "set-door-control", err)
		u.oops(ctx, "set-door-control", "Invalid device/door/state)", StatusBadRequest)
		return
	}

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(*id, *door)
	if err != nil {
		u.warn(ctx, *id, "set-door-control", err)
		u.oops(ctx, "set-door-control", "Error setting door control", StatusInternalServerError)
		return
	}

	states := map[string]uint8{
		"normally open":   1,
		"normally closed": 2,
		"controlled":      3,
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(*id, *door, states[*control], state.Delay)
	if err != nil {
		u.warn(ctx, *id, "set-door-control", err)
		u.oops(ctx, "set-door-control", "Error setting door control", StatusInternalServerError)
		return
	}

	lookup := map[uint8]string{
		1: "normally open",
		2: "normally closed",
		3: "controlled",
	}

	response := DoorControl{
		struct {
			ID           uint32 `json:"id"`
			Door         uint8  `json:"door"`
			ControlState string `json:"control"`
		}{
			ID:           *id,
			Door:         result.Door,
			ControlState: lookup[result.ControlState],
		},
	}

	u.reply(ctx, response)
}
