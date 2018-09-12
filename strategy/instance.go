package strategy

import (
	"errors"
	"hamal/db"
)

type LogChan chan string

type StrategyIF interface {
	ID() uint
	Launch()
	Stop()
}

func NewStrategyInstance(id uint) (StrategyIF, error) {
	si, err := db.GetStrategyInstance(id)
	if err != nil {
		return nil, err
	}

	switch si.StrategyID {
	case "coinex_duiqiao":
		return NewStrategyCoinEXDuiQiao(id)
	case "high_frequency_lang":
		return NewStrategyHighFrequencyLang(id)
	}

	return nil, errors.New("no stragegy instance.")
}

//type StrategyInstance struct {
//	model    *models.StrategyInstance
//	trader   *fcoin.FCoin
//	strategy StrategyIF
//	logch    LogChan
//}
//
//func NewStrategyInstance(id uint, user *models.User) (*StrategyInstance, error) {
//	inst := StrategyInstance{}
//	var err error
//
//	inst.logch = make(LogChan, 200)
//
//	inst.model, err = db.GetStrategyInstanceModel(id)
//	if err != nil {
//		return nil, err
//	}
//
//	inst.trader, err = fcoin.NewFCoin("", "")
//	if err != nil {
//		return nil, err
//	}
//	if inst.model.StrategyID == common.StrategyIDHighFrequencyLang {
//		//inst.strategy = &StrategyHighFrequencyLang{}
//		//inst.strategy.Init(inst.model)
//	}
//
//	return &inst, nil
//}
//
//func (si *StrategyInstance) Launch() error {
//	var err error
//	si.model, err = db.GetStrategyInstanceModel(si.model.ID)
//	if err != nil {
//		return err
//	}
//
//	if err := db.UpdateStrategyInstanceStatus(si.model.ID, common.StrategyInstanceStatusRunning); err != nil {
//		return err
//	}
//	go si.strategy.Launch(si.logch, si.trader, si.model.BaseCurrencyID, si.model.TargetCurrencyID, si.model.Amount, si.model.Interval)
//
//	return nil
//}
//
//func (si *StrategyInstance) Stop() error {
//	if err := db.UpdateStrategyInstanceStatus(si.model.ID, common.StrategyInstanceStatusRest); err != nil {
//		return err
//	}
//	si.strategy.Stop()
//
//	return nil
//}
//
//func (si *StrategyInstance) ID() uint {
//	return si.model.ID
//}
