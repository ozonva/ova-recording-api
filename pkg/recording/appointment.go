package recording

import (
	"fmt"
	"time"
)

type Appointment struct {
	UserID uint64
	AppointmentID uint64
	Name string
	Description string
	StartTime time.Time
	EndTime time.Time
}

func (receiver Appointment) String() string {
	return fmt.Sprintf("Appointment(user=%d,id=%d,name=%s)",
		receiver.UserID,
		receiver.AppointmentID,
		receiver.Name)
}
