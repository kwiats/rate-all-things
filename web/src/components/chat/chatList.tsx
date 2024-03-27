import { useEffect } from "react";
import { Chat } from "../../models/chats";

function ChatList({ chats, onSelectChat }) {
  useEffect(() => {
    console.log(chats);
  });
  return (
    <div>
      <h2>Chat List</h2>
      <ul>
        {chats.map((chat: Chat) => (
          <li key={chat.id}>
            <button onClick={() => onSelectChat(chat.id)}>
              {chat.id} {chat.last_message.content} {chat.last_message.sender}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ChatList;
