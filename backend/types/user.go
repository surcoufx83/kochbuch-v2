package types

import "time"

type UserProfile struct {
	Admin                 bool                   `json:"admin" db:"admin"`
	Created               time.Time              `json:"created" db:"created"`
	DisplayName           string                 `json:"displayname" db:"clouddisplayname"`
	Email                 NullString             `json:"email" db:"email"`
	EmailValidated        NullTime               `json:"-" db:"email_validated"`
	EmailValidationPhrase NullString             `json:"-" db:"email_validationphrase"`
	Enabled               bool                   `json:"enabled" db:"enabled"`
	FirstName             string                 `json:"firstname" db:"firstname"`
	Groups                []Group                `json:"groups"`
	Id                    uint32                 `json:"id" db:"user_id"`
	LastName              string                 `json:"lastname" db:"lastname"`
	Modified              time.Time              `json:"-" db:"modified"`
	NcEnabled             bool                   `json:"-" db:"cloudenabled"`
	NcSyncStatus          int16                  `json:"-" db:"cloudsync_status"`
	NcSyncTime            NullTime               `json:"-" db:"cloudsync"`
	UserName              string                 `json:"username" db:"cloudid"`
	Collections           map[uint32]*Collection `json:"collections" db:"-"`

	SimpleProfile NullUserProfileSimple `json:"-"`
}

type UserProfileSimple struct {
	Id          *uint32 `json:"id" db:"user_id"`
	DisplayName *string `json:"displayname" db:"clouddisplayname"`
}

type Group struct {
	Created     time.Time `json:"-" db:"created"`
	DisplayName string    `json:"displayname" db:"displayname"`
	GrantAccess bool      `json:"-" db:"access_granted"`
	GrantAdmin  string    `json:"-" db:"is_admin"`
	Id          uint16    `json:"id" db:"id"`
	Modified    NullTime  `json:"-" db:"modified"`
	Name        string    `json:"name" db:"ncname"`
}
