import { API_URL } from "./API_URL";

export enum EUserRole {
  ADMIN = "ADMIN",
  STAFF = "STAFF",
  USER = "USER",
}

export class UserService {
  static async login(email: string, password: string) {
    return fetch(API_URL + "/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
  }

  static async register(email: string, password: string, displayname: string) {
    return fetch(API_URL + "/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password, displayname }),
    });
  }

  static async updateRole(
    accessToken: string,
    targetUserID: number,
    role: EUserRole
  ) {
    return fetch(API_URL + "/user/update/role", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      },
      body: JSON.stringify({ userid: targetUserID, role }),
    });
  }

  static async getUserInfo(accessToken: string) {
    return fetch(`${API_URL}/user/me`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
  }
}
