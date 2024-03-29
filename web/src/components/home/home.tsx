import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

function Home() {
  const [chatId, setChatId] = useState("");
  const [token, setToken] = useState("");
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    navigate(`/chat/${token}/`);
  };

  return (
    <div>
      <h1>Welcome</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={token}
          onChange={(e) => setToken(e.target.value)}
          placeholder="Token"
          required
        />
        <button type="submit">Join Chat</button>
      </form>
    </div>
  );
}

export default Home;
