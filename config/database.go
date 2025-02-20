package config

import(
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB()(*gorm.DB, error){
	db, err := gorm.Open(postgres.Open("postgresql://neondb_owner:npg_qXyMNZF4IK3n@ep-damp-heart-a8zcur1e-pooler.eastus2.azure.neon.tech/neondb?sslmode=require"), &gorm.Config{})
	
	if err != nil{
		return nil, err
	}
	return db, nil
}