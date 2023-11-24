import { useEffect } from "react";
import useQueue from "../../hooks/useQueue";
import { DOMAIN, ENDPOINTS } from "../../config";
import Backdrop from "../Backdrop/Backdrop";
import LoadingProgress from "../LoadingProgress/LoadingProgress";
import { APISong } from "../../types";
import Song from "../Song/Song";
import "./styles.css";

const Queue = (): JSX.Element => {
  const { queueData, queueLoading, queueError } = useQueue();

  useEffect((): void => {
    queueError &&
      queueError.type === "UNAUTHORIZED" &&
      window.location.assign(DOMAIN + ENDPOINTS.LOGIN);
    queueError && console.error(queueError.name);
  }, [queueError]);

  return (
    <div className="queue">
      {queueError && (
        <Backdrop>
          <div>{queueError.type}</div>
        </Backdrop>
      )}
      {queueLoading && (
        <Backdrop>
          <LoadingProgress />
        </Backdrop>
      )}
      {!(queueError || queueLoading) && (
        <div className="queue-song-list">
          {queueData.map((song: APISong, index: number) => (
            <Song song={song} key={index} />
          ))}
        </div>
      )}
    </div>
  );
};

export default Queue;
