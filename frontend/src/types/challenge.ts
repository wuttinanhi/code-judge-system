// Generated by https://quicktype.io

import { ITestcaseModify } from "./testcase";
import { User } from "./user";

export interface Challenge {
  challenge_id: number;
  name: string;
  description: string;
  user_id: number;
  testcases: ChallengeTestcase[];
  submission: null;
  user: User;
  submission_status: string;
}

export interface ChallengeTestcase {
  testcase_id: number;
  input: string;
  expected_output: string;
  limit_memory: number;
  limit_time_ms: number;
  submission_testcases: null;
  challenge_id: number;
  challenge: null;
}

export interface ChallengeUpdateDTO {
  challenge_id: number;
  name: string;
  description: string;
  testcases: ITestcaseModify[];
}
