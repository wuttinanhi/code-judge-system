import {
  Box,
  Button,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { EUserRole } from "../apis/user";
import { ChallengeTable } from "../components/ChallengeTable";
import { Navbar } from "../components/Navbar";
import { useUser } from "../contexts/user.provider";

export default function DashboardPage() {
  const navigate = useNavigate();
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

                {user &&
                (user.role === EUserRole.ADMIN ||
                  user.role === EUserRole.STAFF) ? (
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => navigate(`/challenge/create`)}
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
