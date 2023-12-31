import { useEffect } from "react";
import { toast } from "react-toastify";
import { API_URL } from "../apis/API_URL";

export function BackendPinger() {
  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    fetch(API_URL, { signal }).catch((err) => {
      const error = err as Error;

      if (error.message === "Failed to fetch") {
        toast.error("Failed to connect to server");
      }
    });

    return () => {
      controller.abort();
    };
  }, []);
  return null;
}
