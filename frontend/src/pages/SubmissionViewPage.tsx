import {
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { Navbar } from "../components/Navbar";
import { TestcaseRenderer } from "../components/TestcaseRenderer";
import { useSubmission } from "../swrs/submission";
import { ITestcase } from "../types/testcase";

export default function SubmissionViewPage() {
  const params = useParams<{ id: string }>();
  const id = parseInt(params.id as any);

  const { data, isError, isLoading } = useSubmission(id);

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error</div>;

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ padding: 3, mt: 15 }}>
          <Typography variant="h4" component="h1" align="left">
            #{data.submission_id} {data.challenge.name}
          </Typography>
          <Divider sx={{ my: 3 }} />
          <Typography variant="body1" align="left">
            {data.challenge.description}
          </Typography>
        </Paper>

        <Paper sx={{ padding: 3, mt: 5 }}>
          <Typography variant="h4" component="h1" align="left">
            Testcase Results
          </Typography>
          <Divider sx={{ my: 3 }} />
          <Typography variant="body1" align="left">
            <TestcaseRenderer
              testcases={data.submission_testcases.map((t) => {
                return {
                  testcase_id: t.challenge_testcase.testcase_id,
                  correct: t.status,
                  input: t.challenge_testcase.input,
                  output: t.output,
                  expected_output: t.challenge_testcase.expected_output,
                } as ITestcase;
              })}
            />
          </Typography>
        </Paper>
      </Container>
    </Container>
  );
}
