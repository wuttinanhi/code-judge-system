import useSWR from "swr";
import { PaginationResult } from "../types/pagination";
import { User } from "../types/user";
import { fetcherWithAuth } from "./fetcher";

export function usePaginationUser(
  page: number,
  limit: number,
  order: string,
  sort: string
) {
  const { data, error, isLoading } = useSWR(
    () =>
      `/user/pagination?page=${page}&limit=${limit}&sort=${sort}&order=${order}`,
    fetcherWithAuth
  );

  return {
    data: data as PaginationResult<User>,
    isLoading,
    isError: error,
  };
}
