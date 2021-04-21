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
