import CancelIcon from "@mui/icons-material/Cancel";
import CheckBoxIcon from "@mui/icons-material/CheckBox";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import RemoveRedEyeIcon from "@mui/icons-material/RemoveRedEye";
import { Box, Button, TablePagination, TextField } from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { useState } from "react";
import { useDebounce } from "use-debounce";
import { usePaginationSubmission } from "../swrs/submission";
import { Submission } from "../types/submission";

function ShowStatusIcon(status: string) {
  switch (status) {
    case "CORRECT":
      return <CheckBoxIcon color="success" />;
    case "PENDING":
      return <HourglassEmptyIcon color="warning" />;
    default:
      return <CancelIcon color="error" />;
  }
}

export function SubmissionTable() {
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [order, _] = useState("desc");
  const [search, setSearch] = useState("");

  const [searchDebounce] = useDebounce(search, 500);

  const { data, isLoading, isError } = usePaginationSubmission(
    page + 1,
    limit,
    "id",
    order,
    searchDebounce,
    0,
    0
  );

  if (isLoading) return <div>Loading...</div>;

  function renderData() {
    if (isError) {
      return (
        <TableRow sx={{ "&:last-child td, &:last-child th": { border: 0 } }}>
          <TableCell
            component="th"
            colSpan={4}
            align="center"
            sx={{ paddingY: 2 }}
          >
            Error
          </TableCell>
        </TableRow>
      );
    }

    if (!data || !data.items || isLoading) {
      return (
        <TableRow sx={{ "&:last-child td, &:last-child th": { border: 0 } }}>
          <TableCell
            component="th"
            colSpan={4}
            align="center"
            sx={{ paddingY: 2 }}
          >
            Loading
          </TableCell>
        </TableRow>
      );
    }

    return data.items.map((submission: Submission) => (
      <SubmissionTableRow
        submission={submission}
        key={submission.submission_id}
      />
    ));
  }

  return (
    <>
      <Box my={2}>
        <TextField
          fullWidth
          id="fullWidth"
          label="Search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </Box>

      <Paper sx={{ width: "100%", overflow: "hidden" }}>
        <TableContainer>
          <Table sx={{ minWidth: 650 }} aria-label="simple table">
            <TableHead>
              <TableRow>
                <TableCell># Submission ID</TableCell>
                <TableCell>Challenge Name</TableCell>
                <TableCell align="right">Language</TableCell>
                <TableCell align="right">Created By</TableCell>
                <TableCell align="right">Status</TableCell>
                <TableCell align="right">Action</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{renderData()}</TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          align="right"
          component="div"
          count={data.total}
          page={page}
          rowsPerPage={limit}
          onPageChange={(_, newPage) => {
            setPage(newPage);
          }}
          onRowsPerPageChange={(e) => setLimit(parseInt(e.target.value, 10))}
        />
      </Paper>
    </>
  );
}

interface SubmissionTableRowProps {
  submission: Submission;
}

function SubmissionTableRow(props: SubmissionTableRowProps) {
  return (
    <TableRow
      key={props.submission.submission_id}
      sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
    >
      <TableCell component="th" scope="row">
        {props.submission.submission_id}
      </TableCell>
      <TableCell component="th" scope="row">
        {props.submission.challenge.name}
      </TableCell>
      <TableCell component="th" scope="row" align="right">
        {props.submission.language}
      </TableCell>
      <TableCell component="th" scope="row" align="right">
        {props.submission.user.displayname}
      </TableCell>
      <TableCell align="right">
        {ShowStatusIcon(props.submission.status)}
      </TableCell>
      <TableCell align="right">
        <Button
          variant="contained"
          color="primary"
          href={`/submission/${props.submission.submission_id}`}
        >
          <RemoveRedEyeIcon sx={{ marginRight: 2 }} /> View
        </Button>
      </TableCell>
    </TableRow>
  );
}
