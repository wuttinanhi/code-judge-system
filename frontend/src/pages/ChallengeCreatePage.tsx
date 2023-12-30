import {
  Box,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { ChallengeEditor } from "../components/ChallengeEditor";
import { Navbar } from "../components/Navbar";
import { useUser } from "../contexts/user.provider";

export default function ChallengeCreatePage() {
  const { user } = useUser();

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ padding: 3, mt: 15 }}>
          <Box justifyContent="space-between" display="flex">
            <Typography variant="h4" component="h1" align="left">
              Create New Challenge
            </Typography>
          </Box>

          <Divider sx={{ my: 3 }} />

          <ChallengeEditor />
        </Paper>
      </Container>
    </Container>
  );
}
