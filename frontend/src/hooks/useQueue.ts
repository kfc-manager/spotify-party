import { useEffect, useState } from "react"
import { APIError, APIErrorType, Song } from "../types"
import { DOMAIN, ENDPOINTS } from "../config"

const useQueue = (): { data: Song[], loading: boolean, error: APIError | undefined } => {
  const [data, setData] = useState<Song[]>([])
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<APIError | undefined>(undefined)

  const fetchData = async (): Promise<Song[]> => {
    const response = await fetch(DOMAIN + ENDPOINTS.QUEUE, { method: "GET", headers: { "Content-Type": "application/json" } })

    if (response.status === 401) { throw new APIError(APIErrorType.UNAUTHORIZED, response.status) }
    if (response.status === 429) { throw new APIError(APIErrorType.TOO_MANY_REQUEST, response.status) }
    if (!response.ok) { throw new APIError(APIErrorType.INTERNAL_SERVER, response.status) }

    return response.json()
  }

  useEffect((): void => {
    fetchData()
      .then((response: Song[]): void => setData(response))
      .catch((error: APIError): void => setError(error))
      .finally((): void => setLoading(false))
  }, [])

  return { data, loading, error }
}

export default useQueue;
