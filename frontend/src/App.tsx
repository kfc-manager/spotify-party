import { BrowserRouter, Route, Routes } from "react-router-dom";
import Queue from "./components/Queue/Queue";
import Search from "./components/Search/Search";
import Navbar from "./components/Navbar/Navbar";
import { useState } from "react";

const App = (): JSX.Element => {
  const [query, setQuery] = useState<string>("");

  return (
    <BrowserRouter>
      <Navbar query={query} setQuery={setQuery} />
      <Routes>
        <Route path="/" element={<Queue />} />
        <Route path="/search" element={<Search query={query} />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
