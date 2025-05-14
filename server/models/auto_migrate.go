package model

import "godocms/pkg/dbfactory"

func AutoMigrate() error {
	dbfactory.Db.AutoMigrate(&User{})
	dbfactory.Db.AutoMigrate(&UserRole{})
	dbfactory.Db.AutoMigrate(&UserDept{})
	dbfactory.Db.AutoMigrate(&UserThird{})
	return nil
}
