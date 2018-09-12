package db

import (
	"hamal/models"
	"hamal/protocol"

	"github.com/jinzhu/gorm"
)

func GetCurrencies(req *protocol.GetCurrenciesReq) (*protocol.GetCurrenciesRsp, error) {
	var currencies []models.Currency

	offset, limit := OffsetLimit(req.Page, req.Limit)

	query := db.Offset(offset).Limit(limit)
	if req.ID != "" {
		query = query.Where("id LIKE ?", "%"+req.ID+"%")
	}

	var count int
	if err := db.Find(&currencies).Count(&count).Error; err != nil {
		return nil, err
	}

	var resp protocol.GetCurrenciesRsp
	resp.Count = count
	resp.Items = make([]protocol.GetCurrenciesRspItem, len(currencies))
	for idx, currency := range currencies {
		resp.Items[idx].ID = currency.ID
		resp.Items[idx].Name = currency.Name
		resp.Items[idx].CreatedAt = currency.CreatedAt.Unix()
	}

	return &resp, nil

}

func AddOrUpdateCurrency(req *protocol.AddOrUpdateCurrencyReq) (*protocol.AddOrUpdateCurrencyRsp, error) {
	var currency models.Currency
	var isAdd bool
	if err := db.Where("id = ?", req.ID).First(&currency).Error; err == gorm.ErrRecordNotFound {
		isAdd = true
	} else if err != nil {
		return nil, err
	}

	if isAdd {
		currency.ID = req.ID
		currency.Name = req.Name
		if err := db.Create(&currency).Error; err != nil {
			return nil, err
		}
	} else {
		currency.Name = req.Name
		if err := db.Save(&currency).Error; err != nil {
			return nil, err
		}
	}
	var rsp protocol.AddOrUpdateCurrencyRsp
	return &rsp, nil
}

func DeleteCurrency(id string) error {
	if err := db.Where("id = ?", id).Delete(&models.Currency{}).Error; err != nil {
		return err
	}

	return nil
}
