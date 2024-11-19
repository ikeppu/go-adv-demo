package stat

import (
	"arch/ikeppu/github.com/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	repo.DB.Find(&stat, "link_id = ? and date = ?", linkId, datatypes.Date(time.Now()))

	if stat.ID == 0 {
		repo.DB.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   datatypes.Date(time.Now()),
		})
	} else {
		stat.Clicks += 1
		repo.DB.Save(&stat)
	}
}
