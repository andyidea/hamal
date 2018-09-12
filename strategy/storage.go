package strategy

import (
	"hamal/common"
	"hamal/db"
	"log"
	"sync"
)

var instanceMap sync.Map

func AddStrategyInstance(instance StrategyIF) {
	instanceMap.Store(instance.ID(), instance)
}

func GetStrategyInstance(id uint) (StrategyIF, bool) {
	data, ok := instanceMap.Load(id)
	if ok {
		return data.(StrategyIF), true
	}
	return nil, false
}

func LaunchStrategyInstance() {
	sis, err := db.GetStrategyInstancesModel()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("xx")

	for _, si := range sis {
		if si.Status == common.StrategyInstanceStatusRunning {
			log.Println("xxx")
			var inst StrategyIF
			for {
				inst, err = NewStrategyInstance(si.ID)
				if err != nil {
					log.Println(err.Error())

					continue
				}

				break
			}

			log.Println("bb")

			inst.Launch()

			AddStrategyInstance(inst)
			log.Println("aa")
		}
	}

}
