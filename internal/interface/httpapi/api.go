package httpapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/muzio/internal/usecases/account"
	"github.com/mp-hl-2021/muzio/internal/usecases/entity"
	"github.com/mp-hl-2021/muzio/internal/usecases/playlist"
	"net/http"
)

const (
	entityIdUrlPathKey   = "entity_id"
	playlistIdUrlPathKey = "playlist_id"
)

type Api struct {
	AccountUseCases       account.Interface
	MusicalEntityUseCases entity.Interface
	PlaylistUseCases      playlist.Interface
}

func NewApi(a account.Interface, e entity.Interface, p playlist.Interface) *Api {
	return &Api{
		AccountUseCases: a,
		MusicalEntityUseCases: e,
		PlaylistUseCases: p,
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/music/{"+entityIdUrlPathKey+"}", a.getMusicalEntity).Methods(http.MethodGet)

	router.HandleFunc("/playlist/{"+playlistIdUrlPathKey+"}", a.getPlaylist).Methods(http.MethodGet)
	router.HandleFunc("/playlist/{"+playlistIdUrlPathKey+"}", a.authenticate(a.putPlaylist)).Methods(http.MethodPut)
	router.HandleFunc("/playlist/{"+playlistIdUrlPathKey+"}", a.authenticate(a.deletePlaylist)).Methods(http.MethodDelete)

	router.HandleFunc("/drop/music/", a.postMusicalEntity).Methods(http.MethodPost)
	router.HandleFunc("/drop/playlist/", a.authenticate(a.postPlaylist)).Methods(http.MethodPost)

	return router
}

type postSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Api) postSignup(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) postSignin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type link struct {
	ServiceName string `json:"serviceName"`
	Url         string `json:"url"`
}

type getMusicalEntityResponseModel struct {
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Track  string `json:"track"`
	Links  []link `json:"links"`
}

func (a *Api) getMusicalEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eid, ok := vars[entityIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e, err := a.MusicalEntityUseCases.GetMusicalEntityById(eid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := getMusicalEntityResponseModel{
		Artist: e.Artist,
		Album: e.Album,
		Track: e.Track,
		Links: make([]link, 0, len(e.Links)),
	}
	for _, l := range e.Links {
		m.Links = append(m.Links, link{
			ServiceName: l.ServiceName,
			Url: l.Url,
		})
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type getPlaylistResponseModel struct {
	Name    string                          `json:"name"`
	Content []getMusicalEntityResponseModel `json:"content"`
}

func (a *Api) getPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[playlistIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := a.PlaylistUseCases.GetPlaylistById(pid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := getPlaylistResponseModel{
		Name: p.Name,
		Content: make([]getMusicalEntityResponseModel, 0, len(p.Content)),
	}
	for _, c := range p.Content {
		e, err := a.MusicalEntityUseCases.GetMusicalEntityById(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		em := getMusicalEntityResponseModel{
			Artist: e.Artist,
			Album: e.Album,
			Track: e.Track,
			Links: make([]link, 0, len(e.Links)),
		}
		for _, l := range em.Links {
			em.Links = append(em.Links, link{
				ServiceName: l.ServiceName,
				Url: l.Url,
			})
		}
		m.Content = append(m.Content, em)
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type putPlaylistRequestModel struct {
	Name    string   `json:"name"`
	Content []string `json:"content"`
}

func (a *Api) putPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[playlistIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var m putPlaylistRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.PlaylistUseCases.UpdatePlayList(pid, m.Name, m.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deletePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[playlistIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.PlaylistUseCases.DeletePlayList(pid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type postEntityResponseModel struct {
	Id string `json:"name"`
}

func (a *Api) postMusicalEntity(w http.ResponseWriter, r *http.Request) {
	var m getMusicalEntityResponseModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nl := make([]entity.Link, 0, len(m.Links))
	for _, l := range nl {
		nl = append(nl, entity.Link{
			ServiceName: l.ServiceName,
			Url: l.Url,
		})
	}
	eid, err := a.MusicalEntityUseCases.CreateMusicalEntity(m.Artist, m.Album, m.Track, nl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nm := postEntityResponseModel{Id: eid}
	if err := json.NewEncoder(w).Encode(nm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type postPlaylistRequestModel struct {
	Owner   string   `json:"owner"` // TODO: Auth
	Name    string   `json:"name"`
	Content []string `json:"content"`
}

func (a *Api) postPlaylist(w http.ResponseWriter, r *http.Request) {
	var m postPlaylistRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pid, err := a.PlaylistUseCases.CreatePlaylist(m.Owner, m.Name, m.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nm := postEntityResponseModel{Id: pid}
	if err := json.NewEncoder(w).Encode(nm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) authenticate(handler http.HandlerFunc) http.HandlerFunc {
	// TODO implement
	return handler
}
