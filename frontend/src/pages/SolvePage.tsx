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
import ReactMarkdown from "react-markdown";
import { useParams } from "react-router-dom";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";
import { toast } from "react-toastify";
import { SubmissionService } from "../apis/submission";
import { Navbar } from "../components/Navbar";
import { TestcaseRenderer } from "../components/TestcaseRenderer";
import { useUser } from "../contexts/user.provider";
import { handleBadRequest } from "../helpers/badrequest-toast";
import { useChallenge } from "../swrs/challenge";
import { SubmissionSubmitResponse } from "../types/submission";
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

  const onSubmit = async () => {
    console.log(language);
    console.log(code);

    setSubmitButtonDisabled(true);
    const response = await SubmissionService.submit(
      user.accessToken,
      id,
      code,
      language
    ).then();

    if (response.ok) {
      toast.success("Submitted!");
      const data: SubmissionSubmitResponse = await response.json();
      window.location.href = `/submission/${data.submission_id}`;
    } else {
      const data = await response.json();
      handleBadRequest(data);
    }

    setSubmitButtonDisabled(false);
  };

  return (
    <>
      <Navbar />

      <Box
        display="flex"
        flexDirection="row"
        justifyContent="flex-start"
        alignContent="flex-start"
        sx={{ marginTop: 15, marginX: 10 }}
        columnGap={5}
      >
        <Box flex={1}>
          <Paper sx={{ padding: 3 }}>
            <Typography variant="h4" component="h1" align="left">
              #{data.challenge_id} {data.name}
            </Typography>
            <Divider sx={{ my: 3 }} />
            <ReactMarkdown
              children={data.description}
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

          <Box mt={5}>
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
          </Box>
        </Box>

        <Box flex={1} maxHeight="80vh">
          <Typography variant="h4" component="h1" align="left">
            Solution
          </Typography>

          <Divider sx={{ my: 3 }} />

          <Editor
            height="80%"
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
        </Box>
      </Box>

      <Container maxWidth="lg">
        <CssBaseline />
      </Container>
    </>
  );
}
