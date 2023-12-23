export class UserService {
  static async login(email: string, password: string) {
    const response = await fetch("http://localhost:3000/user/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (response.ok) {
      return (await response.json()) as UserLoginResponse;
    } else {
      throw new Error("Something went wrong");
    }
  }

  static async register(
    username: string,
    password: string,
    displayname: string
  ) {
    const response = await fetch("http://localhost:3000/user/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password, displayname }),
    });
    if (response.ok) {
      return (await response.json()) as UserRegisterResponse;
    } else {
      throw new Error("Something went wrong");
    }
  }
}

export interface UserLoginResponse {
  token: string;
  userid: number;
  displayname: string;
  email: string;
}

export interface UserRegisterResponse {
  userid: number;
  displayname: string;
  email: string;
}
