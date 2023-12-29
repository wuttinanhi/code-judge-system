import CancelIcon from "@mui/icons-material/Cancel";
import CheckBoxIcon from "@mui/icons-material/CheckBox";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";

export function ShowStatusIcon(status: string) {
  switch (status) {
    case "CORRECT":
      return <CheckBoxIcon color="success" />;
    case "PENDING":
      return <HourglassEmptyIcon color="warning" />;
    default:
      return <CancelIcon color="error" />;
  }
}
