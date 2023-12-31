import { Box, Button, TablePagination } from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { useState } from "react";
import { useUser } from "../contexts/user.provider";
import { usePaginationUser } from "../swrs/user";
import { User } from "../types/user";
import UserUpdateRoleDialog from "./UserUpdateRoleDialog";

export function UserTable() {
  const { user } = useUser();

  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [order, _] = useState("ASC");

  const { data, isLoading, isError } = usePaginationUser(
    page + 1,
    limit,
    order,
    "id"
  );

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error</div>;
  if (!data || data.items === null) return null;
  if (!user) return null;

  return (
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
          <TableBody>
            {data.items.map((user) => (
              <UserTableRow key={user.id} user={user} />
            ))}
          </TableBody>
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
