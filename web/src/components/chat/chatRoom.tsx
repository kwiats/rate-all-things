import React, { useState, useEffect } from "react";
import axios from "axios";
import { Message } from "../../models/message";

function ChatRoom({ chatId, ws, sub }) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState("");

  useEffect(() => {
    fetchMessages();
    const handleMessage = (event: { data: string }) => {
      const data: Message = JSON.parse(event.data);
      setMessages((prevMessages) => [...prevMessages, data]);
    };

    ws.addEventListener("message", handleMessage);
    return () => ws.removeEventListener("message", handleMessage);
  }, [chatId, ws]);

  const fetchMessages = async () => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api/chat/${chatId}/messages`
      );
      setMessages(response.data);
    } catch (error) {
      console.error("Failed to fetch messages", error);
    }
  };

  const sendMessage = () => {
    if (ws && ws.readyState === WebSocket.OPEN && inputMessage.trim()) {
      const messageToSend = JSON.stringify({
        chat: chatId,
        content: inputMessage,
      });
      ws.send(messageToSend);
      setInputMessage("");
    }
  };
  const formatDate = (dateString: string) => {
    const options = {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  return (
    <div>
      <h2>Chat Room: {chatId}</h2>
      <div style={{ height: "300px", overflowY: "scroll" }}>
        <ul>
          {messages.map((message, index) => (
            <li
              key={index}
              style={{ textAlign: message.sender === sub ? "right" : "left" }}
            >
              <div>
                <strong>
                  {message.sender === sub ? "You" : `User ${message.sender}`}
                </strong>
                <p>{message.content}</p>
                <small>{formatDate(message.created_at)}</small>
              </div>
            </li>
          ))}
        </ul>
      </div>
      <input
        value={inputMessage}
        onChange={(e) => setInputMessage(e.target.value)}
        type="text"
        placeholder="Type your message here..."
      />
      <button onClick={sendMessage}>Send</button>
    </div>
  );
}

export default ChatRoom;
