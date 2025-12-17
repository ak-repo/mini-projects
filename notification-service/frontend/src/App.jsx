import React, { useState, useEffect, useCallback } from "react";
import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";

// --- CONFIGURATION ---
const BACKEND_URL = "http://localhost:8080";
const WS_URL = "ws://localhost:8080";
// ⚠️ IMPORTANT: Replace with a real user ID or pull from auth context
const USER_ID = "user_123";

// ⚠️ IMPORTANT: Replace with your actual Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyDy1E-U_8Tl0H6CcB6pL7S0JLSUEareWB8",
  authDomain: "streamhub-fcm.firebaseapp.com",
  projectId: "streamhub-fcm",
  storageBucket: "streamhub-fcm.firebasestorage.app",
  messagingSenderId: "414322346736",
  appId: "1:414322346736:web:c9b737bf6e78a62cc0d035",
  measurementId: "G-T6B4W2LXMK",
};

// ⚠️ IMPORTANT: Replace with your actual VAPID key from Firebase Console
const VAPID_KEY =
  "BKW1TjDbTXd5v5NsSPG4oJuxwfpp3xo7H9OSwqV9Uhy3AaTbr3EXR7VnV1ILlKdogqZbQ5yA2MpQKv7vleWGPO0";

// --- Firebase Initialization ---
const app = initializeApp(firebaseConfig);
const messaging = getMessaging(app);

// --- Token Registration Function ---
const requestAndRegisterToken = async (userId) => {
  try {
    const currentToken = await getToken(messaging, { vapidKey: VAPID_KEY });
    if (currentToken) {
      console.log("FCM Token registered:", currentToken.slice(0, 10) + "...");
      // Send token to the Go backend API
      await fetch(`${BACKEND_URL}/api/tokens`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: userId,
          token: currentToken,
          platform: "web",
        }),
      });
    }
  } catch (err) {
    console.error("Token retrieval failed. Push notifications disabled.", err);
  }
};

const createNotification = async () => {
  try {
    const payload = {
      type: "custom_alert", // You can choose any type
      user_id: USER_ID,
      data: { info: "This is a custom notification from the UI!" },
    };

    const response = await fetch(`${BACKEND_URL}/api/debug/publish`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!response.ok) throw new Error("Failed to create notification");

    const result = await response.json();
    console.log("Notification published:", result);
  } catch (err) {
    console.error("Error creating notification:", err);
  }
};

// =========================================================================
//                          APP COMPONENT
// =========================================================================

function App() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);

  // --- 1. HTTP API: Fetch Notification History ---
  const fetchNotifications = useCallback(async () => {
    try {
      const response = await fetch(
        `${BACKEND_URL}/api/notifications?user_id=${USER_ID}`
      );
      if (!response.ok) throw new Error("Failed to fetch");
      const data = await response.json();
      setNotifications(data);
    } catch (error) {
      console.error("Failed to fetch notifications:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  // --- 2. HTTP API: Mark Notification as Read (PATCH) ---
  const markAsRead = async (id) => {
    try {
      // Send PATCH request to the Go backend
      await fetch(`${BACKEND_URL}/api/notifications/${id}/read`, {
        method: "PATCH",
      });

      // Optimistically update the local state
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
      );
    } catch (error) {
      console.error("Failed to mark read:", error);
    }
  };

  // --- 3. Effect Hooks: Initialization & Realtime Setup ---
  useEffect(() => {
    // A. Initial Data Load
    fetchNotifications();

    // B. Register FCM Token (starts the push notification capability)
    requestAndRegisterToken(USER_ID);

    // C. Handle Foreground FCM Messages (when app is open)
    const unsubscribeFCM = onMessage(messaging, (payload) => {
      console.log("FCM Foreground Message Received:", payload);
      alert(`FCM: ${payload.notification.title}`);
      // If the notification is received via FCM, add it to the list
      setNotifications((prev) => [
        {
          id: payload.data.notification_id || Date.now(),
          user_id: USER_ID,
          title: payload.notification.title,
          body: payload.notification.body,
          is_read: false,
          created_at: new Date().toISOString(),
        },
        ...prev,
      ]);
    });

    // --- D. WebSocket Connection Setup ---
    const ws = new WebSocket(`${WS_URL}/ws/notifications?user_id=${USER_ID}`);

    ws.onopen = () => console.log("[WS] Connected to backend");

    ws.onmessage = (event) => {
      console.log("[WS] Realtime Message Received");
      const newNotif = JSON.parse(event.data);
      // Add new notification to the top of the list for real-time update
      setNotifications((prev) => [newNotif, ...prev]);
    };

    ws.onclose = () => console.log("[WS] Connection closed.");
    ws.onerror = (err) => console.error("[WS] Error:", err);

    // E. Cleanup
    return () => {
      unsubscribeFCM();
      ws.close();
    };
  }, [fetchNotifications]); // Dependency array ensures fetchNotifications is available

  // --- Render ---
  return (
    <div
      style={{
        padding: "20px",
        maxWidth: "600px",
        margin: "auto",
        fontFamily: "sans-serif",
      }}
    >
      <h1>🔔 Notification Center (User: {USER_ID})</h1>
      <p style={{ display: "flex", gap: "10px", alignItems: "center" }}>
        **Status:** {loading ? "Loading..." : "Live"} |
        <button onClick={fetchNotifications}>Refresh History</button>
      </p>
      <hr />
      <button
        onClick={createNotification}
        style={{
          padding: "10px 15px",
          backgroundColor: "#28a745",
          color: "white",
          border: "none",
          borderRadius: "5px",
          cursor: "pointer",
          marginBottom: "15px",
        }}
      >
        Create Notification
      </button>

      {notifications.map((n) => (
        <div
          key={n.id}
          style={{
            border: "1px solid #ddd",
            padding: "15px",
            margin: "10px 0",
            borderRadius: "8px",
            backgroundColor: n.is_read ? "#f4f4f4" : "#fffbe6",
            boxShadow: n.is_read ? "none" : "0 2px 5px rgba(255, 165, 0, 0.2)",
          }}
        >
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <h3
              style={{
                margin: "0 0 5px 0",
                color: n.is_read ? "#555" : "#333",
              }}
            >
              {n.title}{" "}
              {!n.is_read && (
                <span style={{ color: "red", fontSize: "0.9em" }}>• NEW</span>
              )}
            </h3>
            {!n.is_read && (
              <button
                onClick={() => markAsRead(n.id)}
                style={{
                  padding: "5px 10px",
                  backgroundColor: "#007bff",
                  color: "white",
                  border: "none",
                  borderRadius: "4px",
                  cursor: "pointer",
                }}
              >
                Mark Read
              </button>
            )}
          </div>
          <p style={{ margin: "0 0 10px 0", color: "#666" }}>{n.body}</p>
          <small style={{ color: "#999" }}>
            Received: {new Date(n.created_at).toLocaleString()}
          </small>
        </div>
      ))}

      {notifications.length === 0 && !loading && <p>No notifications yet.</p>}
    </div>
  );
}

export default App;
