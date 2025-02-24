package repositories

import (
    "context"
    "orderservice/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
)

type OrderHistoryRepository struct {
    collection *mongo.Collection
}

func NewOrderHistoryRepository(db *mongo.Database) *OrderHistoryRepository {
    return &OrderHistoryRepository{
        collection: db.Collection("order_history"),
    }
}

func (r *OrderHistoryRepository) CreateOrderHistory(ctx context.Context, history *models.OrderHistory) error {
    _, err := r.collection.InsertOne(ctx, history)
    return err
}

func (r *OrderHistoryRepository) UpdateOrderHistory(ctx context.Context, orderID string, log models.OrderLog) error {
    filter := map[string]interface{}{"order_id": orderID}
    update := map[string]interface{}{
        "$push": map[string]interface{}{
            "log": log,
        },
        "$set": map[string]interface{}{
            "status": log.Status,
        },
    }
    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

func (r *OrderHistoryRepository) GetOrderHistory(ctx context.Context, orderID string) (*models.OrderHistory, error) {
    var history models.OrderHistory
    filter := map[string]interface{}{"order_id": orderID}
    err := r.collection.FindOne(ctx, filter).Decode(&history)
    if err != nil {
        return nil, err
    }
    return &history, nil
}