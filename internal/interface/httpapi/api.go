package httpapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	entityIdUrlPathKey   = "{entity_id}"
	playlistIdUrlPathKey = "{playlist_id}"
)

type Api struct {
	// AccountUseCases       account.Interface
	// MusicalEntityUseCases entity.Interface
	// PlaylistUseCases      playlist.Interface
}

/* TODO
func NewApi(a account.Interface, e entity.Interface, p playlist.Interface) *Api {
	return &Api{
		AccountUseCases: a,
		MusicalEntityUseCases: e,
		PlaylistUseCases: p,
	}
}
 */

func NewApi() *Api {
	return &Api{}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/music/"+entityIdUrlPathKey, a.getMusicalEntity).Methods(http.MethodGet)
	router.HandleFunc("/music/"+entityIdUrlPathKey, a.putMusicalEntity).Methods(http.MethodPut)

	router.HandleFunc("/playlist/"+playlistIdUrlPathKey, a.getPlaylist).Methods(http.MethodGet)
	router.HandleFunc("/playlist/"+playlistIdUrlPathKey, a.putPlaylist).Methods(http.MethodPut)
	router.HandleFunc("/playlist/"+playlistIdUrlPathKey, a.deletePlaylist).Methods(http.MethodDelete)

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
	m := getMusicalEntityResponseModel{
		Artist: "Radiohead", Album: "Hail To the Thief", Track: "2 + 2 = 5",
		Links: []link{link{ServiceName: "Yandex Music", Url: "https://music.yandex.ru/track/333416"}},
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) putMusicalEntity(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type getPlaylistResponseModel struct {
	Content  []getMusicalEntityResponseModel `json:"content"`
}

func (a *Api) getPlaylist(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) putPlaylist(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) deletePlaylist(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) authenticate(handler http.HandlerFunc) http.HandlerFunc {
	panic("implement me")
}
