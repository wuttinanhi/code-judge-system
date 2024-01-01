import {
  Box,
  Button,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import { Navbar } from "../components/Navbar";
import { TestcaseRenderer } from "../components/TestcaseRenderer";
import { useSubmission } from "../swrs/submission";
import { ITestcase } from "../types/testcase";

export default function SubmissionViewPage() {
  const navigate = useNavigate();
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
          <Box justifyContent="space-between" display="flex">
            <Typography variant="h4" align="left">
              Submission #{data.submission_id} {data.challenge.name}
            </Typography>

            <Button
              variant="contained"
              color="primary"
              onClick={() => navigate(`/solve/${data.challenge_id}`)}
            >
              Go to Challenge
            </Button>
          </Box>
          <Divider sx={{ my: 3 }} />
          {data.challenge.description}
        </Paper>

        <Paper sx={{ padding: 3, mt: 5 }}>
          <Typography variant="h6" align="left">
            Testcase Results
          </Typography>
          <Divider sx={{ my: 3 }} />

          <TestcaseRenderer
            testcases={data.submission_testcases.map((t) => {
              return {
                testcase_id: t.challenge_testcase.testcase_id,
                correct: t.status,
                input: t.challenge_testcase.input,
                output: t.output === "" ? "<EMPTY OUTPUT>" : t.output,
                expected_output: t.challenge_testcase.expected_output,
              } as ITestcase;
            })}
          />
        </Paper>
      </Container>
    </Container>
  );
}
