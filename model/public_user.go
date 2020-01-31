package model

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	genid "github.com/srlemon/gen-id"
	userBase "github.com/srlemon/userDetail"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type UserStatus int

const (
	UserStatusNormal    UserStatus = iota // 正常
	UserStatusNotActive                   // 未激活
	UserStatusLoginOut                    // 登出
	UserStatusFreeze                      // 冻结
)

const (
	randStr = `1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM`
)

func (u *UserDetail) TableName() string {
	return "user_detail"
}

// UserDetail
type UserDetail struct {
	Uid              string         `json:"uid" gorm:"primary_key;size:36;index"`
	IsAdmin          bool           `json:"isAdmin" gorm:"default:false;index"`
	Nickname         string         `json:"nickname" gorm:"size:24;index"`
	Username         string         `json:"username" gorm:"size:16;unique_index"` // 最长16位,只支持数字字母大小写
	Status           UserStatus     `json:"status" gorm:"index;default(1)"`       // 用户状态,0正常在线,1账号未激活,2账号冻结
	IsChangeUsername bool           `json:"isChangeUsername"gorm:"default:false"`
	Verified         bool           `json:"verified" gorm:"default:false"` // true:实名认证
	RealName         string         `json:"realName" gorm:"size:36;index"`
	LocNum           string         `json:"locNum" gorm:"unique_index:loc_num_phone;size:6"` // 	所在地区(国家),国际电话区号,不带"+"号
	Phone            string         `json:"phone"gorm:"size:11;unique_index:loc_num_phone"`
	HeadIcon         string         `json:"headIcon"gorm:"size:256"`
	Email            string         `json:"email"gorm:"size:64;index"`
	LoginPassword    string         `json:"loginPassword" gorm:"size:128;column:loginPassword"`
	PayPassword      string         `json:"payPassword"gorm:"size:128"`
	Sign             string         `json:"sign" gorm:"size:64"`                     // 签名
	Role             pq.StringArray `json:"role" gorm:"not null;type:varchar(36)[]"` // 角色, 用户可以拥有多个角色
	Level            int32          `json:"level" gorm:"default:1"`                  // 等级
	Secret           string         `json:"secret" gorm:"not null;type:varchar(32)"` // 用户自己的密钥
	// Religion         string         `json:"religion" gorm:"index;size:24"`           // 宗教
	TimeData
	// 外键关联
	BankCards []*BankCard      `json:"bankCards" gorm:"foreignkey:Uid;"` // association_autoupdate:false"
	Addr      []*AddressDetail `json:"addr" gorm:"foreignkey:Uid;"`      // association_autoupdate:false"
	IDCard    IDCard           `json:"idCard" gorm:"foreignkey:uid"`
}

// checkUserEmailExist 验证邮箱是否已经存在
func checkUserEmailExist(email string) bool {
	var (
		err error
	)
	if _, err = PubUserGetByEmail(email); err != nil {
		return false
	}

	return true
}

// NewUser
func NewUser(u *UserDetail) (ret *UserDetail) {
	if u != nil {
		ret = u
	} else {
		ret = new(UserDetail)
	}
	//
	if len(ret.Uid) != 36 {
		ret.Uid = uuid.NewV4().String()
	}
	//
	if len(ret.Secret) == 0 {
		ret.Secret = ""
	}
	ret.Role = []string{
		"a",
	}
	if len(ret.LocNum) == 0 {
		ret.LocNum = "86"
	}
	if len(ret.Nickname) == 0 {
		g := genid.NewGeneratorData(nil)
		ret.Nickname = g.Name
	}

	if len(ret.Username) == 0 {
		i := 0
		for i < 16 {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			ret.Username += string(randStr[r.Intn(len(randStr))])
			i++
		}

	}
	return
}

// PubUserGet
func PubUserGet(uid string) (ret *UserDetail, err error) {
	ret = new(UserDetail)

	if err = Database.Model(ret).Where("uid = ?", uid).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = userBase.ErrUserNotExist
		}
		return
	}

	return
}

// PubUserGetByEmail
func PubUserGetByEmail(email string) (ret *UserDetail, err error) {
	if len(email) == 0 {
		err = userBase.ErrFormParamInValid
		return
	}
	data := new(UserDetail)
	if err = Database.Model(data).Where("email = ?", email).Find(data).Error; err != nil {
		return
	}

	ret = data
	return
}

// PubUserGetByUsername
func PubUserGetByUsername(username string) (ret *UserDetail, err error) {
	ret = new(UserDetail)

	if err = Database.Model(ret).Where("username = ?", username).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = userBase.ErrUserNotExist
		}
		return
	}
	return
}

// PubUserGetByPhone
func PubUserGetByPhone(locNum, phone string) (ret *UserDetail, err error) {
	if len(locNum) == 0 {
		locNum = "86"
	}
	ret = new(UserDetail)

	if err = Database.Model(ret).Where("phone = ? and loc_num = ?", phone, locNum).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = userBase.ErrUserNotExist
		}
		return
	}
	return
}

// PubUserAdd 添加用户
func PubUserAdd(f *userBase.FormRegister) (ret *UserDetail, err error) {
	if err = f.Valid(); err != nil {
		return
	}
	var (
		data = NewUser(nil)
		pwd  []byte
	)

	data.Phone = f.Phone
	if pwd, err = bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost); err != nil {
		return
	}
	data.LoginPassword = string(pwd)
	// TODO 随机生成数据

	if err = Database.Model(new(UserDetail)).Create(data).Error; err != nil {
		return
	}

	ret = data
	return
}

// PubUserDel 删除用户,软删除
func PubUserDel(uid string) (err error) {
	var (
		data         = new(UserDetail)
		dataAddr     []*AddressDetail
		dataBankCard []*BankCard
		dataID       = new(IDCard)
	)
	data.Uid = uid
	Database.Model(data).Association("Addr").Find(&dataAddr)
	Database.Model(data).Association("BankCards").Find(&dataBankCard)
	Database.Model(data).Association("IDCard").Find(dataID)

	// 存在关联数据，直接删除
	if len(dataAddr) > 0 {
		if err = Database.Model(data).Association("Addr").Delete(dataAddr).Error; err != nil {
			return
		}
	}
	if len(dataID.IDCard) > 0 {
		if err = Database.Model(data).Association("IDCard").Delete(dataID).Error; err != nil {
			return
		}
	}

	//
	if len(dataBankCard) > 0 {
		if err = Database.Model(data).Association("BankCards").Delete(dataBankCard).Error; err != nil {
			return
		}
	}

	if err = Database.Delete(data).Error; err != nil {
		return
	}

	return
}

// PubUserUpdate
func PubUserUpdate(uid string, f *userBase.UpdateUserProfile) (ret *UserDetail, err error) {

	var (
		data *UserDetail
		ok   bool
		m    map[string]interface{}
	)
	if m, err = f.Valid(); err != nil {
		return
	}
	if data, err = PubUserGet(uid); err != nil {
		return
	}
	if _, ok = m["username"]; ok && data.IsChangeUsername {
		err = userBase.ErrActionNotAllow.SetDetail("action not allow, username is changed")
		return
	}
	if err = Database.Model(&UserDetail{}).Where("uid = ?", uid).Updates(m).Error; err != nil {
		return
	}

	return PubUserGet(uid)
}