package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kfc-manager/spotify-party/models"
	"github.com/kfc-manager/spotify-party/utils"
)

const (
    AuthURL = "https://accounts.spotify.com/authorize"
    TokenURL = "https://accounts.spotify.com/api/token"
    ApiURL = "https://api.spotify.com/v1"
    Scopes = "user-read-currently-playing user-read-playback-state user-modify-playback-state"
)

type ApiConnection struct {
    ClientID string
    ClientSecret string
    RedirectURL string
    State string
    Client *http.Client
    Token string
}

func NewConnection() *ApiConnection {
    return &ApiConnection {
        ClientID: os.Getenv("SPOTIFY_ID"),
        ClientSecret: os.Getenv("SPOTIFY_SECRET"),
        RedirectURL: os.Getenv("REDIRECT_URL"),
        State: utils.GenerateRandomString(18),
        Client: &http.Client{},
    }
}

// This route is handling the authorization to the
// spotify api. It follows the authorization flow
// of the spotify web api. If you want to learn more:
// https://developer.spotify.com/documentation/web-api/tutorials/code-flow
func (auth *ApiConnection) HandleLogin(w http.ResponseWriter, r *http.Request) {
    
    req, err := http.NewRequest("GET", AuthURL, nil)

    if err != nil {
        utils.WriteJSON(w, 500, models.ErrorResponse{
            Error: models.ApiError{
                Status: 500,
                Message: err.Error(),
            },
        })
        return
    }

    query := req.URL.Query()
    query.Add("response_type", "code")
    query.Add("client_id", auth.ClientID)
    query.Add("scope", Scopes)
    query.Add("redirect_uri", auth.RedirectURL)
    query.Add("state", auth.State)
    req.URL.RawQuery = query.Encode()

    http.Redirect(w, r, req.URL.String(), http.StatusSeeOther)
}

// This route is the callback for the spotify api
// to send the authorization token with which we 
// can access the api. It is aswell part of the
// authorization flow of the spotify web api.
func (api *ApiConnection) HandleRedirect(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()
    state := r.FormValue("state")

    if api.State != state {
        utils.WriteJSON(w, 502, models.ErrorResponse{
            Error: models.ApiError{
                Status: 502,
                Message: "state missmatch",
            },
        })
        return
    }

    req, err := http.NewRequest("POST", TokenURL, nil) 

    if err != nil {
        utils.WriteJSON(w, 500, models.ErrorResponse{
            Error: models.ApiError{
                Status: 500,
                Message: err.Error(),
            },
        })
        return
    }

    params := req.URL.Query()
    params.Add("code", r.FormValue("code"))
    params.Add("redirect_uri", api.RedirectURL)
    params.Add("grant_type", "authorization_code")
    req.URL.RawQuery = params.Encode()

    req.Header.Add("Authorization", 
        "Basic " + utils.EncodeDetails(api.ClientID, api.ClientSecret))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    res, err := api.Client.Do(req)

    if err != nil {
        utils.WriteJSON(w, 500, models.ErrorResponse{
            Error: models.ApiError{
                Status: 500,
                Message: err.Error(),
            },
        })
        return 
    }

    defer res.Body.Close()

    bytes, err := io.ReadAll(res.Body)

    if err != nil {
        utils.WriteJSON(w, 500, models.ErrorResponse{
            Error: models.ApiError{
                Status: 500,
                Message: err.Error(),
            },
        })
        return 
    }

    data := &models.ApiResponse{}
    err = json.Unmarshal(bytes, data)

    if err != nil {
         utils.WriteJSON(w, 500, models.ErrorResponse{
            Error: models.ApiError{
                Status: 500,
                Message: err.Error(),
            },
        })
        return
    }

    if res.StatusCode != 200 {
        utils.WriteJSON(w, 502, models.ErrorResponse{
            Error: models.ApiError{
                Status: 502,
                Message: data.Error.Message,
            },
        })
    }

    // save token in our connection for the spotify api requests to access
    api.Token = data.Token 

    fmt.Fprintln(w, "Authorization was successful")
}

// Forms a request for the spotify api and retrieves
// the queue of songs which are going to be played 
// by the player of the users account. After
// that the data is transformed into the data type
// suitable for our api response. If you want to learn more:
// https://developer.spotify.com/documentation/web-api/reference/get-queue
func (api *ApiConnection) GetQueue() (*models.SongList, error) {

    req, err := http.NewRequest("GET", ApiURL + "/me/player/queue", nil)

    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", "Bearer " + api.Token)

    res, err := api.Client.Do(req)

    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    bytes, err := io.ReadAll(res.Body)

    if err != nil {
        return nil, err 
    }

    data := &models.ApiResponse{}
    err = json.Unmarshal(bytes, data)

    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        return nil, errors.New(data.Error.Message)
    }

    queue := utils.TransformToSongList(&data.CurrentlyPlaying, &data.Queue) 

    return queue, nil 
}

// Forms a request for the spotify api and adds the
// song with the passed id to the queue of songs
// played by the player of the users account.
// If you want to learn more:
// https://developer.spotify.com/documentation/web-api/reference/add-to-queue
func (api *ApiConnection) AddSongToQueue(id string) error {
    
    req, err := http.NewRequest("POST", ApiURL + "/me/player/queue", nil)

    if err != nil {
        return err
    }

    params := req.URL.Query()
    params.Add("uri", "spotify:track:" + id)
    req.URL.RawQuery = params.Encode()

    req.Header.Add("Authorization", "Bearer " + api.Token)

    res, err := api.Client.Do(req)

    if err != nil {
        return err
    }

    defer res.Body.Close()

    bytes, err := io.ReadAll(res.Body)

    if err != nil {
        return err
    }
    
    // the response body of the request is going
    // to be empty if the request was successful,
    // which means we are going to get a
    // "unexpected end of JSON input" error for
    // decoding the byte array
    data := &models.ApiResponse{}
    err = json.Unmarshal(bytes, data)

    if res.StatusCode != 204 {
        // err only interests us here because the
        // status code of the response is not positive
        // so the byte array should not be empty and 
        // if err != nil, something unexpected went
        // wrong while decoding it
        if err != nil { 
            return err
        }
        return errors.New(data.Error.Message)
    }

    return nil
}

// Forms a request for the spotify api and retrieves
// a list of songs matching the most with the passed
// query. After that the data is transformed into the data type
// suitable for our api response. If you want to learn more:
// https://developer.spotify.com/documentation/web-api/reference/search
func (api *ApiConnection) SearchSong(query string) (*models.SongList, error){

    req, err := http.NewRequest("GET", ApiURL + "/search", nil)

    if err != nil {
        return nil, err
    }

    params := req.URL.Query()
    params.Add("q", query)
    params.Add("type", "track")
    req.URL.RawQuery = params.Encode()

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", "Bearer " + api.Token)

    res, err := api.Client.Do(req)

    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    bytes, err := io.ReadAll(res.Body)

    if err != nil {
        return nil, err
    }

    data := &models.ApiResponse{}
    err = json.Unmarshal(bytes, data)

    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        return nil, errors.New(data.Error.Message)
    }

    result := utils.TransformToSongList(nil, &data.SearchResult.Items)

    return result, nil
}
