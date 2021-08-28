package recording

import (
	"fmt"
	"time"
)

type Appointment struct {
	AppointmentID uint64 	`db:"appointment_id"`
	UserID uint64 			`db:"user_id"`
	Name string 			`db:"name"`
	Description string		`db:"description"`
	StartTime time.Time		`db:"start_time"`
	EndTime time.Time		`db:"end_time"`
}

func (receiver Appointment) String() string {
	return fmt.Sprintf("Appointment(user=%d,id=%d,name=%s)",
		receiver.UserID,
		receiver.AppointmentID,
		receiver.Name)
}
