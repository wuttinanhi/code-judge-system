import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { ThemeProvider, createTheme } from "@mui/material";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "./App.css";
import { BackendPinger } from "./components/BackendPinger";
import { UserProvider } from "./contexts/user.provider";
import ChallengeModifyPage from "./pages/ChallengeModifyPage";
import DashboardPage from "./pages/DashboardPage";
import { SignInPage } from "./pages/SigninPage";
import { SignUpPage } from "./pages/SignupPage";
import SolvePage from "./pages/SolvePage";
import SubmissionPage from "./pages/SubmissionPage";
import SubmissionViewPage from "./pages/SubmissionViewPage";
import UserManagePage from "./pages/UserManagePage";
import UserSettingPage from "./pages/UserSettingPage";

// TODO remove, this demo shouldn't need to reset the theme.
const defaultTheme = createTheme();

const router = createBrowserRouter([
  {
    path: "/",
    element: <DashboardPage />,
  },
  {
    path: "/signin",
    element: <SignInPage />,
  },
  {
    path: "/signup",
    element: <SignUpPage />,
  },
  {
    path: "/challenge",
    element: <DashboardPage />,
  },
  {
    path: "/challenge/:mode",
    element: <ChallengeModifyPage />,
  },
  {
    path: "/challenge/:mode/:id",
    element: <ChallengeModifyPage />,
  },
  {
    path: "/submission",
    element: <SubmissionPage />,
  },
  {
    path: "/submission/:id",
    element: <SubmissionViewPage />,
  },
  {
    path: "/solve/:id",
    element: <SolvePage />,
  },
  {
    path: "/settings",
    element: <UserSettingPage />,
  },
  {
    path: "/admin/user",
    element: <UserManagePage />,
  },
]);

function App() {
  return (
    <>
      <div>
        <ToastContainer />
        <BackendPinger />
        <ThemeProvider theme={defaultTheme}>
          <UserProvider>
            <RouterProvider router={router} />
          </UserProvider>
        </ThemeProvider>
      </div>
    </>
  );
}

export default App;
