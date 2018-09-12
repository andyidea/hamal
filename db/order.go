package db

import "hamal/models"

func GetOrders(userEchangeID uint) ([]models.OrderHistory, error) {
	var orders []models.OrderHistory
	err := db.Find(&orders).Where("user_exchange_id = ?", userEchangeID).Order("id desc").Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}
