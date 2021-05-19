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
