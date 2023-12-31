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
import { useUser } from "../contexts/user.provider";
import { usePaginationChallenge } from "../swrs/challenge";
import { Challenge } from "../types/challenge";
import { ShowStatusIcon } from "./StatusIcon";

export function ChallengeTable() {
  const { user } = useUser();

  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [search, setSearch] = useState("");
  const [order, _] = useState("asc");

  const [searchDebounce] = useDebounce(search, 500);

  const { data, isError: error } = usePaginationChallenge(
    page + 1,
    limit,
    order,
    "id",
    searchDebounce
  );

  if (!user) return null;

  function renderData() {
    if (error) {
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

    if (!data || !data.items) {
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

    return data.items.map((c: Challenge) => (
      <ChallengeTableRow challenge={c} key={c.challenge_id} />
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
                <TableCell>Challenge Name</TableCell>
                <TableCell align="right">Created By</TableCell>
                <TableCell align="right">Status</TableCell>
                <TableCell align="right">Action</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{renderData()}</TableBody>
          </Table>

          {data && (
            <TablePagination
              align="right"
              component="div"
              count={data.total}
              page={page}
              rowsPerPage={limit}
              onPageChange={(_, newPage) => {
                setPage(newPage);
              }}
              onRowsPerPageChange={(e) =>
                setLimit(parseInt(e.target.value, 10))
              }
            />
          )}
        </TableContainer>
      </Paper>
    </>
  );
}

interface ChallengeTableRowProps {
  challenge: Challenge;
}

function ChallengeTableRow({ challenge }: ChallengeTableRowProps) {
  const { user } = useUser();

  return (
    <>
      <TableRow
        key={challenge.challenge_id}
        sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
      >
        <TableCell component="th" scope="row">
          {challenge.name}
        </TableCell>
        <TableCell component="th" scope="row" align="right">
          <strong>{challenge.user.displayname}</strong>
        </TableCell>
        <TableCell align="right">
          {ShowStatusIcon(challenge.submission_status)}
        </TableCell>
        <TableCell align="right">
          <Box display="flex" justifyContent="flex-end" gap={1}>
            {user && user.role === "ADMIN" ? (
              <>
                <Button
                  variant="contained"
                  color="warning"
                  href={`/challenge/edit/${challenge.challenge_id}`}
                >
                  Edit
                </Button>
              </>
            ) : null}

            {user && (
              <Button
                variant="contained"
                color="primary"
                href={`/solve/${challenge.challenge_id}`}
              >
                Solve
              </Button>
            )}
          </Box>
        </TableCell>
      </TableRow>
    </>
  );
}
