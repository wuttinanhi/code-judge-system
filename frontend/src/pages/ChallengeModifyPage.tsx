import {
  Box,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { ChallengeEditor } from "../components/ChallengeEditor";
import { Navbar } from "../components/Navbar";
import { useUser } from "../contexts/user.provider";

export default function ChallengeModifyPage() {
  const params = useParams<{ id: string; mode: string }>();
  const mode = params.mode as "create" | "edit";
  const challengeID = parseInt(params.id as string);

  const { user } = useUser();

  if (!user) {
    return <h1>Not Logged In</h1>;
  }

  if (!mode) {
    return <h1>Invalid Mode</h1>;
  }

  if ((!challengeID || isNaN(challengeID)) && mode === "edit") {
    return <h1>Invalid Challenge ID</h1>;
  }

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ padding: 3, mt: 15 }}>
          <Box justifyContent="space-between" display="flex">
            <Typography variant="h4" component="h1" align="left">
              {mode === "create" ? "Create Challenge" : "Edit Challenge"}
            </Typography>
          </Box>

          <Divider sx={{ my: 3 }} />

          <ChallengeEditor mode={mode} editChallengeID={challengeID} />
        </Paper>
      </Container>
    </Container>
  );
}
