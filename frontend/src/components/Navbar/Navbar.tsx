import { Link, useLocation, useNavigate } from "react-router-dom";
import "./styles.css";
import { ChangeEvent, useEffect, useState } from "react";
import { NavbarProps } from "./types";

const Navbar = ({ query, setQuery }: NavbarProps): JSX.Element => {
  const location = useLocation();
  const navigate = useNavigate();
  const [home, setHome] = useState<boolean>(true);

  const handleChange = (event: ChangeEvent<HTMLInputElement>): void => {
    setQuery(event.target.value);
  };

  useEffect((): void => {
    location.pathname === "/" && query.length > 0 && setQuery("");
    location.pathname === "/" && setHome(true);
    location.pathname !== "/" && setHome(false);
  }, [location]);

  useEffect((): void => {
    if (query.length < 1) {
      navigate("/");
    } else {
      navigate("/search");
    }
  }, [query]);

  return (
    <div className="navbar">
      <div className="navbar-controls">
        <div className={`navbar-link${home ? " navbar-link-selected" : ""}`}>
          <Link
            style={{ display: "flex", alignItems: "center", color: "inherit" }}
            to="/"
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
          </Link>
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
