import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import { Bell } from 'lucide-react';

const NotificationBell = ({ userId, socket }) => {
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    // 1. Load initial
    fetchNotifications();

    // 2. Listen for WebSocket messages
    if (socket) {
      socket.onmessage = (event) => {
        const newNotif = JSON.parse(event.data);
        setNotifications((prev) => [newNotif, ...prev]);
        setUnreadCount((prev) => prev + 1);
        
        // Optional: Trigger browser toast if window visible
        if (Notification.permission === 'granted') {
             new Notification(newNotif.title, { body: newNotif.body });
        }
      };
    }
  }, [socket]);

  const fetchNotifications = async () => {
    try {
      const res = await axios.get(`http://localhost:8080/api/notifications?user_id=${userId}`);
      setNotifications(res.data);
      setUnreadCount(res.data.filter(n => !n.is_read).length);
    } catch (err) {
      console.error(err);
    }
  };

  const markRead = async (id) => {
    try {
      await axios.patch(`http://localhost:8080/api/notifications/${id}/read`);
      setNotifications(prev => prev.map(n => n.id === id ? { ...n, is_read: true } : n));
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="relative">
      <button onClick={() => setIsOpen(!isOpen)} className="p-2 relative">
        <Bell size={24} />
        {unreadCount > 0 && (
          <span className="absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
            {unreadCount}
          </span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white shadow-xl rounded-lg border z-50 max-h-96 overflow-y-auto">
          <div className="p-2 border-b font-bold text-gray-700">Notifications</div>
          {notifications.length === 0 ? (
            <div className="p-4 text-gray-500 text-sm">No notifications</div>
          ) : (
            notifications.map(n => (
              <div 
                key={n.id} 
                onClick={() => !n.is_read && markRead(n.id)}
                className={`p-3 border-b text-sm cursor-pointer hover:bg-gray-50 ${!n.is_read ? 'bg-blue-50' : ''}`}
              >
                <div className="font-semibold">{n.title}</div>
                <div className="text-gray-600">{n.body}</div>
                <div className="text-xs text-gray-400 mt-1">{new Date(n.created_at).toLocaleString()}</div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
};

export default NotificationBell;