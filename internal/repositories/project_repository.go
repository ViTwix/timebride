package repositories

import (
    "database/sql"
    "timebride/internal/models"
)

type ProjectRepository struct {
    DB *sql.DB
}

// GetAll повертає список усіх проєктів з БД
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
    rows, err := r.DB.Query("SELECT id, client_name, instagram, location, date, price, deposit, notes FROM projects")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var projects []models.Project
    for rows.Next() {
        var p models.Project
        if err := rows.Scan(&p.ID, &p.ClientName, &p.Instagram, &p.Location, &p.Date, &p.Price, &p.Deposit, &p.Notes); err != nil {
            return nil, err
        }
        projects = append(projects, p)
    }

    return projects, nil
}
