import { API_URL } from "./API_URL";

export class SubmissionService {
  static async submit(
    accessToken: string,
    challengeID: number,
    code: string,
    language: string
  ) {
    return fetch(API_URL + "/submission/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify({
        challenge_id: challengeID,
        language: language,
        code: code,
      }),
    });
  }
}
