import {
  Box,
  Button,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import ReactMarkdown from "react-markdown";
import { useNavigate, useParams } from "react-router-dom";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";
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

          <ReactMarkdown
            children={data.challenge.description}
            components={{
              code(props) {
                const { children, className, node, ...rest } = props;
                const match = /language-(\w+)/.exec(className || "");
                return match ? (
                  <SyntaxHighlighter
                    PreTag="div"
                    children={String(children).replace(/\n$/, "")}
                    language={match[1]}
                    style={vscDarkPlus}
                  />
                ) : (
                  <code {...rest} className={className}>
                    {children}
                  </code>
                );
              },
            }}
          />
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
