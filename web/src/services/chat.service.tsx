import axios from "axios";

const API_URL = "http://localhost:8080/";

export const fetchChats = async (userId: number) => {
  try {
    const response = await axios.get(`${API_URL}api/chat/?userId=${userId}`);
    return response.data;
  } catch (error) {
    console.error("Failed to fetch chats", error);
    throw error;
  }
};
