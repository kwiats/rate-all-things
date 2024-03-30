import { Message } from "./message";
import { User } from "./user";

export class Chat {
  id: number = 0;
  users: User[] = [];
  last_message: Message = new Message();
}
