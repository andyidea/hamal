package db

import (
	"hamal/common"
	"hamal/models"
	"hamal/protocol"
	"time"
)

func GetStrategies(req *protocol.GetStrategiesReq) (*protocol.GetStrategiesRsp, error) {
	var strategies []models.Strategy
	var count int
	if err := db.Find(&strategies).Count(&count).Error; err != nil {
		return nil, err
	}

	var rsp protocol.GetStrategiesRsp
	rsp.Count = count
	for _, strategy := range strategies {
		var item protocol.GetStrategiesRspItem
		item.ID = strategy.ID
		item.Name = strategy.Name
		item.CreatedAt = strategy.CreatedAt.Unix()

		rsp.Items = append(rsp.Items, item)
	}

	return &rsp, nil
}

func GetStrategyInstancesModel() ([]models.StrategyInstance, error) {
	var strategyInstances []models.StrategyInstance
	if err := db.Find(&strategyInstances).Error; err != nil {
		return nil, err
	}

	return strategyInstances, nil
}

func GetStrategyInstances(user *models.User, strategyID string) (*protocol.GetStrategyInstancesRsp, error) {
	var strategyInstances []models.StrategyInstance
	var count int
	if err := db.Where("user_id = ? and strategy_id = ?", user.ID, strategyID).Find(&strategyInstances).Count(&count).Error; err != nil {
		return nil, err
	}

	var rsp protocol.GetStrategyInstancesRsp
	rsp.Count = count
	for _, strategyInstance := range strategyInstances {
		var item protocol.GetStrategyInstancesRspItem
		item.ID = strategyInstance.ID
		item.StrategyID = strategyInstance.StrategyID
		item.BaseCurrencyID = strategyInstance.BaseCurrencyID
		item.TargetCurrencyID = strategyInstance.TargetCurrencyID
		item.Amount = strategyInstance.Amount
		item.Interval = strategyInstance.Interval
		item.Status = strategyInstance.Status
		item.CreatedAt = strategyInstance.CreatedAt.Unix()
		item.UserExchangeID = strategyInstance.UserExchangeID

		rsp.Items = append(rsp.Items, item)
	}

	return &rsp, nil
}

func GetStrategyInstanceModel(id uint) (*models.StrategyInstance, error) {
	var strategyInstance models.StrategyInstance
	if err := db.First(&strategyInstance, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &strategyInstance, nil
}

func GetStrategyInstance(id uint) (*protocol.GetStrategyInstanceRsp, error) {
	var strategyInstance models.StrategyInstance
	if err := db.First(&strategyInstance, "id = ?", id).Error; err != nil {
		return nil, err
	}

	var rsp protocol.GetStrategyInstanceRsp
	rsp.ID = strategyInstance.ID
	rsp.StrategyID = strategyInstance.StrategyID
	rsp.BaseCurrencyID = strategyInstance.BaseCurrencyID
	rsp.TargetCurrencyID = strategyInstance.TargetCurrencyID
	rsp.Amount = strategyInstance.Amount
	rsp.Status = strategyInstance.Status
	rsp.CreatedAt = strategyInstance.CreatedAt.Unix()

	return &rsp, nil
}

func AddStrategyInstance(req *protocol.AddStrategyInstanceReq, user *models.User) (*protocol.AddStrategyInstanceRsp, error) {
	var strategyInstance models.StrategyInstance

	strategyInstance.StrategyID = req.StrategyID
	strategyInstance.BaseCurrencyID = req.BaseCurrencyID
	strategyInstance.TargetCurrencyID = req.TargetCurrencyID
	strategyInstance.Amount = req.Amount
	strategyInstance.Status = common.StrategyInstanceStatusRest
	strategyInstance.UserID = user.ID
	strategyInstance.Interval = req.Interval
	strategyInstance.UserExchangeID = req.UserExchangeID

	if err := db.Create(&strategyInstance).Error; err != nil {
		return nil, err
	}

	var rsp protocol.AddStrategyInstanceRsp

	rsp.ID = strategyInstance.ID
	rsp.Status = common.StrategyInstanceStatusRest
	rsp.Amount = strategyInstance.Amount
	rsp.BaseCurrencyID = strategyInstance.BaseCurrencyID
	rsp.TargetCurrencyID = strategyInstance.TargetCurrencyID
	rsp.CreatedAt = strategyInstance.CreatedAt.Unix()

	return &rsp, nil
}

func UpdateStrategyInstance(si *models.StrategyInstance) error {
	return db.Save(si).Error
}

func UpdateStrategyInstanceStatus(id uint, status string) error {
	si, err := GetStrategyInstanceModel(id)
	if err != nil {
		return err
	}

	si.Status = status

	if err := db.Save(si).Error; err != nil {
		return err
	}

	return nil
}

func DeleteStrategyInstance(id uint) error {
	if err := db.Where("id = ?", id).Delete(&models.StrategyInstance{}).Error; err != nil {
		return err
	}

	return nil
}

func AddStrategyInstanceRecard(sir *models.StrategyInstanceRecard) error {
	return db.Create(sir).Error
}

func GetStrategyInstanceRecard(id uint) (*models.StrategyInstanceRecard, error) {
	var sir models.StrategyInstanceRecard
	if err := db.First(&sir, "id = ? ", id).Error; err != nil {
		return nil, err
	}

	return &sir, nil
}

type GetStrategyInstanceRecardsParams struct {
	StrategyInstanceID uint
	Status             string
	UserID             uint
	UpdateTimeBegin    *time.Time
	UpdateTimeEnd      *time.Time
	CreateTimeBegin    *time.Time
	CreateTimeEnd      *time.Time
}

func GetStrategyInstanceRecards(params *GetStrategyInstanceRecardsParams) ([]models.StrategyInstanceRecard, error) {
	var sirs []models.StrategyInstanceRecard
	query := db.New()
	if params.StrategyInstanceID > 0 {
		query = query.Where("strategy_instance_id = ?", params.StrategyInstanceID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.UpdateTimeBegin != nil {
		query = query.Where("updated_at >= ?", params.UpdateTimeBegin)
	}
	if params.UpdateTimeEnd != nil {
		query = query.Where("updated_at <= ?", params.UpdateTimeEnd)
	}
	if params.CreateTimeBegin != nil {
		query = query.Where("created_at >= ?", params.CreateTimeBegin)
	}
	if params.UpdateTimeEnd != nil {
		query = query.Where("created_at <= ?", params.CreateTimeEnd)
	}

	if err := query.Find(&sirs).Error; err != nil {
		return nil, err
	}

	return sirs, nil
}

func GetStrategyInstanceRecardByOrderID(orderID string) (*models.StrategyInstanceRecard, error) {
	var sir models.StrategyInstanceRecard
	if err := db.First(&sir, "order_id = ? ", orderID).Error; err != nil {
		return nil, err
	}

	return &sir, nil
}

func UpdateStrategyInstanceRecard(sir *models.StrategyInstanceRecard) error {
	return db.Save(sir).Error
}
