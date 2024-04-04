import React, { useState, useEffect } from "react";
import axios from "axios";
import { Message } from "../../../models/message";

function ChatRoom({ chatId, ws }) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [userId, setUserId] = useState<number | null>(null); // State to hold the user ID

  useEffect(() => {
    const storedUserId = localStorage.getItem("user_id");
    if (storedUserId) {
      setUserId(parseInt(storedUserId));
    }

    fetchMessages();

    const handleMessage = (event: { data: string }) => {
      const data: Message = JSON.parse(event.data);
      if (data.chat === chatId) {
        setMessages((prevMessages) => [...prevMessages, data]);
      }
    };

    ws.addEventListener("message", handleMessage);

    return () => ws.removeEventListener("message", handleMessage);
  }, [chatId, ws]);

  const fetchMessages = async () => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api/chat/${chatId}/messages`
      );
      const sortedMessages = sortMessagesByCreatedAt(response.data);
      setMessages(sortedMessages);
    } catch (error) {
      console.error("Failed to fetch messages", error);
    }
  };

  const sortMessagesByCreatedAt = (messages: Message[]) => {
    return messages.sort(
      (a: Message, b: Message) =>
        new Date(a.created_at) - new Date(b.created_at)
    );
  };

  const sendMessage = () => {
    if (ws && ws.readyState === WebSocket.OPEN && inputMessage.trim()) {
      const messageToSend = JSON.stringify({
        chat: chatId,
        content: inputMessage,
        sender: userId, // Add sender information based on the user ID from Local Storage
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
    <div className="flex flex-col max-h-[calc(100%-150px)] h-full">
      {" "}
      {/* Dostosuj wysokość do swoich potrzeb */}
      <h2 className="text-2xl font-semibold text-gray-800 p-4">
        Chat Room: {chatId}
      </h2>
      <div className="flex-grow overflow-auto p-4 space-y-2 bg-white rounded-lg shadow">
        <ul>
          {messages.map((message, index) => (
            <li
              key={index}
              className={`flex ${
                message.sender === userId ? "justify-end" : "justify-start"
              }`}
            >
              <div
                className={`max-w-xs p-3 rounded-lg ${
                  message.sender === userId ? "bg-blue-100" : "bg-gray-100"
                }`}
              >
                <strong className="font-semibold text-gray-800">
                  {message.sender === userId ? "You" : `User ${message.sender}`}
                </strong>
                <p className="text-gray-800">{message.content}</p>
                <small className="text-xs text-gray-500">
                  {formatDate(message.created_at)}
                </small>
              </div>
            </li>
          ))}
        </ul>
      </div>
      <div className="p-4 bg-white">
        <input
          className="w-full p-3 border-2 border-gray-300 bg-gray-50 text-gray-800 rounded-lg focus:outline-none focus:border-blue-500"
          value={inputMessage}
          onChange={(e) => setInputMessage(e.target.value)}
          type="text"
          placeholder="Type your message here..."
          style={{ boxShadow: "0 2px 5px rgba(0, 0, 0, 0.1)" }}
        />
        <button
          className="mt-3 w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          onClick={sendMessage}
        >
          Send
        </button>
      </div>
    </div>
  );
}

export default ChatRoom;
