type Song = {
  id: string;
  name: string;
  image_url: string;
  artists: string[];
  duration_ms: number;
}

export type { Song }

enum APIErrorType {
  UNAUTHORIZED = "UNAUTHORIZED",
  BAD_REQUEST = "BAD_REQUEST",
  INTERNAL_SERVER = "INTERNAL_SERVER",
  TOO_MANY_REQUEST = "TOO_MANY_REQUEST",
}

class APIError extends Error {
  type: APIErrorType;
  status: number;
  constructor(type: APIErrorType, status: number, ...args: any) {
    super(...args);
    this.type = type;
    this.status = status;
  }

  get name() {
    return `APIError[${this.type}]`
  }
}

export { APIError, APIErrorType };
