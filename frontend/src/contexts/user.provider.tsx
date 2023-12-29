import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { API_URL } from "../apis/API_URL";
import { UserMeResponse } from "../types/user";

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

    // validate token
    const validateURL = API_URL + "/user/me";
    fetch(validateURL, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    })
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((json: UserMeResponse) => {
        setUser({
          accessToken: accessToken as string,
          displayName: json.DisplayName,
          email: json.Email,
        });
      });
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
