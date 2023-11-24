import "./styles.css";
import { BackdropProps } from "./types";

const Backdrop = ({ children }: BackdropProps): JSX.Element => {
  return (
    <div className="backdrop">
      <div className="backdrop-content">{children}</div>
    </div>
  );
};

export default Backdrop;
