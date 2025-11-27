import React, { useEffect, useState, useCallback, useMemo } from 'react';
import { Send, MessageSquare, Loader, XCircle } from 'lucide-react';



const GATEWAY_URL = "ws://10.153.140.105:8080/ws";


// IMPORTANT: Replace with a token that your chat service's verifyToken method accepts (e.g., a test user ID)
const TEST_TOKEN = "user@example.com"; 

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

    const ws = new WebSocket(`${GATEWAY_URL}?token=${encodeURIComponent(token)}`);
    setSocket(ws);
    setError(null);

    ws.onopen = () => {
      console.log("WebSocket: Connection Open");
      setIsConnected(true);
      setError(null);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "message") {
          onMessageReceived(data);
        } else if (data.error) {
          setError(`Server Error: ${data.error}`);
        }
      } catch (e) {
        console.error("WebSocket: Error parsing message:", e);
      }
    };

    ws.onclose = (event) => {
      console.log("WebSocket: Connection Closed", event.code, event.reason);
      setIsConnected(false);
      // Implement robust reconnection logic here for production
    };

    ws.onerror = (e) => {
        console.error("WebSocket: Error", e);
        setError("WebSocket connection failed.");
    };

    return () => {
      ws.close();
    };
  }, [token, onMessageReceived]);

  const sendMessage = useCallback((payload) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({
        type: "message",
        ...payload
      }));
    } else {
      console.warn("WebSocket not open. Message not sent:", payload);
    }
  }, [socket]);

  return { isConnected, sendMessage, error };
};

// --- Main Chat Component ---

export default function App() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [recipient, setRecipient] = useState('friend@example.com');
  const chatID = useMemo(() => {
    // Simple chat ID based on two user IDs/emails
    const users = [TEST_TOKEN, recipient].sort();
    return `chat-${users[0]}-${users[1]}`;
  }, [recipient]);

  const handleMessageReceived = useCallback((message) => {
    setMessages((prev) => [...prev, message]);
  }, []);

  const { isConnected, sendMessage, error } = useChatSocket(
    TEST_TOKEN,
    handleMessageReceived
  );

  const handleSend = () => {
    if (!input.trim() || !isConnected) return;

    const newMessage = {
      client_id: `client-${Date.now()}`, // Unique ID for sender's ACK
      from: TEST_TOKEN,
      to: recipient,
      chat_id: chatID,
      text: input,
      msg_type: "text",
      created_at: Date.now(),
    };

    sendMessage(newMessage);
    setInput('');
  };

  const statusColor = isConnected ? 'bg-green-500' : 'bg-red-500';

  return (
    <div className="p-4 sm:p-6 md:p-8 max-w-2xl mx-auto font-[Inter] h-screen flex flex-col">
      <header className="bg-gray-800 text-white p-4 rounded-t-xl shadow-lg flex items-center justify-between">
        <div className="flex items-center">
          <MessageSquare className="w-6 h-6 mr-3" />
          <h1 className="text-xl font-bold">Real-Time Chat Demo</h1>
        </div>
        <div className="flex items-center text-sm">
          <span className={`w-3 h-3 rounded-full mr-2 ${statusColor} transition-colors`}></span>
          {isConnected ? 'Connected' : 'Disconnected'}
        </div>
      </header>

      {error && (
        <div className="bg-red-100 text-red-700 p-3 rounded-b-lg flex items-center">
          <XCircle className="w-5 h-5 mr-2" />
          {error}
        </div>
      )}

      {/* Message History */}
      <div className="flex-1 overflow-y-auto bg-gray-50 p-4 space-y-3 border-x border-gray-200">
        {!isConnected && !error && (
            <div className='flex justify-center items-center h-full text-gray-500'>
                <Loader className='w-5 h-5 animate-spin mr-2' />
                Connecting to Chat Service...
            </div>
        )}
        {messages.map((m, index) => (
          <div key={m.server_id || m.client_id || index} className={`flex ${m.from === TEST_TOKEN ? 'justify-end' : 'justify-start'}`}>
            <div className={`max-w-[75%] p-3 rounded-lg shadow-sm ${m.from === TEST_TOKEN ? 'bg-blue-500 text-white rounded-br-none' : 'bg-white text-gray-800 rounded-tl-none border border-gray-200'}`}>
              <div className="font-semibold text-xs opacity-70 mb-1">{m.from === TEST_TOKEN ? 'You' : m.from}</div>
              <p className="text-sm break-words">{m.text}</p>
              <div className="text-xs mt-1 opacity-50">{new Date(m.created_at).toLocaleTimeString()}</div>
            </div>
          </div>
        ))}
      </div>

      {/* Input and Recipient Selector */}
      <div className="bg-white p-4 border border-gray-200 rounded-b-xl shadow-inner">
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
            onKeyPress={(e) => e.key === 'Enter' && handleSend()}
            className="flex-1 p-3 border border-gray-300 rounded-xl focus:ring-blue-500 focus:border-blue-500 transition-shadow text-gray-700"
            placeholder="Type your message..."
            disabled={!isConnected}
          />
          <button
            onClick={handleSend}
            disabled={!isConnected || !input.trim()}
            className={`p-3 rounded-xl transition-all duration-200 ${
              isConnected && input.trim()
                ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-md'
                : 'bg-gray-400 text-gray-100 cursor-not-allowed'
            }`}
          >
            <Send className="w-6 h-6" />
          </button>
        </div>
      </div>
    </div>
  );
}