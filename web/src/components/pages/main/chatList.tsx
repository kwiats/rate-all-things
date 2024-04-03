import { Chat } from "../../../models/chats";

function ChatList({ chats, onSelectChat }) {
  function handleSearch(event: ChangeEvent<HTMLInputElement>): void {
    throw new Error("Function not implemented.");
  }

  return (
    <div className="flex items-center p-4">
      <input
        className="flex-grow p-2 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-blue-500"
        placeholder="Search chats..."
        type="text"
        onChange={handleSearch}
      />
      <div className="flex overflow-x-auto p-4 space-x-4">
        {chats.map((chat) => (
          <button
            key={chat.id}
            className="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center"
            onClick={() => onSelectChat(chat.id)}
          >
            <span>Chat {chat.id}</span>
          </button>
        ))}
      </div>
    </div>
  );
}

export default ChatList;
