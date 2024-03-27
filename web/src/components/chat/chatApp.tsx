import React, { useState, useEffect } from "react";
import { useParams, useLocation } from "react-router-dom";
import axios from "axios";
import ChatList from "./chatList";
import ChatRoom from "./chatRoom";
import { jwtDecode } from "jwt-decode";
import { Message } from "../../models/message";
import { Chat } from "../../models/chats";

function ChatApp() {
  const [chats, setChats] = useState<Chat[]>([]);
  const [selectedChatId, setSelectedChatId] = useState<number>();

  const [sub, setSub] = useState(null);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const token = useParams();

  useEffect(() => {
    console.log(chats);
    if (token) {
      const decoded = jwtDecode(token.token);
      if (decoded.sub) {
        setSub(decoded.sub);
        fetchChats(decoded.sub);
      }
    }

    const webSocket = new WebSocket("ws://localhost:8080/ws/");

    webSocket.onopen = () => {
      console.log("WebSocket Connected");
      if (token) {
        webSocket.send(JSON.stringify({ token: `Bearer ${token.token}` }));
      }
    };

    webSocket.onmessage = (event) => {
      const message: Message = JSON.parse(event.data);
      console.log("Message from server ", message);
      console.log(chats);
      updateChats(message);
    };

    webSocket.onerror = (error) => {
      console.error("WebSocket Error: ", error);
    };

    webSocket.onclose = () => {
      console.log("WebSocket Disconnected");
    };

    setWs(webSocket);

    return () => {
      webSocket.close();
    };
  }, [token]);

  const fetchChats = async (userId: number) => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api/chat/?userId=${userId}`
      );
      const chats: Chat[] = response.data;
      console.log("Fetch chats", chats);
      setChats(chats);
    } catch (error) {
      console.error("Failed to fetch chats", error);
    }
  };

  const updateChats = (message: Message) => {
    console.log("Before update", chats);

    const updatedChats = chats.map((chat: Chat) => {
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
    console.log("After update", updatedChats);

    setChats(updatedChats);
  };

  const handleSelectChat = (chatId: number) => {
    setSelectedChatId(chatId);
  };

  return (
    <div>
      {chats && <ChatList chats={chats} onSelectChat={handleSelectChat} />}
      {selectedChatId && <ChatRoom chatId={selectedChatId} ws={ws} sub={sub} />}
    </div>
  );
}

export default ChatApp;
