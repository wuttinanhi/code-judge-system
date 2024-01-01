import { AppBar, Button } from "@mui/material";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { useNavigate } from "react-router-dom";
import { EUserRole } from "../apis/user";
import { useUser } from "../contexts/user.provider";

export function Navbar() {
  const userContext = useUser();
  const navigate = useNavigate();

  return (
    <AppBar sx={{ mb: 10 }}>
      <Toolbar>
        <Typography
          variant="h6"
          noWrap
          component="a"
          onClick={() => navigate("/")}
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

        {userContext.user && (
          <>
            <Button
              sx={{ my: 2, color: "white", display: "block" }}
              onClick={() => navigate("/challenge")}
            >
              Challenge
            </Button>

            <Button
              sx={{ my: 2, color: "white", display: "block" }}
              onClick={() => navigate("/submission")}
            >
              Submission
            </Button>
          </>
        )}

        {userContext.user && userContext.user.role === EUserRole.ADMIN ? (
          <Button
            sx={{ my: 2, color: "white", display: "block" }}
            onClick={() => navigate("/admin/user")}
          >
            User
          </Button>
        ) : null}

        {userContext.user ? (
          <Button
            color="inherit"
            sx={{ marginLeft: "auto" }}
            onClick={() => navigate("/settings")}
          >
            {userContext.user.displayName}
          </Button>
        ) : (
          <Button
            color="inherit"
            sx={{ marginLeft: "auto" }}
            onClick={() => navigate("/signin")}
          >
            Login
          </Button>
        )}
      </Toolbar>
    </AppBar>
  );
}
