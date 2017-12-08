package gamedata

import (
	"adventure/advserver/log"
	"Adventure/common/clog"
	"Adventure/common/csv"
)

var logger *clog.Logger

var AllTemplates Templates

type Templates struct {
	HeroTemplate struct {
		HeroName csv.String `table:"hero" key:"" val:"名字"`
		SkillID  csv.Int    `table:"hero" key:"" val:"技能ID列表"`
	}
}

func Init() {
	logger = log.GetLogger()

	AllTemplates = Templates{}

	csv.LoadTemplates(&AllTemplates)
}
