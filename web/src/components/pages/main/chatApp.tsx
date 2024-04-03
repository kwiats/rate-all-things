import React, { useState, useEffect } from "react";
import axios from "axios";
import { Chat } from "../../../models/chats";
import { Message } from "../../../models/message";
import ChatList from "./chatList";
import ChatRoom from "./chatRoom";

function ChatApp() {
  const [chats, setChats] = useState<Chat[]>([]);
  const [selectedChatId, setSelectedChatId] = useState<number>();
  const [userId, setUserId] = useState<number>();
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const user_id_str = localStorage.getItem("user_id");
    if (user_id_str) {
      const userId = parseInt(user_id_str);
      setUserId(userId);
      fetchChats(userId);
    }

    const webSocket = new WebSocket("ws://localhost:8080/ws/");
    webSocket.onopen = () => {
      console.log("WebSocket Connected");
      const token = localStorage.getItem("token");
      if (token) {
        webSocket.send(JSON.stringify({ token: `Bearer ${token}` }));
      }
    };
    setWs(webSocket);

    return () => webSocket.close();
  }, []);

  useEffect(() => {
    if (ws) {
      ws.onmessage = (event) => {
        const message: Message = JSON.parse(event.data);
        console.log("Message from server ", message);
        updateChats(message);
      };

      ws.onerror = (error) => {
        console.error("WebSocket Error: ", error);
      };

      ws.onclose = () => {
        console.log("WebSocket Disconnected");
      };
    }
  }, [ws]);

  const fetchChats = async (userId: number) => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api/chat/?userId=${userId}`
      );
      setChats(response.data);
    } catch (error) {
      console.error("Failed to fetch chats", error);
    }
  };

  const updateChats = (message: Message) => {
    setChats((currentChats) => {
      const updatedChats = currentChats.map((chat) => {
        const chatLastMessageDate = new Date(chat.last_message.created_at);
        const incomingMessageDate = new Date(message.created_at);
        if (
          chat.id === message.chat &&
          chatLastMessageDate < incomingMessageDate
        ) {
          return { ...chat, last_message: message };
        }
        return chat;
      });
      return updatedChats;
    });
  };

  const handleSelectChat = (chatId: number) => {
    setSelectedChatId(chatId);
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-4xl bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="flex items-center justify-between bg-gray-200 p-4">
          <button className="w-3 h-3 bg-red-500 rounded-full focus:outline-none"></button>
          <span>#{userId}</span>
          <div></div>
        </div>
        <div className="flex flex-col h-[calc(100vh-15rem)]">
          <ChatList chats={chats} onSelectChat={handleSelectChat} />
          {selectedChatId && <ChatRoom chatId={selectedChatId} ws={ws} />}
        </div>
      </div>
    </div>
  );
}

export default ChatApp;
