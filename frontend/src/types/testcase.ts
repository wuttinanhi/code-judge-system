export interface ITestcase {
  testcase_id: any | undefined;
  input: string | undefined;
  output: string | undefined;
  expected_output: string | undefined;
  limit_memory: number;
  limit_time_ms: number;
  correct: string | undefined;
}

export interface ITestcaseModify {
  testcase_id: number | undefined;
  input: string | undefined;
  expected_output: string | undefined;
  limit_memory: number;
  limit_time_ms: number;
  action: "create" | "update" | "delete";
}
