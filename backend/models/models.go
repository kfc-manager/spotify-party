package models

// all data types received from the spotify api
type ApiResponse struct {
    CurrentlyPlaying ApiSong `json:"currently_playing"`
    Queue []ApiSong `json:"queue"`
    SearchResult ApiSearchResult `json:"tracks"`
    Token string `json:"access_token"`
    Error ApiError `json:"error"`
}

type ApiError struct {
    Status int `json:"status"`
    Message string `json:"message"`
}

type ApiSearchResult struct {
    Items []ApiSong `json:"items"`
}

type ApiSong struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Artists []ApiArtist `json:"artists"`
    Album ApiAlbum `json:"album"`
    Duration int `json:"duration_ms"`
}

type ApiAlbum struct {
    Images []ApiImage `json:"images"`
}

type ApiImage struct {
    URL string `json:"url"`
}

type ApiArtist struct {
    Name string `json:"name"`
}

// all data types my api sends to the user 

// I'm reusing ApiError from the api models as it
// has the attributes I need for this error type
type ErrorResponse struct {
    Error ApiError `json:"error"`
}

type QueueResponse struct {
    Queue SongList `json:"queue"`
}

type SearchResponse struct {
    SearchResult SongList `json:"result"`
}

// SongList can be the queue or a list of 
// songs as search result
type SongList struct {
    Songs []Song `json:"items"`
}

type Song struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Image string `json:"image_url"`
    Artists []string `json:"artists"`
    Duration int `json:"duration_ms"`
}
