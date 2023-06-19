export interface ApiResponse {
    error: ApiError;
    queue: ApiItems;
    result: ApiItems;
}

export interface ApiError {
    status: number;
    message: string;
}

export interface ApiItems {
    items: ApiSong[];
}

export interface ApiSong {
    name: string;
    id: string;
    image_url: string;
    artists: string[];
    duration_ms: number;
}
