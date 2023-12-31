import { Card, CardContent, Stack, Typography } from "@mui/material";
import { ITestcase } from "../types/testcase";

interface TestcaseProps {
  testcase: ITestcase;
}

export function Testcase(props: TestcaseProps) {
  return (
    <Stack
      direction={{ xs: "column", md: "row" }}
      justifyContent="space-evenly"
      alignItems="stretch"
      spacing={5}
    >
      <Card sx={{ flexGrow: 1, flexBasis: "50%", width: "full" }}>
        <CardContent>
          <Typography
            sx={{ fontSize: 14, marginBottom: 2 }}
            color="text.secondary"
            gutterBottom
          >
            Input
          </Typography>

          <Typography
            variant="subtitle1"
            component="div"
            sx={{ backgroundColor: "black", color: "white", padding: 2 }}
          >
            {props.testcase.input}
          </Typography>
        </CardContent>
      </Card>

      {props.testcase.output ? (
        <Card sx={{ flexGrow: 1, flexBasis: "50%", width: "full" }}>
          <CardContent>
            <Typography
              sx={{ fontSize: 14, marginBottom: 2 }}
              color="text.secondary"
              gutterBottom
            >
              Output
            </Typography>

            <Typography
              variant="subtitle1"
              component="div"
              sx={{ backgroundColor: "black", color: "white", padding: 2 }}
            >
              {props.testcase.output}
            </Typography>
          </CardContent>
        </Card>
      ) : null}

      <Card sx={{ flexGrow: 1, flexBasis: "50%", width: "full" }}>
        <CardContent>
          <Typography
            sx={{ fontSize: 14, marginBottom: 2 }}
            color="text.secondary"
            gutterBottom
          >
            Expected Output
          </Typography>

          <Typography
            variant="subtitle1"
            component="div"
            sx={{ backgroundColor: "black", color: "white", padding: 2 }}
          >
            {props.testcase.expected_output}
          </Typography>
        </CardContent>
      </Card>
    </Stack>
  );
}
