import Queue from "./components/Queue/Queue";
import Search from "./components/Search/Search";
import Navbar from "./components/Navbar/Navbar";
import { Fragment, useState } from "react";

const App = (): JSX.Element => {
  const [query, setQuery] = useState<string>("");

  return (
    <Fragment>
      <Navbar query={query} setQuery={setQuery} />
      {query.length < 1 ? <Queue /> : <Search query={query} />}
    </Fragment>
  );
};

export default App;
