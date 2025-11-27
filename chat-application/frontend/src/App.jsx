import React, { useEffect, useState, useCallback, useMemo } from "react";
import { Send, MessageSquare, Loader, XCircle } from "lucide-react";

// NOTE: Ensure your backend WebSocket server is running at this address.
const GATEWAY_URL = "ws://10.153.140.105:8080/ws";



// Two distinct user tokens/IDs for simulation
const TEST_TOKEN_1 = "user1@example.com"; // Client 1's identity
const TEST_TOKEN_2 = "user2@example.com"; // Client 2's identity

// Custom Hook for WebSocket connection logic
const useChatSocket = (token, onMessageReceived) => {
  const [socket, setSocket] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!token) {
      setError("Token is missing. Cannot connect.");
      return;
    }

    // Set a flag to prevent setting state on unmounted components
    let isMounted = true; 

    const ws = new WebSocket(
      `${GATEWAY_URL}?token=${encodeURIComponent(token)}`
    );
    if(isMounted) {
      setSocket(ws);
      setError(null);
    }

    ws.onopen = () => {
      console.log(`WebSocket (${token}): Connection Open`);
      if(isMounted) {
        setIsConnected(true);
        setError(null);
      }
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "message") {
          onMessageReceived(data);
        } else if (data.error) {
          if(isMounted) {
            setError(`Server Error: ${data.error}`);
          }
        }
      } catch (e) {
        console.error(`WebSocket (${token}): Error parsing message:`, e);
      }
    };

    ws.onclose = (event) => {
      console.log(`WebSocket (${token}): Connection Closed`, event.code, event.reason);
      if(isMounted) {
        setIsConnected(false);
      }
    };

    ws.onerror = (e) => {
      console.error(`WebSocket (${token}): Error`, e);
      if(isMounted) {
        setError("WebSocket connection failed.");
      }
    };

    return () => {
      isMounted = false;
      ws.close();
    };
  }, [token, onMessageReceived]);

  const sendMessage = useCallback(
    (payload) => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(
          JSON.stringify({
            type: "message",
            ...payload,
          })
        );
      } else {
        console.warn("WebSocket not open. Message not sent:", payload);
      }
    },
    [socket]
  );

  return { isConnected, sendMessage, error };
};

