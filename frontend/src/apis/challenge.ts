import { ChallengeUpdateDTO } from "../types/challenge";
import { ITestcaseModify } from "../types/testcase";
import { API_URL } from "./API_URL";

export class ChallengeService {
  static async delete(accessToken: string, challengeID: number) {
    return fetch(API_URL + "/challenge/delete/" + challengeID, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
    });
  }

  static async create(
    accessToken: string,
    challengeName: string,
    challengeDescription: string,
    testcases: ITestcaseModify[]
  ) {
    return fetch(`${API_URL}/challenge/create`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify({
        name: challengeName,
        description: challengeDescription,
        testcases: testcases,
      }),
    });
  }

  static async edit(accessToken: string, challenge: ChallengeUpdateDTO) {
    return fetch(`${API_URL}/challenge/update/${challenge.challenge_id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify(challenge),
    });
  }

  static async get(accessToken: string, challengeID: any) {
    return fetch(`${API_URL}/challenge/get/${challengeID}`, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
  }
}
