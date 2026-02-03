package repository

import (
	"database/sql"
	"errors"
	
	"subscription-service/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}


func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	var endDate interface{}
	if sub.EndDate != nil {
		endDate = sub.EndDate.Time
	}

	return r.db.QueryRow(
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate.Time,
		endDate,
	).Scan(&sub.ID)
}


func (r *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var s models.Subscription
		var endDate sql.NullTime
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate.Time, &endDate); err != nil {
			return nil, err
		}
		if endDate.Valid {
			s.EndDate = &models.CustomDate{Time: endDate.Time}
		}
		subs = append(subs, s)
	}
	
	if subs == nil {
		subs = []models.Subscription{}
	}
	return subs, nil
}


func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	var s models.Subscription
	var endDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate.Time, &endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}
	if endDate.Valid {
		s.EndDate = &models.CustomDate{Time: endDate.Time}
	}
	return &s, nil
}


func (r *SubscriptionRepository) Update(sub *models.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET service_name = $1, price = $2, start_date = $3, end_date = $4, updated_at = NOW()
		WHERE id = $5`

	var endDate interface{}
	if sub.EndDate != nil {
		endDate = sub.EndDate.Time
	}

	res, err := r.db.Exec(query, sub.ServiceName, sub.Price, sub.StartDate.Time, endDate, sub.ID)
	if err != nil {
		return err
	}
	
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("subscription not found")
	}
	return nil
}


func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("subscription not found")
	}
	return nil
}


type CostFilter struct {
	UserID      uuid.UUID
	ServiceName string
	StartDate   models.CustomDate
	EndDate     models.CustomDate
}

func (r *SubscriptionRepository) GetTotalCost(filter CostFilter) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0) FROM subscriptions
		WHERE user_id = $1
		AND ($2 = '' OR service_name = $2)
		AND start_date <= $4
		AND (end_date IS NULL OR end_date >= $3)
	`
	var total int
	err := r.db.QueryRow(query, filter.UserID, filter.ServiceName, filter.StartDate.Time, filter.EndDate.Time).Scan(&total)
	return total, err
}
