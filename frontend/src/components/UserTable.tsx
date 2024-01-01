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
import { usePaginationUser } from "../swrs/user";
import { User } from "../types/user";
import UserUpdateRoleDialog from "./UserUpdateRoleDialog";

export function UserTable() {
  const { user } = useUser();

  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [search, setSearch] = useState("");
  const [order, _] = useState("ASC");

  const [searchDebounce] = useDebounce(search, 500);

  const { data, isLoading, isError } = usePaginationUser(
    page + 1,
    limit,
    order,
    "id",
    searchDebounce
  );

  if (!user) return null;

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

    return data.items.map((user: User) => (
      <UserTableRow user={user} key={user.id} />
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
                <TableCell>ID</TableCell>
                <TableCell align="right">Display Name</TableCell>
                <TableCell align="right">Email</TableCell>
                <TableCell align="right">Action</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{renderData()}</TableBody>
          </Table>
        </TableContainer>
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
            onRowsPerPageChange={(e) => setLimit(parseInt(e.target.value, 10))}
          />
        )}
      </Paper>
    </>
  );
}

interface UserTableRowProps {
  user: User;
}

function UserTableRow({ user }: UserTableRowProps) {
  const [dialogOpen, setDialogOpen] = useState(false);

  return (
    <>
      <TableRow
        key={user.id}
        sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
      >
        <TableCell component="th" scope="row">
          {user.id}
        </TableCell>
        <TableCell component="th" scope="row" align="right">
          {user.displayname}
        </TableCell>
        <TableCell component="th" scope="row" align="right">
          {user.email}
        </TableCell>
        <TableCell align="right">
          <Box display="flex" justifyContent="flex-end" gap={1}>
            <Button
              variant="contained"
              color="warning"
              onClick={() => {
                setDialogOpen(true);
              }}
            >
              Update
            </Button>
          </Box>
        </TableCell>
      </TableRow>

      <UserUpdateRoleDialog
        open={dialogOpen}
        user={user}
        response={() => {
          setDialogOpen(false);
        }}
      />
    </>
  );
}
