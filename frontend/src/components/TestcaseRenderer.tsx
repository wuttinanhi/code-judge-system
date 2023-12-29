import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Typography,
} from "@mui/material";
import { ITestcase } from "../types/testcase";
import { ShowStatusIcon } from "./StatusIcon";
import { Testcase } from "./Testcase";

interface TestcaseRendererProps {
  testcases: ITestcase[];
}

export function TestcaseRenderer(props: TestcaseRendererProps) {
  return props.testcases.map((testcase, index) => {
    return (
      <Accordion defaultExpanded key={index}>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>Testcase #{testcase.testcase_id}</Typography>
          {testcase.correct ? ShowStatusIcon(testcase.correct) : null}
        </AccordionSummary>
        <AccordionDetails>
          <Testcase testcase={testcase} />
        </AccordionDetails>
      </Accordion>
    );
  });
}
