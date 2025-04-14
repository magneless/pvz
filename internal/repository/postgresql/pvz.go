package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/domain/repository"
)

type PVZRepository struct {
	db *sql.DB
}

func NewPVZRepository(db *sql.DB) repository.PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

func (r *PVZRepository) Create(ctx context.Context, pvz models.PVZ) error {
	const query = `
		INSERT INTO pvz (id, registration_date, city)
		VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(
		ctx,
		query,
		pvz.ID,
		*pvz.RegistrationDate,
		pvz.City,
	)
	return err
}

type rowData struct {
	PVZID            uuid.UUID   `db:"pvz_id"`
	RegistrationDate *time.Time  `db:"pvz_registrationDate"`
	City             models.City `db:"pvz_city"`

	ReceptionID       *uuid.UUID     `db:"reception_id"`
	ReceptionDateTime *time.Time     `db:"reception_dateTime"`
	ReceptionStatus   *models.Status `db:"reception_status"`

	ProductID       *uuid.UUID   `db:"product_id"`
	ProductDateTime *time.Time   `db:"product_dateTime"`
	ProductType     *models.Type `db:"product_type"`
}

func (r *PVZRepository) GetAllInfo(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.PVZFullInfo, error) {
	const query = `
		SELECT
			pvz.id AS pvz_id,
			pvz.registration_date AS pvz_registrationDate,
			pvz.city AS pvz_city,
			r.id AS reception_id,
			r.date_time AS reception_dateTime,
			r.status AS reception_status,
			pr.id AS product_id,
			pr.date_time AS product_dateTime,
			pr.type AS product_type
		FROM PVZ pvz
		LEFT JOIN receptions r 
			ON pvz.id = r.pvz_id 
			AND r.date_time BETWEEN $1 AND $2
		LEFT JOIN products pr 
			ON r.id = pr.reception_id
		ORDER BY pvz.registration_date
		LIMIT $4
		OFFSET (($3 - 1) * $4);`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		startDate,
		endDate,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pvzMap := make(map[uuid.UUID]*models.PVZFullInfo)

	for rows.Next() {
		var row rowData
		err := rows.Scan(
			&row.PVZID,
			&row.RegistrationDate,
			&row.City,
			&row.ReceptionID,
			&row.ReceptionDateTime,
			&row.ReceptionStatus,
			&row.ProductID,
			&row.ProductDateTime,
			&row.ProductType,
		)
		if err != nil {
			return nil, err
		}

		pvzInfo, exists := pvzMap[row.PVZID]
		if !exists {
			pvzInfo = &models.PVZFullInfo{
				PVZ: models.PVZ{
					ID:               row.PVZID,
					RegistrationDate: row.RegistrationDate,
					City:             row.City,
				},
				Receptions: []models.ReceptionFullInfo{},
			}
			pvzMap[row.PVZID] = pvzInfo
		}

		if row.ReceptionID != nil {
			var receptionInfo *models.ReceptionFullInfo
			for i := range pvzInfo.Receptions {
				if pvzInfo.Receptions[i].ID == *row.ReceptionID {
					receptionInfo = &pvzInfo.Receptions[i]
					break
				}
			}
			if receptionInfo == nil {
				reception := models.ReceptionFullInfo{
					Reception: models.Reception{
						ID:       *row.ReceptionID,
						DateTime: row.ReceptionDateTime,
						PVZID:    row.PVZID,
						Status:   *row.ReceptionStatus,
					},
					Products: []models.Product{},
				}
				pvzInfo.Receptions = append(pvzInfo.Receptions, reception)

				receptionInfo = &pvzInfo.Receptions[len(pvzInfo.Receptions)-1]
			}

			if row.ProductID != nil {
				product := models.Product{
					ID:          *row.ProductID,
					DateTime:    row.ProductDateTime,
					ReceptionID: *row.ReceptionID,
					Type:        *row.ProductType,
				}
				receptionInfo.Products = append(receptionInfo.Products, product)
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var result []models.PVZFullInfo
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	return result, nil
}
