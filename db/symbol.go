package db

import (
	"hamal/market"
	"hamal/models"
	"hamal/protocol"

	"github.com/jinzhu/gorm"
)

func GetSymbols() (*protocol.GetSymbolsRsp, error) {
	var symbols []models.Symbol
	var count int
	if err := db.Find(&symbols).Count(&count).Error; err != nil {
		return nil, err
	}

	var resp protocol.GetSymbolsRsp
	resp.Count = count

	for _, symbol := range symbols {
		var item protocol.GetSymbolsRspItem
		item.ID = symbol.ID
		item.Name = symbol.Name
		item.BaseCurrencyID = symbol.BaseCurrencyID
		item.QuoteCurrencyID = symbol.QuoteCurrencyID
		item.CreatedAt = symbol.CreatedAt.Unix()

		current, ok := market.GetCurrentDetail(symbol.ID)
		if ok {
			if len(current.Tick.Data) > 0 {
				item.LastetPrice = current.Tick.Data[0].Price
			}

		}

		resp.Items = append(resp.Items, item)
	}

	return &resp, nil
}

func GetSymbolsModel() ([]models.Symbol, error) {
	var symbols []models.Symbol
	if err := db.Find(&symbols).Error; err != nil {
		return nil, err
	}

	return symbols, nil
}

func AddOrUpdateSymbol(s *models.Symbol) error {
	var symbol models.Symbol
	var isAdd bool
	if err := db.Where("id = ?", s.ID).First(&symbol).Error; err == gorm.ErrRecordNotFound {
		isAdd = true
	} else if err != nil {
		return err
	}

	if isAdd {
		if err := db.Create(s).Error; err != nil {
			return err
		}
		//market.HuobiAPI.SubscribeDepth(s.ID)
		//market.HuobiAPI.SubscribeDetail(s.ID)
	} else {
		symbol.Name = s.Name
		symbol.QuoteCurrencyID = s.QuoteCurrencyID
		symbol.BaseCurrencyID = s.BaseCurrencyID
		if err := db.Save(&symbol).Error; err != nil {
			return err
		}
	}
	return nil
}
