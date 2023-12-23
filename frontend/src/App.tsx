import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { ThemeProvider, createTheme } from "@mui/material";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "./App.css";
import { UserProvider } from "./contexts/user.provider";
import DashboardPage from "./pages/DashboardPage";
import { SignInPage } from "./pages/SigninPage";
import { SignUpPage } from "./pages/SignupPage";

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
]);

function App() {
  return (
    <>
      <div>
        <ToastContainer />
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
