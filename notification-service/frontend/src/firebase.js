import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";
import { getAnalytics } from "firebase/analytics";

// Your web app's Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyDy1E-U_8Tl0H6CcB6pL7S0JLSUEareWB8",
  authDomain: "streamhub-fcm.firebaseapp.com",
  projectId: "streamhub-fcm",
  storageBucket: "streamhub-fcm.firebasestorage.app",
  messagingSenderId: "414322346736",
  appId: "1:414322346736:web:c9b737bf6e78a62cc0d035",
  measurementId: "G-T6B4W2LXMK",
};

// 1. Initialize App FIRST
const app = initializeApp(firebaseConfig);

// 2. Then initialize Analytics and Messaging
const analytics = getAnalytics(app);
export const messaging = getMessaging(app);

export const requestForToken = async (userId) => {
  try {
    // 3. IMPORTANT: Replace this string with your actual Key Pair from Firebase Console
    // Go to Project Settings -> Cloud Messaging -> Web Configuration -> Generate Key Pair
    const currentToken = await getToken(messaging, {
      vapidKey: "YOUR_ACTUAL_VAPID_KEY_HERE",
    });

    if (currentToken) {
      // Send token to backend
      await fetch("http://localhost:8080/api/tokens", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: userId,
          token: currentToken,
          platform: "web",
        }),
      });
      console.log("FCM Token registered");
    }
  } catch (err) {
    console.log("An error occurred while retrieving token. ", err);
  }
};