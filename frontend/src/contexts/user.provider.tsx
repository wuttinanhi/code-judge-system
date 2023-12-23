import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { UserLoginResponse } from "../apis/user";

type UserDataType = {
  accessToken: string;
  displayName: string;
  email: string;
};

interface UserContextType {
  user: UserDataType | undefined;
  setUser: (user: UserDataType | undefined) => void;
}

const UserContext = createContext<UserContextType>({
  setUser: () => {},
  user: undefined,
});

export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<UserDataType | undefined>(undefined);

  // try load accessToken from local storage
  useEffect(() => {
    const accessToken = localStorage.getItem("accessToken");
    const userData = localStorage.getItem("user");

    if (userData && accessToken) {
      const parsedUserData = JSON.parse(userData) as UserLoginResponse;

      setUser({
        accessToken,
        displayName: parsedUserData.displayname,
        email: parsedUserData.email,
      });
    }
  }, []);

  return (
    <UserContext.Provider value={{ user, setUser }}>
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
