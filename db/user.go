package db

import (
	"hamal/models"
	"hamal/protocol"
	"hamal/util"

	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

func UserRegister() {

}

func UserLogin(req *protocol.LoginReq) (*protocol.LoginRsp, error) {
	var user models.User
	if err := db.First(&user, "username = ?", req.AccountName).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			return nil, protocol.NewError(protocol.CodeUsernameNoExist, protocol.MSgUsernameNoExist, err.Error())
		}

	}

	if user.ID == 0 {
		return nil, protocol.NewError(protocol.CodeUsernameNoExist, protocol.MSgUsernameNoExist, "")
	}

	if user.PasswordMd5 != util.Md5(req.Password) {
		return nil, protocol.NewError(protocol.CodePasswordError, protocol.MsgPasswordError, "")
	}

	var tokenUser util.TokenUser
	tokenUser.ID = user.ID

	m := structs.Map(&tokenUser)
	m["token_created_at"] = time.Now().Unix()

	token, err := util.GenerateToken(m)
	if err != nil {
		return nil, err
	}

	var rsp protocol.LoginRsp
	rsp.Token = token
	return &rsp, nil
}

func GetUserInfoModel(userID uint) (*models.User, error) {
	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserInfo(reqUser *models.User) (*protocol.GetUserInfoRsp, error) {
	var user models.User
	if err := db.First(&user, "id = ?", reqUser.ID).Error; err != nil {
		return nil, err
	}

	var rsp protocol.GetUserInfoRsp
	rsp.Name = user.Name
	rsp.Avatar = "http://7xqrd2.com1.z0.glb.clouddn.com/FiWvpofL8-QjDfjrgLLGAUFZCj_n.jpeg"
	if user.Avatar != "" {
		rsp.Avatar = user.Avatar
	}
	rsp.Roles = []string{"admin"}

	return &rsp, nil

}

func UpdateUserInfo(user *models.User) error {
	return db.Save(user).Error
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserExchages(userID uint, exchange string) ([]models.UserExchange, error) {
	var ues []models.UserExchange

	query := db.Where("user_id = ?", userID)

	if exchange != "" {
		query = query.Where("exchange = ?", exchange)
	}
	err := query.Find(&ues).Error
	if err != nil {
		return nil, err
	}

	return ues, nil
}

func GetUserExchange(id uint) (*models.UserExchange, error) {
	var ue models.UserExchange
	err := db.Find(&ue).Where("id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &ue, nil
}

func AddUserExchange(ue *models.UserExchange) error {
	return db.Create(ue).Error
}

func DeleteUserExchange(id uint, userID uint) error {
	if err := db.Where("id = ? and user_id = ?", id, userID).Delete(&models.UserExchange{}).Error; err != nil {
		return err
	}

	return nil
}
