import {
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { ChallengeTable } from "../components/ChallengeTable";
import { Navbar } from "../components/Navbar";

export default function DashboardPage() {
  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ mt: 20, padding: 5 }}>
          <Typography variant="h4" component="h1" align="left">
            Challenge
          </Typography>
          <Divider sx={{ my: 3 }} />

          <ChallengeTable />
        </Paper>
      </Container>
    </Container>
  );
}
