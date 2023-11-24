import "./styles.css";

const LoadingProgress = (): JSX.Element => {
  return (
    <span className="loading-progress">
      <svg className="loading-progress-svg" viewBox="22 22 44 44">
        <circle
          className="loading-progress-svg-circle"
          cx="44"
          cy="44"
          r="20.2"
          fill="none"
          strokeWidth="3.6"
        />
      </svg>
    </span>
  );
};

export default LoadingProgress;
