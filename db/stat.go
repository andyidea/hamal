package db

import (
	"hamal/models"
	"time"

	"github.com/jinzhu/gorm"
)

func AddOrUpdateDailyStat(stat *models.DailyStat) error {
	if err := db.Save(stat).Error; err != nil {
		return err
	}

	return nil
}

func GetOrAddDailyStat(userID uint, date time.Time) (*models.DailyStat, error) {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	var stat models.DailyStat
	stat.UserID = userID
	stat.Date = date
	if err := db.Where("date=DATE(?) and user_id = ?", date, userID).First(&stat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err2 := db.Create(&stat).Error; err2 != nil {
				return nil, err2
			}
		} else {
			return nil, err
		}
	}

	return &stat, nil
}

func GetDailyStats(userID uint) ([]models.DailyStat, error) {
	var dailyStats []models.DailyStat
	if err := db.Where("user_id = ?", userID).Find(&dailyStats).Error; err != nil {
		return nil, err
	}

	return dailyStats, nil
}
