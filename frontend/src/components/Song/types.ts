import { APISong } from "../../types";

type SongProps = {
  song: APISong;
  addSong?: (id: string) => void;
};

export type { SongProps };
