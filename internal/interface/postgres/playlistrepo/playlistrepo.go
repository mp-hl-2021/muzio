package playlistrepo

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/domain/playlist"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const queryCreatePlaylist = `
	INSERT INTO playlists(
	    owner,
		name,
		content
	) VALUES ($1, $2, $3)
	RETURNING id
`

func (p *Postgres) CreatePlaylist(owner, name string, content []string) (playlist.Playlist, error) {
	pl := playlist.Playlist{
		Owner: owner,
		Name: name,
		Content: content,
	}
	row := p.conn.QueryRow(
		queryCreatePlaylist, owner, name, pq.Array(content),
	)
	err := row.Scan(&pl.Id)
	return pl, err
}

const queryGetPlayListById = `
	SELECT
		id,
		owner,
		name,
		content
	FROM playlists
	WHERE id = $1
`

func (p *Postgres) GetPlaylistById(id string) (playlist.Playlist, error) {
	pl := playlist.Playlist{}
	row := p.conn.QueryRow(queryGetPlayListById, id)
	err := row.Scan(&pl.Id, &pl.Owner, &pl.Name, pq.Array(&pl.Content))
	if err != nil && err == sql.ErrNoRows {
		return pl, domain.ErrNotFound
	}
	return pl, err
}

const queryUpdatePlaylist = `
	UPDATE playlists
	SET (name, content) = ($2, $3)
	WHERE id = $1
	RETURNING id
`

func (p *Postgres) UpdatePlaylist(id, name string, content []string) error {
	var uid string
	row := p.conn.QueryRow(queryUpdatePlaylist, id, name, content)
	err := row.Scan(&uid)
	if err != nil && err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

const queryDeletePlaylist = `
	DELETE
	FROM playlists
	WHERE id = $1
	RETURNING id
`

func (p *Postgres) DeletePlaylist(id string) error {
	var uid string
	row := p.conn.QueryRow(queryDeletePlaylist, id)
	err := row.Scan(&uid)
	if err != nil && err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}
