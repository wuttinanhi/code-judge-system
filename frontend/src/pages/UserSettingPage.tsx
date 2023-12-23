import {
  Button,
  Container,
  CssBaseline,
  Divider,
  Paper,
  Typography,
} from "@mui/material";
import { Navbar } from "../components/Navbar";
import { useUser } from "../contexts/user.provider";

export default function UserSettingPage() {
  const userContext = useUser();

  const onLogoutClick = () => {
    userContext.setUser(undefined);
    localStorage.removeItem("accessToken");
    localStorage.removeItem("user");
    window.location.href = "/";
  };

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />

      <Navbar />

      <Paper sx={{ mt: 20, padding: 5 }}>
        <Typography variant="h4" component="h1" align="left">
          User Setting
        </Typography>
        <Divider sx={{ my: 3 }} />

        <Button
          variant="contained"
          color="error"
          sx={{ ml: "auto" }}
          onClick={onLogoutClick}
        >
          Logout
        </Button>
      </Paper>
    </Container>
  );
}
