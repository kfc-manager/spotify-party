import { useEffect, useState } from "react";
import { APIError, APIErrorType, APISong } from "../types";
import { DOMAIN, ENDPOINTS } from "../config";

type QueueResponse = {
  queue: APISong[];
};

const useQueue = (): {
  queueData: APISong[];
  queueLoading: boolean;
  queueError: APIError | undefined;
} => {
  const [queueData, setData] = useState<APISong[]>([]);
  const [queueLoading, setLoading] = useState<boolean>(true);
  const [queueError, setError] = useState<APIError | undefined>(undefined);

  const fetchQueue = async (): Promise<QueueResponse> => {
    const response: Response = await fetch(DOMAIN + ENDPOINTS.QUEUE);

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

  useEffect((): void => {
    fetchQueue()
      .then((response: QueueResponse): void => setData(response.queue))
      .catch((error: APIError): void => setError(error))
      .finally((): void => setLoading(false));
  }, []);

  return { queueData, queueLoading, queueError };
};

export default useQueue;
