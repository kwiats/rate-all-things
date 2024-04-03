import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode"; // Ensure correct import for jwtDecode

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);

  const navigate = useNavigate(); // Hook to get the navigate function

  const handleLogin = async (e) => {
    e.preventDefault();

    const loginUrl = "http://localhost:8080/api/auth/login";

    try {
      const response = await fetch(loginUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username: username, password: password }),
      });

      if (!response.ok) {
        throw new Error("Login failed");
      }

      const { token } = await response.json();
      const decoded = jwtDecode(token); 

      if (decoded.sub) {
        localStorage.setItem("user_id", decoded.sub);
      }
      localStorage.setItem("token", token);

      navigate("/"); 
    } catch (error) {
      setError(error.message);
    }
  };

  return (
    <div className="login-container bg-gray-800 min-h-screen flex justify-center items-center">
      <form
        onSubmit={handleLogin}
        className="bg-gray-700 p-10 rounded-lg shadow-lg"
      >
        {error && <p className="text-red-500">{error}</p>}
        <div className="mb-4">
          <label
            htmlFor="username"
            className="block text-white text-sm font-bold mb-2"
          >
            Username
          </label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>
        <div className="mb-6">
          <label
            htmlFor="password"
            className="block text-white text-sm font-bold mb-2"
          >
            Password
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Login
        </button>
      </form>
    </div>
  );
};

export default Login;
