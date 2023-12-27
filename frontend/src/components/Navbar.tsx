import { AppBar, Button } from "@mui/material";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { useUser } from "../contexts/user.provider";

export function Navbar() {
  const userContext = useUser();

  return (
    <AppBar sx={{ mb: 10 }}>
      <Toolbar>
        <Typography
          variant="h6"
          noWrap
          component="a"
          href="/"
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

        <Button
          sx={{ my: 2, color: "white", display: "block" }}
          href="/challenge"
        >
          Challenge
        </Button>

        <Button
          sx={{ my: 2, color: "white", display: "block" }}
          href="/submission"
        >
          Submission
        </Button>

        {userContext.user ? (
          <Button color="inherit" sx={{ marginLeft: "auto" }} href="/settings">
            {userContext.user.displayName}
          </Button>
        ) : (
          <Button color="inherit" sx={{ marginLeft: "auto" }} href="/signin">
            Login
          </Button>
        )}
      </Toolbar>
    </AppBar>
  );
}
