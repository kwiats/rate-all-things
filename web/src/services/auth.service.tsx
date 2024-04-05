import { jwtDecode } from "jwt-decode";
interface TokenPayload {
  sub: string; // Subject (User ID)
  name: string; // Username
  exp: number; // Expiry
}
export const logout = (navigate: Function) => {
  localStorage.removeItem("user_id");
  localStorage.removeItem("token");
  navigate("/login");
};

export const isUserLoggedIn = (): boolean => {
  const token = localStorage.getItem("token");
  if (!token) {
    return false;
  }
  try {
    const decoded: any = jwtDecode(token);
    const currentTime = Date.now() / 1000;
    return decoded.exp > currentTime;
  } catch (error) {
    console.error("Problem with token decoding:", error);
    return false;
  }
};
export const getUserData = (): TokenPayload | null => {
  const token = localStorage.getItem("token");
  if (!token) return null;

  try {
    return jwtDecode<TokenPayload>(token);
  } catch (error) {
    console.error("Problem with token decoding:", error);
    return null;
  }
};

export const getUsername = (): string | undefined => {
  return getUserData()?.name;
};

export const getUserid = (): number | undefined => {
  const sub = getUserData()?.sub;
  return sub ? parseInt(sub) : undefined;
};

export const getToken = () => {
  return localStorage.getItem("token");
};
