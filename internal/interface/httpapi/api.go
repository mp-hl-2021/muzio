package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/usecases/account"
	"github.com/mp-hl-2021/muzio/internal/usecases/entity"
	"github.com/mp-hl-2021/muzio/internal/usecases/playlist"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

const (
	entityIdUrlPathKey   = "entity_id"
	playlistIdUrlPathKey = "playlist_id"
	accountIdContextKey  = "account_id"
)

type Api struct {
	AccountUseCases       account.Interface
	MusicalEntityUseCases entity.Interface
	PlaylistUseCases      playlist.Interface
	Logger 				  zerolog.Logger
}

func NewApi(a account.Interface, e entity.Interface, p playlist.Interface) *Api {
	return &Api{
		AccountUseCases: a,
		MusicalEntityUseCases: e,
		PlaylistUseCases: p,
		Logger: log.With().Str("module", "http-server").Logger(),
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler())
	router.Use(measurer())

	router.HandleFunc("/blank", a.blank).Methods(http.MethodGet)
	router.HandleFunc("/blanka", a.authenticate(a.blanka)).Methods(http.MethodGet)

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/music/{"+entityIdUrlPathKey+"}", a.getMusicalEntity).Methods(http.MethodGet)

	router.HandleFunc("/playlist/{"+playlistIdUrlPathKey+"}", a.getPlaylist).Methods(http.MethodGet)
	router.HandleFunc("/playlist/{"+playlistIdUrlPathKey+"}", a.authenticate(a.putPlaylist)).Methods(http.MethodPut)
	router.HandleFunc("/playlist/{"+playlistIdUrlPathKey+"}", a.authenticate(a.deletePlaylist)).Methods(http.MethodDelete)

	router.HandleFunc("/drop/music", a.postMusicalEntity).Methods(http.MethodPost)
	router.HandleFunc("/drop/playlist", a.authenticate(a.postPlaylist)).Methods(http.MethodPost)

	router.Use(a.logger)

	return router
}

type PostSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Api) blank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "loc")
	psrm := PostSignupRequestModel{Login: "login1", Password: "password123"}
	b, err := json.Marshal(psrm)
	if err != nil{
		fmt.Println("JSONSHIT")
		return
	}
	fmt.Println(b)
	w.Write(b)
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) blanka(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "loc")
	fmt.Println("Authed")
	w.Write([]byte("Authed"))
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) postSignup(w http.ResponseWriter, r *http.Request) {
	var model PostSignupRequestModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		fmt.Println("Decode failed")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := a.AccountUseCases.CreateAccount(model.Login, model.Password)
	if err != nil {
		fmt.Println("Acc creation failed")
		fmt.Println(err)
		handleError(err, w)
		return
	}

	location := fmt.Sprintf("/accounts/%s", acc.Id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) postSignin(w http.ResponseWriter, r *http.Request) {
	var model PostSignupRequestModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.AccountUseCases.LoginToAccount(model.Login, model.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/jwt")
	w.Write([]byte(token))
}

type GetMusicalEntityResponseModel struct {
	Artist string        `json:"artist"`
	Album  string        `json:"album"`
	Track  string        `json:"track"`
	Links  []common.Link `json:"links"`
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
		handleError(err, w)
		return
	}
	m := GetMusicalEntityResponseModel{
		Artist: e.Artist,
		Album: e.Album,
		Track: e.Track,
		Links: e.Links,
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type getPlaylistResponseModel struct {
	Name    string                          `json:"name"`
	Content []GetMusicalEntityResponseModel `json:"content"`
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
		handleError(err, w)
		return
	}
	m := getPlaylistResponseModel{
		Name: p.Name,
		Content: make([]GetMusicalEntityResponseModel, 0, len(p.Content)),
	}
	for _, c := range p.Content {
		e, err := a.MusicalEntityUseCases.GetMusicalEntityById(c)
		if err != nil {
			handleError(err, w)
			return
		}
		em := GetMusicalEntityResponseModel{
			Artist: e.Artist,
			Album: e.Album,
			Track: e.Track,
			Links: e.Links,
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
	uid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := a.PlaylistUseCases.UpdatePlayList(uid, pid, m.Name, m.Content)
	if err != nil {
		handleError(err, w)
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
	uid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := a.PlaylistUseCases.DeletePlayList(uid, pid)
	if err != nil {
		handleError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type link struct {
	ServiceName string `json:"serviceName"`
	Url         string `json:"url"`
}

type postMusicalEntityRequestModel struct {
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Track  string `json:"track"`
	Links  []link `json:"links"`
}

type PostMusicalEntityResponseModel struct {
	Id string `json:"id"`
}

func (a *Api) postMusicalEntity(w http.ResponseWriter, r *http.Request) {
	var m postMusicalEntityRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	links := make([]common.Link, 0, len(m.Links))
	for _, l := range m.Links {
		links = append(links, common.Link{ServiceName: l.ServiceName, Url: l.Url, IsAvailable: true})
	}
	eid, err := a.MusicalEntityUseCases.CreateMusicalEntity(m.Artist, m.Album, m.Track, links)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nm := PostMusicalEntityResponseModel{Id: eid}
	if err := json.NewEncoder(w).Encode(nm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type postPlaylistRequestModel struct {
	Name    string   `json:"name"`
	Content []string `json:"content"`
}

func (a *Api) postPlaylist(w http.ResponseWriter, r *http.Request) {
	var m postPlaylistRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pid, err := a.PlaylistUseCases.CreatePlaylist(uid, m.Name, m.Content)
	if err != nil {
		handleError(err, w)
		return
	}
	nm := PostMusicalEntityResponseModel{Id: pid}
	if err := json.NewEncoder(w).Encode(nm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) accessibility(owner string, ctx context.Context) error {
	uid := ctx.Value(accountIdContextKey)
	if uid != owner {
		return domain.ErrForbidden
	}
	return nil
}

func (a *Api) authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		strArr := strings.Split(bearer, " ")
		if len(strArr) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := strArr[1]
		accId, err := a.AccountUseCases.Authenticate(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), accountIdContextKey, accId)
		handler(w, r.WithContext(ctx))
	}
}

func handleError(err error, w http.ResponseWriter) {
	if err == domain.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == domain.ErrUnauthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err == domain.ErrForbidden {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

type responseWriterObserver struct {
	http.ResponseWriter
	status 		int
	wroteHeader bool
}

func (o *responseWriterObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (o *responseWriterObserver) StatusCode() int {
	if !o.wroteHeader {
		return http.StatusOK
	}
	return o.status
}

func (a *Api) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		o := &responseWriterObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		fmt.Printf("method: %s; url: %s; status-code: %d; remote-addr: %s; duration: %v;\n",
			r.Method, r.URL.String(), o.StatusCode(), r.RemoteAddr, time.Since(start))
	})
}
