import {
  Box,
  Button,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { ChallengeTable } from "../components/ChallengeTable";
import { Navbar } from "../components/Navbar";
import { useUser } from "../contexts/user.provider";

export default function DashboardPage() {
  const { user } = useUser();

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        {user ? (
          <>
            <Paper sx={{ padding: 3, mt: 15 }}>
              <Box justifyContent="space-between" display="flex">
                <Typography variant="h4" component="h1" align="left">
                  Challenge
                </Typography>

                {user && user.role === "ADMIN" ? (
                  <Button
                    variant="contained"
                    color="primary"
                    href={`/challenge/create`}
                  >
                    Create
                  </Button>
                ) : null}
              </Box>

              <Divider sx={{ my: 3 }} />

              <ChallengeTable />
            </Paper>
          </>
        ) : (
          <>
            <Paper sx={{ padding: 3, mt: 15 }}>
              <Typography variant="h4" component="h1" align="left">
                Welcome to Code Judge System
              </Typography>

              <Divider sx={{ my: 3 }} />

              <Typography variant="body1" align="left">
                Please login to continue
              </Typography>
            </Paper>
          </>
        )}
      </Container>
    </Container>
  );
}
