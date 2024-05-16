package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/SEC-Jobstreet/backend-employer-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EmployerRepo struct {
	DB []gorm.DB
}

func NewEmployerRepo(dbConnections []string) *EmployerRepo {
	res := []gorm.DB{}
	for _, dbConnection := range dbConnections {

		sqlDB, err := sql.Open("pgx", dbConnection)
		if err != nil {
			log.Fatal().Msg("cannot connect to db")
		}
		store, err := gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Fatal().Msg("cannot connect to db")
		}

		err = models.MigrateEnterprises(store)
		if err != nil {
			log.Fatal().Msg("could not migrate db")
		}
		res = append(res, *store)
	}

	return &EmployerRepo{
		DB: res,
	}
}

func getDBShardKey(id uuid.UUID) uint32 {
	shardKey := id.ID() % 3
	fmt.Println(shardKey)
	return shardKey
}

func (repo *EmployerRepo) Create(enterprise *models.Enterprises) error {
	repo.DB[getDBShardKey(enterprise.ID)].Create(enterprise)
	return nil
}

func (repo *EmployerRepo) FindById(id uuid.UUID) (*models.Enterprises, error) {
	enterprise := &models.Enterprises{}

	err := repo.DB[getDBShardKey(id)].Where("id = ?", id).Find(enterprise).Error
	if err != nil {
		return nil, err
	}
	return enterprise, nil
}

func (repo *EmployerRepo) FindByEmployerId(id uuid.UUID) (*[]models.Enterprises, error) {
	enterprises0 := []models.Enterprises{}
	enterprises1 := []models.Enterprises{}
	enterprises2 := []models.Enterprises{}

	err := repo.DB[0].Where("employer_id = ?", id).Find(&enterprises0).Error
	if err != nil {
		return nil, err
	}
	err = repo.DB[1].Where("employer_id = ?", id).Find(&enterprises1).Error
	if err != nil {
		return nil, err
	}
	err = repo.DB[2].Where("employer_id = ?", id).Find(&enterprises2).Error
	if err != nil {
		return nil, err
	}
	enterprises0 = append(enterprises0, enterprises1...)
	enterprises0 = append(enterprises0, enterprises2...)
	return &enterprises0, nil
}
