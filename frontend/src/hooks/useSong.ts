import { APIError, APIErrorType } from "../types";
import { DOMAIN, ENDPOINTS } from "../config";
import { useState } from "react";

const useSong = (): {
  addSong: (id: string) => void;
  songLoading: boolean;
  songError: APIError | undefined;
} => {
  const [songLoading, setLoading] = useState<boolean>(false);
  const [songError, setError] = useState<APIError | undefined>(undefined);

  const addSong = async (id: string): Promise<void> => {
    setLoading(true);
    const params = new URLSearchParams();
    params.append("id", id);
    const response: Response = await fetch(
      DOMAIN + ENDPOINTS.ADD_SONG + "?" + params.toString(),
      {
        method: "POST",
      },
    );

    if (response.status === 401) {
      setError(new APIError(APIErrorType.UNAUTHORIZED, response.status));
    } else if (response.status === 429) {
      setError(new APIError(APIErrorType.TOO_MANY_REQUEST, response.status));
    } else if (!response.ok) {
      setError(new APIError(APIErrorType.INTERNAL_SERVER, response.status));
    }

    setLoading(false);
  };

  return { addSong, songLoading, songError };
};

export default useSong;
