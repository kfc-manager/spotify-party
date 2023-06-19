package utils

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/kfc-manager/spotify-party/models"
)

const chars string = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz1234567890"

// Method to generate a random string of a certian length.
// This is used to generate the state, which is crucial
// for the authorization flow of the spotify api.
// If you want to learn more:
// https://developer.spotify.com/documentation/web-api/tutorials/code-flow
func GenerateRandomString(length int) string {
    builder := ""
    
    for i := 0 ; i < length ; i++ {
        index := rand.Intn(len(chars))
        builder += string(chars[index])
    }

    return builder
}

// Method to base64-encoding the CLIENT_ID and CLIENT_SECRET.
// This is required by the authorization flow of the spotify api. 
func EncodeDetails(id string, secret string) string {
    original := id + ":" + secret
    return base64.StdEncoding.EncodeToString([]byte(original)) 
}

// Method to write JSON of the passed data to the passed
// response writer. Also sets the status of the response
// and allows Access of Origin, which is important for
// testing of the application in a localhost environment.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
    w.Header().Add("Content-Type", "application/json")
    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

// Method to transform a whole list of songs of the api
// response to a list of songs suitable for our response
// to the user. Transform each song with the method
// "TransformToSong" induvidually. Also accepts a single
// song "currentlyPlaying" if we transform a queue. If
// currentlyPlaying == nil, we transform a search result.
func TransformToSongList(currentlyPlaying *models.ApiSong, songs *[]models.ApiSong) *models.SongList {
    list := &models.SongList{}
    // currentlyPlaying.ID can be an empty string if the 
    // user is currently not listening to music but the
    // queue is requested from the spotify api
    if currentlyPlaying != nil && len(currentlyPlaying.ID) > 0 {
        list.Songs = append(list.Songs, *TransformToSong(currentlyPlaying))
    }
    for _, elem := range *songs {
        list.Songs= append(list.Songs, *TransformToSong(&elem))
    }
    return list 
}

// Method to transform a song of the api response to
// a song suitable for our response to the user.
// Mostly stripping the api response of paramaters,
// which we are not interested in.
func TransformToSong(apiSong *models.ApiSong) *models.Song {
    song := &models.Song{}
    song.Name = apiSong.Name
    song.ID = apiSong.ID
    if len(apiSong.Album.Images) > 0 {
        song.Image = apiSong.Album.Images[0].URL
    }
    song.Duration = apiSong.Duration
    for _, elem := range apiSong.Artists {
        song.Artists = append(song.Artists, elem.Name)
    }
    return song
}
