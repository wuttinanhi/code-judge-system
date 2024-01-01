import {
  Box,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { Navbar } from "../components/Navbar";
import { UserTable } from "../components/UserTable";

export default function UserManagePage() {
  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Container>
        <Paper sx={{ padding: 3, mt: 15 }}>
          <Box justifyContent="space-between" display="flex">
            <Typography variant="h4" component="h1" align="left">
              User Management
            </Typography>
          </Box>

          <Divider sx={{ my: 3 }} />

          <UserTable />
        </Paper>
      </Container>
    </Container>
  );
}
