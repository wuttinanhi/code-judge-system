import useSWR from "swr";
import { Challenge } from "../types/challenge";
import { PaginationResult } from "../types/pagination";
import { fetcherWithAuth } from "./fetcher";

export function useChallenge(
  page: number,
  limit: number,
  order: string,
  sort: string
) {
  const { data, error, isLoading } = useSWR(
    () =>
      `/challenge/pagination?page=${page}&limit=${limit}&sort=${sort}&order=${order}`,
    fetcherWithAuth
  );

  return {
    data: data as PaginationResult<Challenge>,
    isLoading,
    isError: error,
  };
}
