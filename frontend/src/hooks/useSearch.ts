import { useEffect, useState } from "react";
import { APIError, APIErrorType, APISong } from "../types";
import { DOMAIN, ENDPOINTS } from "../config";

type SearchResponse = {
  tracks: APISong[];
};

const useSearch = (
  query: string,
): {
  searchData: APISong[];
  searchLoading: boolean;
  searchError: APIError | undefined;
} => {
  const [searchData, setData] = useState<APISong[]>([]);
  const [searchLoading, setLoading] = useState<boolean>(true);
  const [searchError, setError] = useState<APIError | undefined>(undefined);
  const [abortController, setAbortController] = useState<AbortController>(
    new AbortController(),
  );

  const fetchSearch = async (
    query: string,
    abortController: AbortController,
  ): Promise<SearchResponse> => {
    const params = new URLSearchParams();
    params.append("query", query);
    const response: Response = await fetch(
      DOMAIN + ENDPOINTS.SEARCH_SONG + "?" + params.toString(),
      { signal: abortController.signal },
    );

    if (response.status === 401) {
      throw new APIError(APIErrorType.UNAUTHORIZED, response.status);
    }
    if (response.status === 429) {
      throw new APIError(APIErrorType.TOO_MANY_REQUEST, response.status);
    }
    if (!response.ok) {
      throw new APIError(APIErrorType.INTERNAL_SERVER, response.status);
    }

    return response.json();
  };

  useEffect(() => {
    if (query.length < 1) return;
    abortController.abort();
    const newAbortController: AbortController = new AbortController();
    setAbortController(newAbortController);
    setLoading(true);
    fetchSearch(query, newAbortController)
      .then((response: SearchResponse) => setData(response.tracks))
      .catch((error: Error) => {
        error instanceof APIError && setError(error);
      })
      .finally(() => setLoading(false));
  }, [query]);

  return { searchData, searchLoading, searchError };
};

export default useSearch;