// --- Generic Client Component ---
const ChatClient = ({ clientToken, defaultRecipientToken, clientName }) => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  // Recipient defaults to the other user's token
  const [recipient, setRecipient] = useState(defaultRecipientToken); 

  // Unique chat ID for the conversation pair
  const chatID = useMemo(() => {
    const users = [clientToken, defaultRecipientToken].sort();
    return `chat-${users[0]}-${users[1]}`;
  }, [clientToken, defaultRecipientToken]);

  const handleMessageReceived = useCallback((message) => {
    // Only display messages relevant to this client (either sent to or from them)
    if (message.to === clientToken || message.from === clientToken) {
        setMessages((prev) => [...prev, message]);
    }
  }, [clientToken]);

  const { isConnected, sendMessage, error } = useChatSocket(
    clientToken,
    handleMessageReceived
  );

  const handleSend = () => {
    if (!input.trim() || !isConnected) return;

    const newMessage = {
      client_id: `${clientToken}-${Date.now()}`, // Unique ID for sender's ACK
      from: clientToken,
      to: recipient,
      chat_id: chatID,
      text: input,
      msg_type: "text",
      created_at: Date.now(),
    };

    sendMessage(newMessage);
    setInput("");
  };

  const statusColor = isConnected ? "bg-green-500" : "bg-red-500";

  return (
    <div className="p-4 flex flex-col bg-white rounded-xl shadow-2xl h-full border border-gray-100">
      <header className="bg-blue-600 text-white p-4 rounded-t-xl shadow-lg flex items-center justify-between">
        <div className="flex items-center">
          <MessageSquare className="w-6 h-6 mr-3" />
          <h1 className="text-xl font-bold">{clientName} ({clientToken})</h1>
        </div>
        <div className="flex items-center text-sm">
          <span
            className={`w-3 h-3 rounded-full mr-2 ${statusColor} transition-colors`}
          ></span>
          {isConnected ? "Connected" : "Disconnected"}
        </div>
      </header>

      {error && (
        <div className="bg-red-100 text-red-700 p-3 flex items-center text-sm font-medium">
          <XCircle className="w-5 h-5 mr-2" />
          {error}
        </div>
      )}

      {/* Message History */}
      <div className="flex-1 overflow-y-auto bg-gray-50 p-4 space-y-3 border-x border-gray-100 min-h-[300px] max-h-[60vh]">
        {!isConnected && !error && (
          <div className="flex justify-center items-center h-full text-gray-500">
            <Loader className="w-5 h-5 animate-spin mr-2" />
            Connecting to Chat Service...
          </div>
        )}
        {messages.map((m, index) => (
          <div
            key={m.server_id || m.client_id || index}
            className={`flex ${
              m.from === clientToken ? "justify-end" : "justify-start"
            }`}
          >
            <div
              className={`max-w-[75%] p-3 rounded-xl shadow-md transition-all ${
                m.from === clientToken
                  ? "bg-blue-500 text-white rounded-br-none"
                  : "bg-white text-gray-800 rounded-tl-none border border-gray-200"
              }`}
            >
              <div className="font-semibold text-xs opacity-80 mb-1">
                {m.from === clientToken ? "You" : m.from}
              </div>
              <p className="text-sm break-words">{m.text}</p>
              <div className="text-xs mt-1 opacity-60 text-right">
                {new Date(m.created_at).toLocaleTimeString()}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Input and Recipient Selector */}
      <div className="bg-white p-4 border-t border-gray-200 rounded-b-xl shadow-inner">
        <div className="mb-3 text-sm flex items-center">
          <label className="font-semibold mr-2 text-gray-600">To:</label>
          <input
            type="email"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            className="flex-1 p-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 text-gray-700"
            placeholder="Recipient User ID"
          />
        </div>
        <div className="flex space-x-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && handleSend()}
            className="flex-1 p-3 border border-gray-300 rounded-xl focus:ring-blue-500 focus:border-blue-500 transition-shadow text-gray-700"
            placeholder="Type your message..."
            disabled={!isConnected}
          />
          <button
            onClick={handleSend}
            disabled={!isConnected || !input.trim()}
            className={`p-3 rounded-xl transition-all duration-200 shadow-lg ${
              isConnected && input.trim()
                ? "bg-blue-600 hover:bg-blue-700 text-white"
                : "bg-gray-400 text-gray-100 cursor-not-allowed"
            }`}
          >
            <Send className="w-6 h-6" />
          </button>
        </div>
      </div>
    </div>
  );
};


// --- Main App Component ---
export default function App() {
  return (
    <div className="p-4 bg-gray-100 min-h-screen font-[Inter]">
      <h2 className="text-2xl font-extrabold text-gray-800 text-center mb-6">
        Two-Client WebSocket Chat Simulation
      </h2>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 max-w-7xl mx-auto">
        {/* Client 1: user1@example.com talks to user2@example.com */}
        <ChatClient
          clientToken={TEST_TOKEN_1}
          defaultRecipientToken={TEST_TOKEN_2}
          clientName="Client 1"
        />

        {/* Client 2: user2@example.com talks to user1@example.com */}
        <ChatClient
          clientToken={TEST_TOKEN_2}
          defaultRecipientToken={TEST_TOKEN_1}
          clientName="Client 2"
        />
      </div>
      <footer className="text-center text-sm text-gray-500 mt-8">
        *Note: This demo requires a local WebSocket chat server running on 
        <code className="bg-gray-200 p-1 rounded">ws://localhost:8080/ws</code> to function.
      </footer>
    </div>
  );
}