import {
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { Navbar } from "../components/Navbar";
import { SubmissionTable } from "../components/SubmissionTable";

export default function SubmissionPage() {
  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ padding: 3, mt: 15 }}>
          <Typography variant="h4" component="h1" align="left">
            Submission
          </Typography>
          <Divider sx={{ my: 3 }} />

          <SubmissionTable />
        </Paper>
      </Container>
    </Container>
  );
}
