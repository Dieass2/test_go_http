package models

import (
	"encoding/json"
	// "fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)


type CustomDate struct {
	time.Time
}

const customLayout = "01-2006"

func (c *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(customLayout, s)
	if err != nil {
		return err
	}
	c.Time = t
	return nil
}

func (c CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Time.Format(customLayout))
}

type Subscription struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	ServiceName string      `json:"service_name" binding:"required" db:"service_name"`
	Price       int         `json:"price" binding:"required,min=0" db:"price"`
	UserID      uuid.UUID   `json:"user_id" binding:"required" db:"user_id"`
	StartDate   CustomDate  `json:"start_date" binding:"required" db:"start_date"`
	EndDate     *CustomDate `json:"end_date,omitempty" db:"end_date"`
}
