import { useEffect } from "react";
import useSearch from "../../hooks/useSearch";
import useSong from "../../hooks/useSong";
import { APISong } from "../../types";
import Backdrop from "../Backdrop/Backdrop";
import LoadingProgress from "../LoadingProgress/LoadingProgress";
import Song from "../Song/Song";
import { SearchProps } from "./types";
import { DOMAIN, ENDPOINTS } from "../../config";
import "./styles.css";

const Search = ({ query }: SearchProps): JSX.Element => {
  const { searchData, searchLoading, searchError } = useSearch(query);
  const { addSong, songLoading, songError } = useSong();

  useEffect((): void => {
    searchError &&
      searchError.type === "UNAUTHORIZED" &&
      window.location.assign(DOMAIN + ENDPOINTS.LOGIN);
    searchError && console.error(searchError.name);

    songError && console.error(songError.name);
    songError &&
      songError.type === "UNAUTHORIZED" &&
      window.location.assign(DOMAIN + ENDPOINTS.LOGIN);
  }, [searchError, songError]);

  return (
    <div className="search">
      {(searchLoading || songLoading) && (
        <Backdrop>
          <LoadingProgress />
        </Backdrop>
      )}
      {!searchError && (
        <div className="search-song-list">
          {searchData.map((song: APISong, index: number) => (
            <Song song={song} addSong={addSong} key={index} />
          ))}
        </div>
      )}
    </div>
  );
};

export default Search;
