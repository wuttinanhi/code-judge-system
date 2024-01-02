import { toast } from "react-toastify";
import { BadRequestResponse } from "../types/http";

export function handleBadRequest(data: BadRequestResponse) {
  // loop through the errors array and display each error
  data.errors.forEach((error) => {
    toast.error(error);
  });
}
