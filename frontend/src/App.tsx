import useQueue from "./hooks/useQueue";

const App = (): JSX.Element => {
  const { data, loading, error } = useQueue();

  return <>{!loading && !error && data.map((song) => song.name)}</>;
};

export default App;
