import { Box, MenuItem, Select } from "@mui/material";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import { useState } from "react";
import { toast } from "react-toastify";
import { EUserRole, UserService } from "../apis/user";
import { useUser } from "../contexts/user.provider";
import { User } from "../types/user";

interface UserUpdateRoleDialogProps {
  open: boolean;
  user: User;
  response: (user: User) => void;
}

export default function UserUpdateRoleDialog(props: UserUpdateRoleDialogProps) {
  const { user } = useUser();
  const [userRole, setUserRole] = useState<EUserRole>(
    props.user.role as EUserRole
  );

  if (!user) return null;

  const handleClose = () => {
    props.response(props.user);
  };

  const handleUserUpdateRole = async () => {
    try {
      const response = await UserService.updateRole(
        user.accessToken,
        props.user.id,
        userRole
      );

      if (response.ok) {
        props.response(props.user);
        toast.success("User updated");
      } else {
        const data = await response.json();
        toast.error(data.message);
      }
    } catch (error) {
      toast.error("Something went wrong");
    }
  };

  return (
    <>
      <Dialog
        open={props.open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle id="alert-dialog-title">Update User</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Update user {props.user.displayname}
          </DialogContentText>

          <Box mt={2}>
            <Select
              value={userRole}
              onChange={(e) => setUserRole(e.target.value as EUserRole)}
            >
              <MenuItem value={EUserRole.ADMIN}>ADMIN</MenuItem>
              <MenuItem value={EUserRole.STAFF}>STAFF</MenuItem>
              <MenuItem value={EUserRole.USER}>USER</MenuItem>
            </Select>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => handleUserUpdateRole()}>Update</Button>
          <Button onClick={() => handleClose()}>Cancel</Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
