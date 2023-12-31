import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { UserService } from "../apis/user";

type UserDataType = {
  accessToken: string;
  displayName: string;
  email: string;
  role: string;
};

interface UserContextType {
  user: UserDataType | undefined;
  setUser: (user: UserDataType | undefined) => void;
  logout: () => void;
}

const UserContext = createContext<UserContextType>({
  setUser: () => {},
  user: undefined,
  logout: () => {},
});

export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<UserDataType | undefined>(undefined);

  const logout = () => {
    localStorage.removeItem("accessToken");
    setUser(undefined);
  };

  // try load accessToken from local storage
  useEffect(() => {
    async function validateToken() {
      const accessToken = localStorage.getItem("accessToken");
      if (!accessToken) {
        return;
      }
      const res = await UserService.getUserInfo(accessToken);

      if (res.status === 200) {
        const data = await res.json();

        setUser({
          accessToken: accessToken as string,
          displayName: data.displayname,
          email: data.email,
          role: data.role,
        });
      } else {
        logout();
      }
    }

    validateToken();
  }, []);

  return (
    <UserContext.Provider value={{ user, setUser, logout }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);

  if (context === undefined) {
    throw new Error("useUser must be used within a UserProvider");
  }

  return context;
};
