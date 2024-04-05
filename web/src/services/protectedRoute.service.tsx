import { Navigate } from "react-router-dom";
import { isUserLoggedIn } from "./auth.service";

const ProtectedRoute = ({ children }) => {
  const isLoggedIn = isUserLoggedIn();
  return isLoggedIn ? children : <Navigate to="/login" />;
};

export default ProtectedRoute;
