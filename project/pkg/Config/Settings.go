package Config

import (
	"github.com/antoniocs/go_api/pkg/Util"
)

type Settings struct {
	DatabaseDriver   string
	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	AppPort          string
	VenueUsed        string
}

func NewSettings() Settings {
	return Util.PopulateWithEnv(Settings{})
}
