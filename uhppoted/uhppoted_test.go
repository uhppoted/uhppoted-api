package uhppoted

import (
	"fmt"
	"os"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type stub struct {
	recordSpecialEvents func(deviceID uint32, enable bool) (bool, error)
}

func (m *stub) DeviceList() map[uint32]*uhppote.Device {
	return nil
}

func (m *stub) FindDevices() ([]types.Device, error) {
	return nil, nil
}

func (m *stub) FindDevice(deviceID uint32) (*types.Device, error) {
	return nil, nil
}

func (m *stub) GetTime(serialNumber uint32) (*types.Time, error) {
	return nil, nil
}

func (m *stub) SetTime(serialNumber uint32, datetime time.Time) (*types.Time, error) {
	return nil, nil
}

func (m *stub) GetStatus(serialNumber uint32) (*types.Status, error) {
	return nil, nil
}

func (m *stub) GetCards(deviceID uint32) (uint32, error) {
	return 0, nil
}

func (m *stub) GetCardByIndex(deviceID, index uint32) (*types.Card, error) {
	return nil, nil
}

func (m *stub) GetCardByID(deviceID, cardNumber uint32) (*types.Card, error) {
	return nil, nil
}

func (m *stub) PutCard(deviceID uint32, card types.Card) (bool, error) {
	return false, nil
}

func (m *stub) DeleteCard(deviceID uint32, cardNumber uint32) (bool, error) {
	return false, nil
}

func (m *stub) DeleteCards(deviceID uint32) (bool, error) {
	return false, nil
}

func (m *stub) GetDoorControlState(deviceID uint32, door byte) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *stub) SetDoorControlState(deviceID uint32, door uint8, state uint8, delay uint8) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *stub) OpenDoor(deviceID uint32, door uint8) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetEvent(deviceID, index uint32) (*types.Event, error) {
	return nil, nil
}

func (m *stub) RecordSpecialEvents(deviceID uint32, enable bool) (bool, error) {
	if m.recordSpecialEvents != nil {
		return m.recordSpecialEvents(deviceID, enable)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) Listen(listener uhppote.Listener, q chan os.Signal) error {
	return nil
}