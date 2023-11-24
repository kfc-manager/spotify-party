import { SongProps } from "./types";
import "./styles.css";

const Song = ({ song, addSong }: SongProps): JSX.Element => {
  const transformArtists = (artists: string[]): string => {
    let result: string = "";
    if (artists.length < 1) return result;
    artists.map((artist: string) => {
      result += artist + ", ";
    });

    return result.slice(0, -2);
  };

  const transformDuration = (duration: number): string => {
    let minutes: number = Math.floor(duration / 60000);
    let seconds: string = (duration % 60000) + "";

    while (seconds.length < 5) {
      seconds = "0" + seconds;
    }

    return minutes + ":" + seconds.slice(0, 2);
  };

  const handleClick = (): void => {
    addSong && addSong(song.id);
  };

  return (
    <div className="song">
      <img src={song.image_url} width={80} height={80} />
      <div className="song-text-info">
        <div className="song-name">{song.name}</div>
        {song.artists && (
          <div className="song-artists">{transformArtists(song.artists)}</div>
        )}
      </div>
      <div className="song-end-container">
        {addSong ? (
          <button className="song-add-button" onClick={handleClick}>
            <svg
              stroke="currentColor"
              fill="currentColor"
              strokeWidth="0"
              viewBox="0 0 1024 1024"
              height="1em"
              width="1em"
              style={{ fontSize: "24px" }}
            >
              <path d="M482 152h60q8 0 8 8v704q0 8-8 8h-60q-8 0-8-8V160q0-8 8-8Z"></path>
              <path d="M192 474h672q8 0 8 8v60q0 8-8 8H160q-8 0-8-8v-60q0-8 8-8Z"></path>
            </svg>
          </button>
        ) : (
          <div className="song-duration">
            {transformDuration(song.duration_ms)}
          </div>
        )}
      </div>
    </div>
  );
};

export default Song;
