package uhppoted

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
)

type GetTimeProfileRequest struct {
	DeviceID  uint32
	ProfileID uint8
}

type GetTimeProfileResponse struct {
	DeviceID    DeviceID          `json:"device-id"`
	TimeProfile types.TimeProfile `json:"time-profile"`
}

func (u *UHPPOTED) GetTimeProfile(request GetTimeProfileRequest) (*GetTimeProfileResponse, error) {
	u.debug("get-time-profile", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	profileID := request.ProfileID

	profile, err := u.UHPPOTE.GetTimeProfile(deviceID, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving time profile %v from %v (%w)", profileID, deviceID, err))
	}

	if profile == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("Error retrieving time profile %v from %v", profileID, deviceID))
	}

	response := GetTimeProfileResponse{
		DeviceID:    DeviceID(deviceID),
		TimeProfile: *profile,
	}

	u.debug("get-time-profile", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type SetTimeProfileRequest struct {
	DeviceID    uint32
	TimeProfile types.TimeProfile
}

type SetTimeProfileResponse struct {
	DeviceID    DeviceID          `json:"device-id"`
	TimeProfile types.TimeProfile `json:"time-profile"`
}

func (u *UHPPOTED) SetTimeProfile(request SetTimeProfileRequest) (*SetTimeProfileResponse, error) {
	u.debug("set-time-profile", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	profile := request.TimeProfile
	linked := profile.LinkedProfileID

	if profile.ID < 2 || profile.ID > 254 {
		return nil, fmt.Errorf("Invalid time profile ID (%v) - valid range is [1..254]", profile.ID)
	}

	if linked != 0 {
		if linked == profile.ID {
			return nil, fmt.Errorf("Link to self creates circular reference")
		}

		if p, err := u.UHPPOTE.GetTimeProfile(deviceID, linked); err != nil {
			return nil, err
		} else if p == nil {
			return nil, fmt.Errorf("Linked time profile %v is not defined", linked)
		}

		profiles := map[uint8]bool{profile.ID: true}
		links := []uint8{profile.ID}
		for l := linked; l != 0; {
			if p, err := u.UHPPOTE.GetTimeProfile(deviceID, l); err != nil {
				return nil, err
			} else if p == nil {
				return nil, fmt.Errorf("Linked time profile %v is not defined", l)
			} else {
				links = append(links, p.ID)
				if profiles[p.ID] {
					return nil, fmt.Errorf("Linking to time profile %v creates a circular reference (%v)", linked, links)
				}

				profiles[p.ID] = true
				l = p.LinkedProfileID
			}
		}
	}

	ok, err := u.UHPPOTE.SetTimeProfile(deviceID, profile)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error writing time profile %v to %v (%w)", profile.ID, deviceID, err))
	}

	if !ok {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Failed to write time profile %v to %v", profile.ID, deviceID))
	}

	response := SetTimeProfileResponse{
		DeviceID:    DeviceID(deviceID),
		TimeProfile: profile,
	}

	u.debug("set-time-profile", fmt.Sprintf("response %+v", response))

	return &response, nil
}
