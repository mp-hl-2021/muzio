package entityrepo

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/domain/entity"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const queryCreateMusicalEntity = `
	INSERT INTO entities(
	    artist,
		album,
		track
		links
	) VALUES ($1, $2, $3, $4::link[])
	RETURNING id
`

func (p *Postgres) CreateMusicalEntity(artist, album, track string, links []common.Link) (entity.MusicalEntity, error) {
	e := entity.MusicalEntity{
		Artist: artist,
		Album: album,
		Track: track,
		Links: links,
	}
	row := p.conn.QueryRow(
		queryCreateMusicalEntity, artist, album, track, pq.Array(links),
	)
	err := row.Scan(&e.Id)
	return e, err
}

const queryGetMusicalEntityById = `
	SELECT
		id,
		artist,
		album,
		track,
		links::link[]
	FROM entities
	WHERE id = $1
`

func (p *Postgres) GetMusicalEntityById(id string) (entity.MusicalEntity, error) {
	e := entity.MusicalEntity{}
	row := p.conn.QueryRow(queryGetMusicalEntityById, id)
	err := row.Scan(&e.Id, &e.Artist, &e.Album, &e.Track, pq.Array(&e.Links))
	if err != nil && err == sql.ErrNoRows {
		return e, domain.ErrNotFound
	}
	return e, err
}

const queryGetBatchToCheck = `
	UPDATE entities
	SET checkedAt = now()
	WHERE id IN (
		SELECT id FROM entities 
		ORDER BY checkedAt ASC  
		LIMIT $1
	)
	RETURNING id, artist, album, track, links::link[]
`

func (p *Postgres) GetBatchToCheck(number int) ([]entity.MusicalEntity, error) {
	result := make([]entity.MusicalEntity, 0, number)
	rows, err := p.conn.Query(queryGetBatchToCheck, number)
	if err != nil {
		return result, err
	}
	for {
		if !rows.Next() {
			return result, rows.Err()
		}
		e := entity.MusicalEntity{}
		err := rows.Scan(&e.Id, &e.Artist, &e.Album, &e.Track, pq.Array(&e.Links))
		if err != nil {
			return result, domain.ErrNotFound
		}
		result = append(result, e)
	}
}

const queryUpdateLinks = `
	UPDATE entities
	SET links = $2::link[]
	WHERE id = $1
	RETURNING id
`

func (p *Postgres) UpdateLinks(id string, links []common.Link) error {
	var uid string
	row := p.conn.QueryRow(queryUpdateLinks, id, pq.Array(links))
	err := row.Scan(&uid)
	if err != nil && err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}
