import { Dispatch, SetStateAction } from "react";

type NavbarProps = {
  query: string;
  setQuery: Dispatch<SetStateAction<string>>;
};

export type { NavbarProps };
