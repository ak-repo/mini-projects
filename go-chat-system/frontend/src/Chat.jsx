import React, { useState, useEffect, useRef } from 'react';
import { MessageCircle, Send, Phone, Video, Users, LogOut, User, Settings, Search, Paperclip, Smile, MoreVertical, Check, CheckCheck } from 'lucide-react';

const API_URL = 'http://localhost:8080/api';
const WS_URL = 'ws://localhost:8080/api/ws';

export default function ChatApp() {
  const [currentView, setCurrentView] = useState('login'); // login, register, chat
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [conversations, setConversations] = useState([]);
  const [activeConversation, setActiveConversation] = useState(null);
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState('');
  const [ws, setWs] = useState(null);
  const [onlineUsers, setOnlineUsers] = useState(new Set());
  const [typingUsers, setTypingUsers] = useState(new Set());
  const [searchQuery, setSearchQuery] = useState('');
  const messagesEndRef = useRef(null);
  const wsRef = useRef(null);

  // Auth handlers
  const handleLogin = async (email, password) => {
    try {
      const res = await fetch(`${API_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      const data = await res.json();
      if (res.ok) {
        setToken(data.token);
        setUser(data.user);
        localStorage.setItem('token', data.token);
        setCurrentView('chat');
        connectWebSocket(data.token);
      } else {
        alert(data.error || 'Login failed');
      }
    } catch (err) {
      alert('Connection error');
    }
  };

  const handleRegister = async (username, email, password, displayName) => {
    try {
      const res = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password, display_name: displayName })
      });
      const data = await res.json();
      if (res.ok) {
        setToken(data.token);
        setUser(data.user);
        localStorage.setItem('token', data.token);
        setCurrentView('chat');
        connectWebSocket(data.token);
      } else {
        alert(data.error || 'Registration failed');
      }
    } catch (err) {
      alert('Connection error');
    }
  };

  const handleLogout = () => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    setToken(null);
    setUser(null);
    setConversations([]);
    setMessages([]);
    setActiveConversation(null);
    localStorage.removeItem('token');
    setCurrentView('login');
  };

  // WebSocket connection
  const connectWebSocket = (authToken) => {
    const websocket = new WebSocket(`${WS_URL}?token=${authToken}`);
    
    websocket.onopen = () => {
      console.log('WebSocket connected');
    };

    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      handleWebSocketMessage(data);
    };

    websocket.onclose = () => {
      console.log('WebSocket disconnected');
      setTimeout(() => connectWebSocket(authToken), 3000);
    };

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    wsRef.current = websocket;
    setWs(websocket);
  };

  const handleWebSocketMessage = (data) => {
    switch (data.type) {
      case 'new_message':
        setMessages(prev => [...prev, {
          id: data.payload.id,
          content: data.payload.content,
          sender_id: data.payload.sender_id,
          created_at: data.payload.created_at,
          status: 'sent'
        }]);
        break;
      case 'typing_indicator':
        if (data.payload.is_typing) {
          setTypingUsers(prev => new Set(prev).add(data.payload.user_id));
        } else {
          setTypingUsers(prev => {
            const next = new Set(prev);
            next.delete(data.payload.user_id);
            return next;
          });
        }
        break;
      case 'presence_update':
        if (data.payload.status === 'online') {
          setOnlineUsers(prev => new Set(prev).add(data.payload.user_id));
        } else {
          setOnlineUsers(prev => {
            const next = new Set(prev);
            next.delete(data.payload.user_id);
            return next;
          });
        }
        break;
    }
  };

  const sendMessage = () => {
    if (!messageInput.trim() || !ws || !activeConversation) return;

    const message = {
      type: 'chat_message',
      payload: {
        conversation_id: activeConversation.id,
        content: messageInput,
        message_type: 'text'
      }
    };

    ws.send(JSON.stringify(message));
    
    // Optimistic update
    setMessages(prev => [...prev, {
      id: `temp-${Date.now()}`,
      content: messageInput,
      sender_id: user.id,
      created_at: new Date().toISOString(),
      status: 'sending'
    }]);
    
    setMessageInput('');
  };

  const sendTypingIndicator = (isTyping) => {
    if (!ws || !activeConversation) return;
    ws.send(JSON.stringify({
      type: 'typing',
      payload: {
        conversation_id: activeConversation.id,
        is_typing: isTyping
      }
    }));
  };

  const loadMessages = async (conversationId) => {
    try {
      const res = await fetch(`${API_URL}/conversations/${conversationId}/messages?limit=50`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await res.json();
      setMessages(data.messages || []);
    } catch (err) {
      console.error('Failed to load messages');
    }
  };

  const createConversation = async (otherUserId) => {
    try {
      const res = await fetch(`${API_URL}/conversations`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          type: 'one_to_one',
          member_ids: [otherUserId]
        })
      });
      const data = await res.json();
      if (res.ok) {
        setConversations(prev => [...prev, data]);
        setActiveConversation(data);
        loadMessages(data.id);
      }
    } catch (err) {
      alert('Failed to create conversation');
    }
  };

  useEffect(() => {
    if (token && currentView === 'chat') {
      connectWebSocket(token);
    }
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [token, currentView]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  useEffect(() => {
    if (activeConversation) {
      loadMessages(activeConversation.id);
    }
  }, [activeConversation]);

  // Login View
  if (currentView === 'login') {
    return <LoginView onLogin={handleLogin} onSwitchToRegister={() => setCurrentView('register')} />;
  }

  // Register View
  if (currentView === 'register') {
    return <RegisterView onRegister={handleRegister} onSwitchToLogin={() => setCurrentView('login')} />;
  }

  // Main Chat View
  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <div className="w-80 bg-white border-r border-gray-200 flex flex-col">
        {/* Header */}
        <div className="p-4 border-b border-gray-200 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
              {user?.display_name?.[0]?.toUpperCase() || 'U'}
            </div>
            <div>
              <div className="font-semibold text-sm">{user?.display_name}</div>
              <div className="text-xs text-green-500">Online</div>
            </div>
          </div>
          <div className="flex gap-2">
            <button className="p-2 hover:bg-gray-100 rounded-full">
              <Settings size={20} className="text-gray-600" />
            </button>
            <button onClick={handleLogout} className="p-2 hover:bg-gray-100 rounded-full">
              <LogOut size={20} className="text-gray-600" />
            </button>
          </div>
        </div>

        {/* Search */}
        <div className="p-3 border-b border-gray-200">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 text-gray-400" size={18} />
            <input
              type="text"
              placeholder="Search conversations..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        {/* Conversations List */}
        <div className="flex-1 overflow-y-auto">
          {conversations.length === 0 ? (
            <div className="p-4 text-center text-gray-500">
              <MessageCircle size={48} className="mx-auto mb-2 text-gray-300" />
              <p className="text-sm">No conversations yet</p>
              <button
                onClick={() => createConversation('demo-user-id')}
                className="mt-3 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm hover:bg-blue-600"
              >
                Start New Chat
              </button>
            </div>
          ) : (
            conversations.map(conv => (
              <ConversationItem
                key={conv.id}
                conversation={conv}
                isActive={activeConversation?.id === conv.id}
                onClick={() => setActiveConversation(conv)}
                isOnline={onlineUsers.has(conv.other_user_id)}
              />
            ))
          )}
        </div>
      </div>

      {/* Main Chat Area */}
      {activeConversation ? (
        <div className="flex-1 flex flex-col">
          {/* Chat Header */}
          <div className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-4">
            <div className="flex items-center gap-3">
              <div className="relative">
                <div className="w-10 h-10 rounded-full bg-purple-500 flex items-center justify-center text-white font-semibold">
                  {activeConversation.name?.[0]?.toUpperCase() || 'C'}
                </div>
                {onlineUsers.has(activeConversation.other_user_id) && (
                  <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
                )}
              </div>
              <div>
                <div className="font-semibold">{activeConversation.name || 'Chat'}</div>
                {typingUsers.size > 0 && (
                  <div className="text-xs text-blue-500">typing...</div>
                )}
              </div>
            </div>
            <div className="flex gap-2">
              <button className="p-2 hover:bg-gray-100 rounded-full">
                <Phone size={20} className="text-gray-600" />
              </button>
              <button className="p-2 hover:bg-gray-100 rounded-full">
                <Video size={20} className="text-gray-600" />
              </button>
              <button className="p-2 hover:bg-gray-100 rounded-full">
                <MoreVertical size={20} className="text-gray-600" />
              </button>
            </div>
          </div>

          {/* Messages Area */}
          <div className="flex-1 overflow-y-auto p-4 bg-gray-50">
            {messages.map((msg, idx) => (
              <MessageBubble
                key={msg.id}
                message={msg}
                isOwn={msg.sender_id === user?.id}
                showAvatar={idx === 0 || messages[idx - 1].sender_id !== msg.sender_id}
              />
            ))}
            <div ref={messagesEndRef} />
          </div>

          {/* Input Area */}
          <div className="bg-white border-t border-gray-200 p-4">
            <div className="flex items-center gap-2">
              <button className="p-2 hover:bg-gray-100 rounded-full">
                <Paperclip size={20} className="text-gray-600" />
              </button>
              <input
                type="text"
                value={messageInput}
                onChange={(e) => {
                  setMessageInput(e.target.value);
                  sendTypingIndicator(true);
                  setTimeout(() => sendTypingIndicator(false), 1000);
                }}
                onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                placeholder="Type a message..."
                className="flex-1 px-4 py-2 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button className="p-2 hover:bg-gray-100 rounded-full">
                <Smile size={20} className="text-gray-600" />
              </button>
              <button
                onClick={sendMessage}
                disabled={!messageInput.trim()}
                className="p-3 bg-blue-500 text-white rounded-full hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Send size={20} />
              </button>
            </div>
          </div>
        </div>
      ) : (
        <div className="flex-1 flex items-center justify-center bg-gray-50">
          <div className="text-center">
            <MessageCircle size={64} className="mx-auto mb-4 text-gray-300" />
            <h3 className="text-xl font-semibold text-gray-700 mb-2">Welcome to Chat</h3>
            <p className="text-gray-500">Select a conversation to start messaging</p>
          </div>
        </div>
      )}
    </div>
  );
}

function LoginView({ onLogin, onSwitchToRegister }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        <div className="text-center mb-8">
          <MessageCircle size={48} className="mx-auto mb-4 text-blue-500" />
          <h1 className="text-3xl font-bold text-gray-800">Welcome Back</h1>
          <p className="text-gray-600 mt-2">Sign in to continue chatting</p>
        </div>
        <div className="space-y-4">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && onLogin(email, password)}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={() => onLogin(email, password)}
            className="w-full py-3 bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-600 transition"
          >
            Sign In
          </button>
        </div>
        <p className="text-center mt-6 text-gray-600">
          Don't have an account?{' '}
          <button onClick={onSwitchToRegister} className="text-blue-500 font-semibold hover:underline">
            Sign Up
          </button>
        </p>
      </div>
    </div>
  );
}

function RegisterView({ onRegister, onSwitchToLogin }) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [displayName, setDisplayName] = useState('');

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        <div className="text-center mb-8">
          <Users size={48} className="mx-auto mb-4 text-purple-500" />
          <h1 className="text-3xl font-bold text-gray-800">Create Account</h1>
          <p className="text-gray-600 mt-2">Join us and start chatting</p>
        </div>
        <div className="space-y-4">
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
          <input
            type="text"
            placeholder="Display Name"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
          <input
            type="password"
            placeholder="Password (min 8 characters)"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          />
          <button
            onClick={() => onRegister(username, email, password, displayName)}
            className="w-full py-3 bg-purple-500 text-white rounded-lg font-semibold hover:bg-purple-600 transition"
          >
            Sign Up
          </button>
        </div>
        <p className="text-center mt-6 text-gray-600">
          Already have an account?{' '}
          <button onClick={onSwitchToLogin} className="text-purple-500 font-semibold hover:underline">
            Sign In
          </button>
        </p>
      </div>
    </div>
  );
}

function ConversationItem({ conversation, isActive, onClick, isOnline }) {
  return (
    <div
      onClick={onClick}
      className={`p-4 border-b border-gray-100 cursor-pointer hover:bg-gray-50 ${
        isActive ? 'bg-blue-50' : ''
      }`}
    >
      <div className="flex items-center gap-3">
        <div className="relative">
          <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white font-semibold">
            {conversation.name?.[0]?.toUpperCase() || 'C'}
          </div>
          {isOnline && (
            <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
          )}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex justify-between items-baseline">
            <h3 className="font-semibold text-sm truncate">{conversation.name || 'Conversation'}</h3>
            <span className="text-xs text-gray-500">2m</span>
          </div>
          <p className="text-sm text-gray-600 truncate">Last message preview...</p>
        </div>
      </div>
    </div>
  );
}

function MessageBubble({ message, isOwn, showAvatar }) {
  return (
    <div className={`flex gap-2 mb-4 ${isOwn ? 'flex-row-reverse' : ''}`}>
      {showAvatar && !isOwn && (
        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-purple-400 to-pink-500 flex items-center justify-center text-white text-xs font-semibold flex-shrink-0">
          U
        </div>
      )}
      {!showAvatar && !isOwn && <div className="w-8" />}
      <div className={`max-w-xs lg:max-w-md ${isOwn ? 'items-end' : 'items-start'} flex flex-col`}>
        <div
          className={`px-4 py-2 rounded-2xl ${
            isOwn
              ? 'bg-blue-500 text-white rounded-br-none'
              : 'bg-white text-gray-800 rounded-bl-none shadow'
          }`}
        >
          <p className="break-words">{message.content}</p>
        </div>
        <div className="flex items-center gap-1 mt-1 px-2">
          <span className="text-xs text-gray-500">
            {new Date(message.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
          </span>
          {isOwn && (
            <span className="text-blue-500">
              {message.status === 'read' ? <CheckCheck size={14} /> : <Check size={14} />}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}