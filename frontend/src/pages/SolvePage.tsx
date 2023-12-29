import { Editor } from "@monaco-editor/react";
import {
  Box,
  Button,
  Container,
  CssBaseline,
  Divider,
  MenuItem,
  Paper,
  Select,
  Typography,
} from "@mui/material";
import { useState } from "react";
import { useParams } from "react-router-dom";
import { toast } from "react-toastify";
import { SubmissionService } from "../apis/submission";
import { Navbar } from "../components/Navbar";
import { TestcaseRenderer } from "../components/TestcaseRenderer";
import { useUser } from "../contexts/user.provider";
import { useChallenge } from "../swrs/challenge";
import { ITestcase } from "../types/testcase";

export default function SolvePage() {
  const params = useParams<{ id: string }>();
  const id = parseInt(params.id as any);

  const { data, isError, isLoading } = useChallenge(id);
  const { user } = useUser();

  const [language, setLanguage] = useState("python");
  const [code, setCode] = useState<string>("");

  const [submitButtonDisabled, setSubmitButtonDisabled] = useState(false);

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error</div>;
  if (!user) return <div>Not logged in</div>;

  const onSubmit = () => {
    console.log(language);
    console.log(code);

    try {
      setSubmitButtonDisabled(true);
      SubmissionService.submit(user.accessToken, id, code, language);
      toast.success("Submitted!");
    } catch (error) {
      toast.error("Failed to submit " + error);
    } finally {
      setSubmitButtonDisabled(false);
    }
  };

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ padding: 3, mt: 15 }}>
          <Typography variant="h4" component="h1" align="left">
            #{data.challenge_id} {data.name}
          </Typography>
          <Divider sx={{ my: 3 }} />
          <Typography variant="body1" align="left">
            {data.description}
          </Typography>
        </Paper>

        <Paper sx={{ padding: 3, mt: 5 }}>
          <Typography variant="h4" component="h1" align="left">
            Testcases
          </Typography>
          <Divider sx={{ my: 3 }} />
          <TestcaseRenderer
            testcases={data.testcases.map((t) => {
              return {
                testcase_id: t.testcase_id,
                input: t.input,
                expected_output: t.expected_output,
              } as ITestcase;
            })}
          />
        </Paper>

        <Paper sx={{ padding: 3, mt: 5 }}>
          <Typography variant="h4" component="h1" align="left">
            Solution
          </Typography>
          <Divider sx={{ my: 3 }} />

          <Editor
            height="50vh"
            theme="vs-dark"
            language={language}
            value={code}
            options={{
              minimap: { enabled: false },
            }}
            onChange={(value, _) => setCode(value ?? "")}
          />

          <Box
            display="flex"
            justifyContent="flex-end"
            alignItems="center"
            mt={5}
            gap={2}
          >
            <Select
              value={language}
              onChange={(e) => setLanguage(e.target.value as string)}
            >
              <MenuItem value="python">Python</MenuItem>
              <MenuItem value="c">C</MenuItem>
              <MenuItem value="go">Go</MenuItem>
            </Select>

            <Button
              variant="contained"
              color="primary"
              size="large"
              onClick={onSubmit}
              disabled={submitButtonDisabled}
            >
              Submit
            </Button>
          </Box>
        </Paper>
      </Container>
    </Container>
  );
}
