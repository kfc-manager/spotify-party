import "./styles.css";
import { ChangeEvent, useEffect, useState } from "react";
import { NavbarProps } from "./types";

const Navbar = ({ query, setQuery }: NavbarProps): JSX.Element => {
  const [home, setHome] = useState<boolean>(true);

  const handleChange = (event: ChangeEvent<HTMLInputElement>): void => {
    setQuery(event.target.value);
  };

  useEffect((): void => {
    query.length > 0 && setHome(false);
    query.length < 1 && setHome(true);
  }, [query]);

  return (
    <div className="navbar">
      <div className="navbar-controls">
        <div className={`navbar-link${home ? " navbar-link-selected" : ""}`}>
          <button
            style={{
              display: "flex",
              alignItems: "center",
              color: "inherit",
              backgroundColor: "transparent",
              fontSize: "32px",
            }}
            onClick={() => setQuery("")}
          >
            <svg
              stroke="currentColor"
              fill="none"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
              aria-hidden="true"
              height="1em"
              width="1em"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5"
              ></path>
            </svg>
          </button>
        </div>
        <input
          type="text"
          placeholder="Search song.."
          onChange={handleChange}
          value={query}
          className="navbar-search-input"
        />
      </div>
    </div>
  );
};

export default Navbar;
