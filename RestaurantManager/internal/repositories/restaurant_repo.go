package repositories

import (
	"restaurantservice/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) Create(restaurant *models.Restaurant) error {
	if restaurant.ID == uuid.Nil {
		restaurant.ID = uuid.New()
	}
	return r.db.Create(restaurant).Error
}

func (r *RestaurantRepository) GetByID(id uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := r.db.Where("id = ?", id).First(&restaurant).Error
	return &restaurant, err
}

func (r *RestaurantRepository) Update(restaurant *models.Restaurant) error {
	return r.db.Where("id = ?", restaurant.ID).Updates(restaurant).Error
}

func (r *RestaurantRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Restaurant{}).Error
}

func (r *RestaurantRepository) List() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Find(&restaurants).Error
	return restaurants, err
}
