import api from "../api";

// Public Auth Routes
export const userLogin = async (loginData) => {
  try {
    const res = await api.post("auth/login", loginData);
    return res.data;
  } catch (error) {
    console.error("Login error:", error.message);
    throw error;
  }
};

export const userRegister = async (registerData) => {
  try {
    const res = await api.post("auth/register", registerData);
    return res.data;
  } catch (error) {
    console.error("Registration error:", error.message);
    throw error;
  }
};

// Protected Auth Routes
export const changePassword = async (passwordData) => {
  try {
    const res = await api.post("auth/password-change", passwordData);
    return res.data;
  } catch (error) {
    console.error("Password change error:", error.message);
    throw error;
  }
};

export const sendOTP = async (otpData) => {
  try {
    const res = await api.post("auth/send-otp", otpData);
    return res.data;
  } catch (error) {
    console.error("Send OTP error:", error.message);
    throw error;
  }
};

export const verifyOTP = async (otpData) => {
  try {
    const res = await api.post("auth/verify-otp", otpData);
    return res.data;
  } catch (error) {
    console.error("Verify OTP error:", error.message);
    throw error;
  }
};

export const getMe = async () => {
  try {
    const res = await api.get("/auth/me");
    return res.data;
  } catch (error) {
    console.error("Get user info error:", error.message);
    throw error;
  }
};
