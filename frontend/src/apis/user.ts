import { API_URL } from "./API_URL";

export class UserService {
  static async login(email: string, password: string) {
    const response = await fetch(API_URL + "/auth/login", {
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

  static async register(email: string, password: string, displayname: string) {
    const response = await fetch(API_URL + "/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password, displayname }),
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
  role: string;
}

export interface UserRegisterResponse {
  userid: number;
  displayname: string;
  email: string;
}
