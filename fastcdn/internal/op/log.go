package op

import (
	// "fmt"
	// "time"

	// "github.com/pkg/errors"
	// "gorm.io/gorm"

	"fastcdn/internal/conf"
	"fastcdn/internal/db"
	// "fastcdn/internal/model"
	// utils "fastcdn/internal/utils"
)

func AddLog(uid int64, content string) error {
	if conf.Security.InstallLock {
		return db.AddLog(nil, uid, content)
	}
	return nil

}

func SysLog(content string) error {
	return db.AddLog(nil, 0, content)
}
