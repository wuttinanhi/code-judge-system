import { toast } from "react-toastify";
import { API_URL } from "./API_URL";

export class ChallengeService {
  static async delete(accessToken: string, challengeID: number) {
    const response = await fetch(API_URL + "/challenge/delete/" + challengeID, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (response.ok) {
      toast.success("Challenge deleted successfully");
      return true;
    } else {
      const data = await response.json();
      toast.error(`Something went wrong ${data.message}`);
      return false;
    }
  }
}
