package uhppoted

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

const (
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
)

var (
	BadRequest          = errors.New("Bad Request")
	NotFound            = errors.New("Not Found")
	InternalServerError = errors.New("INTERNAL SERVER ERROR")
)

type IUhppote interface {
	DeviceList() map[uint32]*uhppote.Device
	GetDevices() ([]types.Device, error)
	GetDevice(deviceID uint32) (*types.Device, error)

	GetTime(serialNumber uint32) (*types.Time, error)
	SetTime(serialNumber uint32, datetime time.Time) (*types.Time, error)

	GetStatus(serialNumber uint32) (*types.Status, error)

	GetCards(deviceID uint32) (uint32, error)
	GetCardByIndex(deviceID, index uint32) (*types.Card, error)
	GetCardByID(deviceID, cardNumber uint32) (*types.Card, error)
	PutCard(deviceID uint32, card types.Card) (bool, error)
	DeleteCard(deviceID uint32, cardNumber uint32) (bool, error)
	DeleteCards(deviceID uint32) (bool, error)

	GetDoorControlState(deviceID uint32, door byte) (*types.DoorControlState, error)
	SetDoorControlState(deviceID uint32, door uint8, state uint8, delay uint8) (*types.DoorControlState, error)
	OpenDoor(deviceID uint32, door uint8) (*types.Result, error)

	GetEvent(deviceID, index uint32) (*types.Event, error)
	RecordSpecialEvents(deviceID uint32, enable bool) (bool, error)
	Listen(listener uhppote.Listener, q chan os.Signal) error
}

type UHPPOTED struct {
	Uhppote         IUhppote
	ListenBatchSize int
	Log             *log.Logger
}

func (u *UHPPOTED) debug(tag string, msg interface{}) {
	if u != nil && u.Log != nil {
		u.Log.Printf("DEBUG %-12s %v", tag, msg)
	}
}

func (u *UHPPOTED) info(tag string, msg interface{}) {
	if u != nil && u.Log != nil {
		u.Log.Printf("INFO  %-12s %v", tag, msg)
	}
}

func (u *UHPPOTED) warn(tag string, err error) {
	if u != nil && u.Log != nil {
		u.Log.Printf("WARN  %-12s %v", tag, err)
	}
}
