import { AppBar, Button, Container, CssBaseline } from "@mui/material";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { useUser } from "../contexts/user.provider";

export function Dashboard() {
  const userContext = useUser();

  return (
    <Container sx={{ width: "100%" }} disableGutters>
      <CssBaseline />
      <AppBar>
        <Toolbar>
          <Typography
            variant="h6"
            noWrap
            component="a"
            href="#app-bar-with-responsive-menu"
            sx={{
              mr: 5,
              display: { xs: "none", md: "flex" },
              fontWeight: 700,
              color: "inherit",
              textDecoration: "none",
            }}
          >
            CODE JUDGE SYSTEM
          </Typography>

          <Button sx={{ my: 2, color: "white", display: "block" }}>
            Challenge
          </Button>

          <Button sx={{ my: 2, color: "white", display: "block" }}>
            Submission
          </Button>

          {userContext.user ? (
            <Button color="inherit" sx={{ marginLeft: "auto" }}>
              {userContext.user.username}
            </Button>
          ) : (
            <Button color="inherit" sx={{ marginLeft: "auto" }}>
              Login
            </Button>
          )}
        </Toolbar>
      </AppBar>
    </Container>
  );
}
