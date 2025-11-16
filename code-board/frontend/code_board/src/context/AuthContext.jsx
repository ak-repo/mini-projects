import { useState, useEffect } from "react";
import { AuthContext } from "./context";
import { getMe, userLogin, userRegister } from "../api/services/authService";

function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [isAuthenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        await getMe;
        setAuthenticated(true);
      } catch (error) {
        console.log("errr :", error);
      }
    })();
  }, []);

  useEffect(() => {
    if (user != null) {
      setAuthenticated(true);
    } else {
      setAuthenticated(false);
    }
  }, [user]);

  const login = async (email, password) => {
    const data = await userLogin(email, password);
    setUser(data.user);
  };

  const register = async (form) => {
    const data = await userRegister(form);
    setUser(data.user);
  };

  const logout = () => {
    // logoutService();
    setUser(null);
  };
  return (
    <AuthContext.Provider
      value={{
        user,
        setUser,
        login,
        logout,
        isAuthenticated,
        register,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export default AuthProvider;
