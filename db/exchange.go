package db

import "hamal/models"

func GetExchanges() ([]models.Exchange, error) {
	var exs []models.Exchange
	err := db.Find(&exs).Error
	if err != nil {
		return nil, err
	}

	return exs, nil
}
