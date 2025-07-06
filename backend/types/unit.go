package types

import "time"

type Unit struct {
	CreatedTime    time.Time                   `json:"created"`
	DecimalPlaces  uint8                       `json:"decimalPlaces"`
	Id             uint8                       `json:"id"`
	Localization   map[string]UnitLocalization `json:"localization"`
	ModifiedTime   time.Time                   `json:"modified"`
	ReplacedById   uint8                       `json:"replacedBy"`
	SavedAsFactor  float32                     `json:"savedAsFactor"`
	SavedAsId      uint8                       `json:"savedAs"`
	ShowAsFraction bool                        `json:"showAsFraction"`
}

type UnitLocalization struct {
	NameSingular string `json:"singular"`
	NamePlural   string `json:"plural"`
}
