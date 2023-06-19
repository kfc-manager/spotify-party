package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kfc-manager/spotify-party/models"
	"github.com/kfc-manager/spotify-party/spotify"
	"github.com/kfc-manager/spotify-party/utils"
)

type Server struct {
	Port string
    API *spotify.ApiConnection
}

func New(port string) *Server {
	return &Server{
		Port: port,
	}
}

// Initializes the server by creating an api connection to
// be able to interact with the spotify api. Also adds all 
// necessary routes to the router. After that it serves 
// those on the passed port.
func (s *Server) Run() error {

    router := mux.NewRouter()
    s.API = spotify.NewConnection()

    // routes 
    router.HandleFunc("/login", s.API.HandleLogin) // part of the spotify authorization flow

    router.HandleFunc("/callback", s.API.HandleRedirect) // part of the spotify authorization flow

    router.HandleFunc("/queue", s.HandleGetQueue)

    router.HandleFunc("/queue/{id}", s.HandleAddSong)

    router.HandleFunc("/search/{query}", s.HandleSearch)

    log.Println("Server is listening on port " + s.Port)

	err := http.ListenAndServe(s.Port, router)

    if err != nil {
        return err
    }

	return nil
}

// @return: JSON of the queue of songs, which are going 
// to be played by the users player
// @failure 405: The wrong http method has been used
// (not GET)
// @failure 502: Something went wrong while retrieving
// the data from the spotify api
func (s *Server) HandleGetQueue(w http.ResponseWriter, r *http.Request) {

    if r.Method != "GET" {
        utils.WriteJSON(w, 405, models.ErrorResponse{
            Error: models.ApiError{
                Status: 405,
                Message: "Wrong HTTP method",
            },
        })
        return 
    }

    queue, err := s.API.GetQueue() // get the queue from the spotify api

    if err != nil {
        utils.WriteJSON(w, 502, models.ErrorResponse{ 
            Error: models.ApiError{
                Status: 502,
                Message: err.Error(),
            },
        })
        return
    }

    res := &models.QueueResponse{ Queue: *queue, } // form a response
    err = utils.WriteJSON(w, 200, res)

    if err != nil {
        log.Println(err)
    }
}

// @param id (url-encoded): ID of the song which shall 
// be added to the players queue
// @return: No JSON if successful, only writes
// the header to status code 204 (command received)
// @failure 400: The ID string in the url was empty
// @failure 405: The wrong http method has been used
// (not POST)
// @failure 502: Something went wrong while forwarding
// the song that shall be added to the player's queue
// to the spoify api
func (s *Server) HandleAddSong(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        utils.WriteJSON(w, 405, models.ErrorResponse{
            Error: models.ApiError{
                Status: 405,
                Message: "Wrong HTTP method",
            },
        })
        return
    }

    id := mux.Vars(r)["id"] // get the id from the url of the request

    if len(id) < 1 {
        utils.WriteJSON(w, 400, models.ErrorResponse{
            Error: models.ApiError{
                Status: 400,
                Message: "ID is missing",
            },
        })
        return
    }

    err := s.API.AddSongToQueue(id) // add song to the queue via the spotify api

    if err != nil {
        utils.WriteJSON(w, 502, models.ErrorResponse{
            Error: models.ApiError{
                Status: 502,
                Message: err.Error(),
            },
        })
        return
    }

    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.WriteHeader(204) // status for command received
}

// @param query (url-encoded): Query of which matching 
// songs shall be returned
// @return: JSON of a list of songs that were matching
// the most with the passed query
// @failure 400: The query string in the url was empty
// @failure 405: The wrong http method has been used
// (not GET)
// @failure 502: Something went wrong while forwarding
// the query to the spotify api
func (s *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {

    if r.Method != "GET" {
        utils.WriteJSON(w, 405, models.ErrorResponse{
            Error: models.ApiError{
                Status: 405,
                Message: "Wrong HTTP method",
            },
        })
        return 
    }

    query := mux.Vars(r)["query"] // get the query from the url of the request

    if len(query) < 1 {
        utils.WriteJSON(w, 400, models.ErrorResponse{
            Error: models.ApiError{
                Status: 400,
                Message: "Empty query",
            },
        })
        return
    }

    result, err := s.API.SearchSong(query) // get list of songs from the spotify api

    if err != nil {
        utils.WriteJSON(w, 502, models.ErrorResponse{
            Error: models.ApiError{
                Status: 502,
                Message: err.Error(),
            },
        })
        return
    }

    res := &models.SearchResponse{ SearchResult: *result } // form a response
    err = utils.WriteJSON(w, 200, res)

    if err != nil {
        log.Println(err)
    }
}
