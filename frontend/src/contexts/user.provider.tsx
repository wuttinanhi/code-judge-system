import { ReactNode, createContext, useContext, useState } from "react";

type UserDataType = {
  accessToken: string;
  username: string;
};

interface UserContextType {
  user: UserDataType | undefined;
  setUser: (user: UserDataType) => void;
}

const UserContext = createContext<UserContextType>({
  setUser: () => {},
  user: undefined,
});

export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<UserDataType | undefined>(undefined);

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
