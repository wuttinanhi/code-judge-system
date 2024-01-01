import DeleteForeverIcon from "@mui/icons-material/DeleteForever";
import {
  Box,
  Button,
  Card,
  CardContent,
  Divider,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { ChallengeService } from "../apis/challenge";
import { useUser } from "../contexts/user.provider";
import { ChallengeUpdateDTO } from "../types/challenge";
import { ITestcaseModify } from "../types/testcase";
import AlertDialog from "./AlertDialog";

interface ChallengeEditorProps {
  mode: "create" | "edit";
  editChallengeID?: number;
}

export function ChallengeEditor(props: ChallengeEditorProps) {
  const [submitButtonDisabled, setSubmitButtonDisabled] = useState(false);
  const [testcases, setTestcases] = useState<ITestcaseModify[]>([]);
  const [challengeName, setChallengeName] = useState("");
  const [challengeDescription, setChallengeDescription] = useState("");

  const [deleteButtonDisabled, setDeleteButtonDisabled] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const { user } = useUser();

  if (!user) {
    return <div>Not logged in</div>;
  }

  if (props.mode === "edit" && !props.editChallengeID) {
    return <div>Invalid challenge ID</div>;
  }

  const createMode = async () => {
    const res = await ChallengeService.create(
      user.accessToken,
      challengeName,
      challengeDescription,
      testcases
    );

    const data = await res.json();

    if (!res.ok) {
      toast.error(data.message);
    } else {
      // redirect to challenge page
      window.location.href = `/challenge`;
      toast.success(data.message);
    }
  };

  const editMode = async () => {
    const updateData = {
      challenge_id: props.editChallengeID,
      name: challengeName,
      description: challengeDescription,
      testcases: testcases,
    } as ChallengeUpdateDTO;

    console.log(updateData);

    const res = await ChallengeService.edit(user.accessToken, updateData);
    const data = await res.json();

    if (!res.ok) {
      toast.error(data.message);
    } else {
      toast.success("Challenge updated successfully");
    }
  };

  const handleSubmit = async () => {
    try {
      setSubmitButtonDisabled(true);
      if (props.mode === "create") {
        await createMode();
      } else if (props.mode === "edit") {
        await editMode();
      }
    } catch (err) {
      const error = err as Error;
      toast.error(error.message);
    } finally {
      setSubmitButtonDisabled(false);
    }
  };

  function onCreateTestcase(): void {
    setTestcases((prev) => [
      ...prev,
      {
        input: "",
        expected_output: "",
        testcase_id: 0,
        limit_memory: 268435456,
        limit_time_ms: 1000,
        action: "create",
      },
    ]);
  }

  useEffect(() => {
    if (props.mode === "edit") {
      async function loadEditChallenge() {
        if (!user || !props.editChallengeID) return;

        const res = await ChallengeService.get(
          user.accessToken,
          props.editChallengeID
        );
        const data = await res.json();

        if (!res.ok) {
          toast.error(data.message);
        } else {
          setChallengeName(data.name);
          setChallengeDescription(data.description);
          setTestcases(data.testcases);
        }
      }

      loadEditChallenge();
    }
  }, []);

  return (
    <>
      <AlertDialog
        title="Delete Challenge"
        description="Are you sure you want to delete this challenge?"
        open={deleteDialogOpen}
        response={async (res) => {
          if (res) {
            setDeleteDialogOpen(false);
            setDeleteButtonDisabled(true);

            const response = await ChallengeService.delete(
              user.accessToken,
              props.editChallengeID!
            );

            if (response.ok) {
              toast.success("Challenge deleted successfully");
              window.location.href = "/challenge";
            } else {
              const data = await response.json();
              toast.error(`Something went wrong ${data.message}`);
              setDeleteDialogOpen(false);
              setDeleteButtonDisabled(false);
            }
          }
        }}
      />

      <Box display="flex" flexDirection="column" gap={5}>
        <TextField
          label="Challenge Name"
          helperText="Name of the challenge"
          value={challengeName}
          onChange={(e) => setChallengeName(e.target.value)}
        />

        <TextField
          label="Challenge Description"
          multiline
          rows={20}
          helperText="Description of the challenge"
          value={challengeDescription}
          onChange={(e) => setChallengeDescription(e.target.value)}
        />

        <Divider />

        <Box
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          mt={2}
          gap={2}
        >
          <Typography variant="h5" align="left">
            Testcases
          </Typography>

          <Button
            variant="contained"
            color="primary"
            size="medium"
            disabled={submitButtonDisabled}
            onClick={onCreateTestcase}
          >
            Create Testcase
          </Button>
        </Box>

        <Box
          display={"flex"}
          flexDirection={"column"}
          justifyContent={"stretch"}
          gap={2}
        >
          {testcases.map((testcase, index) => (
            <TestcaseEditor
              key={index}
              testcase={testcase}
              onChange={(testcase) => {
                setTestcases((prev) => {
                  const newTestcases = [...prev];
                  newTestcases[index] = testcase;
                  return newTestcases;
                });
              }}
            />
          ))}
        </Box>

        <Box
          display="flex"
          justifyContent="flex-end"
          alignItems="center"
          mt={2}
          gap={2}
        >
          <Button
            variant="contained"
            color="error"
            size="large"
            onClick={() => setDeleteDialogOpen(true)}
            disabled={deleteButtonDisabled}
          >
            Delete
          </Button>

          <Button
            variant="contained"
            color="primary"
            size="large"
            onClick={handleSubmit}
            disabled={submitButtonDisabled}
          >
            Submit
          </Button>
        </Box>
      </Box>
    </>
  );
}

interface TestcaseEditorProps {
  testcase: ITestcaseModify;
  onChange?: (testcase: ITestcaseModify) => void;
}

export function TestcaseEditor(props: TestcaseEditorProps) {
  const [data, setData] = useState<ITestcaseModify>(props.testcase);

  useEffect(() => {
    if (props.onChange) props.onChange(data);
  }, [data, props]);

  if (data.action === "delete") {
    return null;
  }

  return (
    <Stack
      direction="row"
      justifyContent="space-evenly"
      alignItems="stretch"
      spacing={5}
    >
      <Card sx={{ flexGrow: 1, flexBasis: "50%", width: "full" }}>
        <CardContent>
          <Box justifyContent="space-between" display="flex" mb={3}>
            <Typography variant="h6">
              Testcase #{data.testcase_id ? data.testcase_id : "New"}
            </Typography>

            <Box display="flex" justifyContent="flex-end" gap={2}>
              <TextField
                label="Limit Memory"
                variant="outlined"
                value={data.limit_memory}
                onChange={(e) => {
                  setData((p) => ({
                    ...p,
                    limit_memory: parseInt(e.target.value),
                    action: p.action === "create" ? "create" : "update",
                  }));
                }}
              />

              <TextField
                label="Limit Time (ms)"
                variant="outlined"
                value={data.limit_time_ms}
                onChange={(e) => {
                  setData((p) => ({
                    ...p,
                    limit_time_ms: parseInt(e.target.value),
                    action: p.action === "create" ? "create" : "update",
                  }));
                }}
              />

              <Button
                variant="contained"
                color="error"
                onClick={() => {
                  setData((p) => ({ ...p, action: "delete" }));
                }}
              >
                <DeleteForeverIcon />
                Delete
              </Button>
            </Box>
          </Box>

          <Box
            flexDirection={"row"}
            display={"flex"}
            justifyContent={"stretch"}
            width={"full"}
            gap={3}
          >
            <Box sx={{ width: "50%" }}>
              <Typography
                sx={{ fontSize: 14, marginBottom: 2 }}
                color="text.secondary"
                gutterBottom
              >
                Input
              </Typography>

              <TextField
                multiline
                rows={4}
                fullWidth
                value={data.input}
                onChange={(e) =>
                  setData((p) => ({
                    ...p,
                    input: e.target.value,
                    action: p.action === "create" ? "create" : "update",
                  }))
                }
              />
            </Box>

            <Box sx={{ width: "50%" }}>
              <Typography
                sx={{ fontSize: 14, marginBottom: 2 }}
                color="text.secondary"
                gutterBottom
              >
                Expected Output
              </Typography>

              <TextField
                multiline
                rows={4}
                fullWidth
                value={data.expected_output}
                onChange={(e) =>
                  setData((p) => ({
                    ...p,
                    expected_output: e.target.value,
                    action: p.action === "create" ? "create" : "update",
                  }))
                }
              />
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Stack>
  );
}
