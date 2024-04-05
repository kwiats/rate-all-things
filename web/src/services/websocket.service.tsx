export const initializeWebSocket = (token: string | null) => {
  const webSocket = new WebSocket("ws://localhost:8080/ws/");
  webSocket.onopen = () => {
    console.log("WebSocket Connected");
    if (token) {
      webSocket.send(JSON.stringify({ token: `Bearer ${token}` }));
    }
  };
  webSocket.onerror = (error) => {
    console.error("WebSocket Error: ", error);
  };
  webSocket.onclose = () => {
    console.log("WebSocket Disconnected");
  };
  return webSocket;
};
