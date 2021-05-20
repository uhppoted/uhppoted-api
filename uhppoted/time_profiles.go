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
